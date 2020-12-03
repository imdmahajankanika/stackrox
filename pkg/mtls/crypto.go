package mtls

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"math/big"
	"time"

	"github.com/cloudflare/cfssl/config"
	cfcsr "github.com/cloudflare/cfssl/csr"
	"github.com/cloudflare/cfssl/helpers"
	cflog "github.com/cloudflare/cfssl/log"
	cfsigner "github.com/cloudflare/cfssl/signer"
	"github.com/cloudflare/cfssl/signer/local"
	"github.com/pkg/errors"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/errorhelpers"
	"github.com/stackrox/rox/pkg/sync"
)

const (
	certsPrefix = "/run/secrets/stackrox.io/certs/"
	// defaultCACertFilePath is where the certificate is stored.
	defaultCACertFilePath = certsPrefix + "ca.pem"
	// defaultCAKeyFilePath is where the key is stored.
	defaultCAKeyFilePath = certsPrefix + "ca-key.pem"

	// defaultCertFilePath is where the certificate is stored.
	defaultCertFilePath = certsPrefix + "cert.pem"
	// defaultKeyFilePath is where the key is stored.
	defaultKeyFilePath = certsPrefix + "key.pem"

	// To account for clock skew, set certificates to be valid some time in the past.
	beforeGracePeriod = 1 * time.Hour

	certLifetime = 365 * 24 * time.Hour
)

var (
	// serialMax is the max value to be used with `rand.Int` to obtain a `*big.Int` with 64 bits of random data
	// (i.e., 1 << 64).
	serialMax = func() *big.Int {
		max := big.NewInt(1)
		max.Lsh(max, 64)
		return max
	}()
)

func init() {
	// The cfssl library prints logs at Info level when it processes a
	// Certificate Signing Request (CSR) or issues a new certificate.
	// These logs do not help the user understand anything, so here
	// we adjust the log level to exclude them.
	cflog.Level = cflog.LevelWarning
}

var (
	// CentralSubject is the identity used in certificates for Central.
	CentralSubject = Subject{ServiceType: storage.ServiceType_CENTRAL_SERVICE, Identifier: "Central"}

	// SensorSubject is the identity used in certificates for Sensor.
	SensorSubject = Subject{ServiceType: storage.ServiceType_SENSOR_SERVICE, Identifier: "Sensor"}

	// ScannerSubject is the identity used in certificates for Scanner.
	ScannerSubject = Subject{ServiceType: storage.ServiceType_SCANNER_SERVICE, Identifier: "Scanner"}

	// ScannerDBSubject is the identity used in certificates for Scanners Postgres DB
	ScannerDBSubject = Subject{ServiceType: storage.ServiceType_SCANNER_DB_SERVICE, Identifier: "Scanner DB"}

	readCACertOnce     sync.Once
	caCert             *x509.Certificate
	caCertDER          []byte
	caCertFileContents []byte
	caCertErr          error

	readCAKeyOnce     sync.Once
	caKeyFileContents []byte
	caKeyErr          error
)

// IssuedCert is a representation of an issued certificate
type IssuedCert struct {
	CertPEM []byte
	KeyPEM  []byte
	ID      *storage.ServiceIdentity
}

// LeafCertificateFromFile reads a tls.Certificate (including private key and cert).
func LeafCertificateFromFile() (tls.Certificate, error) {
	return tls.LoadX509KeyPair(certFilePathSetting.Setting(), keyFilePathSetting.Setting())
}

// ConvertPEMToDERs converts the given certBytes to DER.
// Returns multiple DERs if multiple PEMs were passed.
func ConvertPEMToDERs(certBytes []byte) ([][]byte, error) {
	var result [][]byte

	restBytes := certBytes
	for {
		var decoded *pem.Block
		decoded, restBytes = pem.Decode(restBytes)

		if decoded == nil && len(result) == 0 {
			return nil, errors.New("invalid PEM")
		} else if decoded == nil {
			return result, nil
		}

		result = append(result, decoded.Bytes)
		if len(restBytes) == 0 {
			return result, nil
		}
	}
}

// CACertPEM returns the PEM-encoded CA certificate.
func CACertPEM() ([]byte, error) {
	_, caDER, err := CACert()
	if err != nil {
		return nil, errors.Wrap(err, "CA cert loading")
	}
	return pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caDER,
	}), nil
}

func readCAKey() ([]byte, error) {
	readCAKeyOnce.Do(func() {
		caKeyBytes, err := ioutil.ReadFile(caKeyFilePathSetting.Setting())
		if err != nil {
			caKeyErr = errors.Wrap(err, "reading CA key")
			return
		}
		caKeyFileContents = caKeyBytes
	})
	return caKeyFileContents, caKeyErr
}

func readCA() (*x509.Certificate, []byte, []byte, error) {
	readCACertOnce.Do(func() {
		caBytes, err := ioutil.ReadFile(caFilePathSetting.Setting())
		if err != nil {
			caCertErr = errors.Wrap(err, "reading CA file")
			return
		}

		der, err := ConvertPEMToDERs(caBytes)
		if err != nil {
			caCertErr = errors.Wrap(err, "CA cert could not be decoded")
			return
		}
		if len(der) == 0 {
			caCertErr = errors.New("reading CA file failed")
			return
		}

		cert, err := x509.ParseCertificate(der[0])
		if err != nil {
			caCertErr = errors.Wrap(err, "CA cert could not be parsed")
			return
		}
		caCertFileContents = caBytes
		caCert = cert
		caCertDER = der[0]
	})
	return caCert, caCertFileContents, caCertDER, caCertErr
}

// CACert reads the cert from the local file system and returns the cert and the DER encoding.
func CACert() (*x509.Certificate, []byte, error) {
	caCert, _, caCertDER, caCertErr := readCA()
	return caCert, caCertDER, caCertErr
}

func signerFromCABytes(caCert, caKey []byte) (cfsigner.Signer, error) {
	parsedCa, err := helpers.ParseCertificatePEM(caCert)
	if err != nil {
		return nil, err
	}

	priv, err := helpers.ParsePrivateKeyPEMWithPassword(caKey, nil)
	if err != nil {
		return nil, err
	}

	return local.NewSigner(priv, parsedCa, cfsigner.DefaultSigAlgo(priv), signingPolicy())
}

func signer() (cfsigner.Signer, error) {
	return local.NewSignerFromFile(caFilePathSetting.Setting(), caKeyFilePathSetting.Setting(), signingPolicy())
}

func signingPolicy() *config.Signing {
	return &config.Signing{
		Default: &config.SigningProfile{
			Usage:    []string{"signing", "key encipherment", "server auth", "client auth"},
			Expiry:   certLifetime + beforeGracePeriod,
			Backdate: beforeGracePeriod,
			CSRWhitelist: &config.CSRWhitelist{
				PublicKey:          true,
				PublicKeyAlgorithm: true,
				SignatureAlgorithm: true,
			},
		},
	}
}

// CACertAndKey returns the contents of the ca cert and ca key files.
func CACertAndKey() ([]byte, []byte, error) {
	_, caCertFileContents, _, err := readCA()
	if err != nil {
		return nil, nil, err
	}
	caKeyFileContents, err := readCAKey()
	if err != nil {
		return nil, nil, err
	}
	return caCertFileContents, caKeyFileContents, nil
}

// IssueNewCertFromCA issues a certificate from the CA that is passed in
func IssueNewCertFromCA(subj Subject, caCert, caKey []byte) (cert *IssuedCert, err error) {
	s, err := signerFromCABytes(caCert, caKey)
	if err != nil {
		return nil, errors.Wrap(err, "signer creation")
	}

	return issueNewCertFromSigner(subj, s)
}

func validateSubject(subj Subject) error {
	errorList := errorhelpers.NewErrorList("")
	if subj.ServiceType == storage.ServiceType_UNKNOWN_SERVICE {
		errorList.AddString("Subject service type must be known")
	}
	if subj.Identifier == "" {
		errorList.AddString("Subject Identifier must be non-empty")
	}
	return errorList.ToError()
}

func issueNewCertFromSigner(subj Subject, signer cfsigner.Signer) (*IssuedCert, error) {
	if err := validateSubject(subj); err != nil {
		// Purposefully didn't use returnErr because errorList.ToError() returned from validateSubject is already prefixed
		return nil, err
	}

	serial, err := RandomSerial()
	if err != nil {
		return nil, errors.Wrap(err, "serial generation")
	}
	csr := &cfcsr.CertificateRequest{
		KeyRequest: cfcsr.NewBasicKeyRequest(),
	}
	csrBytes, keyBytes, err := cfcsr.ParseRequest(csr)
	if err != nil {
		return nil, errors.Wrap(err, "request parsing")
	}

	req := cfsigner.SignRequest{
		Hosts:   subj.AllHostnames(),
		Request: string(csrBytes),
		Subject: &cfsigner.Subject{
			CN:           subj.CN(),
			Names:        []cfcsr.Name{subj.Name()},
			SerialNumber: serial.String(),
		},
	}
	certBytes, err := signer.Sign(req)
	if err != nil {
		return nil, errors.Wrap(err, "signing")
	}

	id := generateIdentity(subj, serial)

	return &IssuedCert{
		CertPEM: certBytes,
		KeyPEM:  keyBytes,
		ID:      id,
	}, nil

}

// IssueNewCert generates a new key and certificate chain for a sensor.
func IssueNewCert(subj Subject) (cert *IssuedCert, err error) {
	s, err := signer()
	if err != nil {
		return nil, errors.Wrap(err, "signer creation")
	}
	return issueNewCertFromSigner(subj, s)
}

// RandomSerial returns a new integer that can be used as a certificate serial number (i.e., it is positive and contains
// 64 bits of random data).
func RandomSerial() (*big.Int, error) {
	serial, err := rand.Int(rand.Reader, serialMax)
	if err != nil {
		return nil, errors.Wrap(err, "serial number generation")
	}
	serial.Add(serial, big.NewInt(1)) // Serial numbers must be positive.
	return serial, nil
}

func generateIdentity(subj Subject, serial *big.Int) *storage.ServiceIdentity {
	return &storage.ServiceIdentity{
		Id:   subj.Identifier,
		Type: subj.ServiceType,
		Srl: &storage.ServiceIdentity_SerialStr{
			SerialStr: serial.String(),
		},
	}
}

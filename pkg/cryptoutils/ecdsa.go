package cryptoutils

import (
	"crypto"
	"crypto/ecdsa"
	"encoding/asn1"
	"errors"
	"fmt"
	"io"
	"math/big"
)

// NewECDSAVerifier returns a new verifier using the ECDSA algorithm with the given private key and hash
// function.
func NewECDSAVerifier(publicKey *ecdsa.PublicKey, hash crypto.Hash) SignatureVerifier {
	return &ecdsaVerifier{
		publicKey: publicKey,
		hash:      hash,
	}
}

type ecdsaVerifier struct {
	hash      crypto.Hash
	publicKey *ecdsa.PublicKey
}

type ecdsaSig struct {
	R, S *big.Int
}

func (v *ecdsaVerifier) Verify(data, sig []byte) error {
	digest, err := ComputeDigest(data, v.hash)
	if err != nil {
		return fmt.Errorf("computing digest: %v", err)
	}

	var ecdsaSig ecdsaSig
	if rest, err := asn1.Unmarshal(sig, &ecdsaSig); err != nil {
		return fmt.Errorf("unmarshalling signature: %v", err)
	} else if len(rest) != 0 {
		return fmt.Errorf("unmarshalling signature: %d extra bytes", len(rest))
	}

	if !ecdsa.Verify(v.publicKey, digest, ecdsaSig.R, ecdsaSig.S) {
		return errors.New("signature verification failed")
	}
	return nil
}

// NewECDSASigner returns a new signer using the ECDSA algorithm with the given private key and hash
// function.
func NewECDSASigner(pk *ecdsa.PrivateKey, hash crypto.Hash) Signer {
	return &ecdsaSigner{
		ecdsaVerifier: ecdsaVerifier{
			hash:      hash,
			publicKey: &pk.PublicKey,
		},
		priv: pk,
	}
}

type ecdsaSigner struct {
	ecdsaVerifier
	priv *ecdsa.PrivateKey
}

func (es *ecdsaSigner) Sign(data []byte, entropySrc io.Reader) ([]byte, error) {
	digest, err := ComputeDigest(data, es.hash)
	if err != nil {
		return nil, fmt.Errorf("computing digest: %v", err)
	}

	r, s, err := ecdsa.Sign(entropySrc, es.priv, digest)
	if err != nil {
		return nil, fmt.Errorf("signing: %v", err)
	}
	sig, err := asn1.Marshal(ecdsaSig{R: r, S: s})
	if err != nil {
		return nil, fmt.Errorf("marshalling: %v", err)
	}
	return sig, nil
}

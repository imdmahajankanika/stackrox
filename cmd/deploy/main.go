package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"os"

	"bitbucket.org/stack-rox/apollo/generated/api/v1"
	"bitbucket.org/stack-rox/apollo/pkg/central"
	"bitbucket.org/stack-rox/apollo/pkg/logging"
	zipPkg "bitbucket.org/stack-rox/apollo/pkg/zip"
	"github.com/cloudflare/cfssl/csr"
	"github.com/cloudflare/cfssl/initca"
	"github.com/spf13/cobra"
)

var (
	logger = logging.LoggerForModule()
)

var (
	clairifyTag   = "0.3.1"
	clairifyImage = "clairify:" + clairifyTag
	preventTag    = "1.3"
	preventImage  = "prevent:" + preventTag
)

// ServeHTTP serves a ZIP file for the cluster upon request.
func outputZip(config central.Config) error {
	buf := new(bytes.Buffer)
	zipW := zip.NewWriter(buf)

	d, ok := central.Deployers[config.ClusterType]
	if !ok {
		return fmt.Errorf("Undefined cluster deployment generator: %s", config.ClusterType)
	}

	files, err := d.Render(config)
	if err != nil {
		return fmt.Errorf("Could not render files: %s", err)
	}
	for _, f := range files {
		if err := zipPkg.AddFile(zipW, f); err != nil {
			return fmt.Errorf("Failed to write '%s': %s", f.Name, err)
		}
	}
	// Add MTLS files
	req := csr.CertificateRequest{
		CN:         "StackRox Prevent Certificate Authority",
		KeyRequest: csr.NewBasicKeyRequest(),
	}
	cert, _, key, err := initca.New(&req)
	if err != nil {
		return fmt.Errorf("Could not generate keypair: %s", err)
	}
	if err := zipPkg.AddFile(zipW, zipPkg.NewFile("ca.pem", string(cert), false)); err != nil {
		return fmt.Errorf("Failed to write cert.pem: %s", err)
	}
	if err := zipPkg.AddFile(zipW, zipPkg.NewFile("ca-key.pem", string(key), false)); err != nil {
		return fmt.Errorf("Failed to write key.pem: %s", err)
	}

	err = zipW.Close()
	if err != nil {
		return fmt.Errorf("Couldn't close zip writer: %s", err)
	}

	_, err = os.Stdout.Write(buf.Bytes())
	if err != nil {
		return fmt.Errorf("Couldn't write zip file: %s", err)
	}
	return err
}

func root() *cobra.Command {
	c := &cobra.Command{
		Use:          "root",
		SilenceUsage: true,
	}
	c.AddCommand(interactive())
	c.AddCommand(cmd())
	return c
}

func interactive() *cobra.Command {
	return &cobra.Command{
		Use: "interactive",
		RunE: func(c *cobra.Command, args []string) error {
			c = cmd()
			c.SilenceUsage = true
			return runInteractive(c)
		},
		SilenceUsage: true,
	}
}

func cmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy generates deployment files for StackRox Prevent Central",
		Long: `Deploy generates deployment files for StackRox Prevent Central.
Output is a zip file printed to stdout.`,
		Run: func(*cobra.Command, []string) {
			printToStderr("Orchestrator is required\n")
		},
	}
	c.AddCommand(k8s())
	c.AddCommand(openshift())
	c.AddCommand(dockerBasedOrchestrator("dockeree", "Docker EE", v1.ClusterType_DOCKER_EE_CLUSTER))
	c.AddCommand(dockerBasedOrchestrator("swarm", "Docker Swarm", v1.ClusterType_SWARM_CLUSTER))
	return c
}

func runInteractive(cmd *cobra.Command) error {
	// Overwrite os.Args because cobra uses them
	os.Args = walkTree(cmd)
	return cmd.Execute()
}

func main() {
	root().Execute()
}

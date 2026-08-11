package ca

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"time"

	"github.com/spf13/cobra"
)

const (
	generateFuncName = "generate"
	generateCmdDesc  = "Generate the CA certificate and key"
)

func GenerateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   generateFuncName,
		Short: generateCmdDesc,
		RunE:  generateCa,
	}

	return cmd
}

func generateCa(cmd *cobra.Command, args []string) error {
	// Certificate Authority
	caCert, caPrivKey, err := generateCaPki()
	if err != nil {
		return err
	}

	caBytes, err := signCertificate(caCert, caPrivKey, caCert, caPrivKey)
	if err != nil {
		return err
	}

	caPEM, caPrivKeyPEM, err := encodeCertificateAndKey(caPrivKey, caBytes)
	if err != nil {
		return err
	}

	// Server
	serverCert, serverPrivKey, err := generateClientPki("Server")
	if err != nil {
		return err
	}

	certBytes, err := signCertificate(serverCert, serverPrivKey, caCert, caPrivKey)
	if err != nil {
		return err
	}

	serverCertPEM, serverCertPrivKeyPEM, err := encodeCertificateAndKey(serverPrivKey, certBytes)
	if err != nil {
		return err
	}


	// Client
	clientCert, clientPrivKey, err := generateClientPki("Client")
	if err != nil {
		return err
	}

	clientCertBytes, err := signCertificate(clientCert, clientPrivKey, caCert, caPrivKey)
	if err != nil {
		return err
	}

	clientCertPEM, clientCertPrivKeyPEM, err := encodeCertificateAndKey(clientPrivKey, clientCertBytes)
	if err != nil {
		return err
	}

	err = writePki("ca", caPEM, caPrivKeyPEM)
	if err != nil {
		return err
	}

	err = writePki("server", serverCertPEM, serverCertPrivKeyPEM)
	if err != nil {
		return err
	}

	err = writePki("client", clientCertPEM, clientCertPrivKeyPEM)
	if err != nil {
		return err
	}

	fmt.Println("PKI material generated successfully.")

	return nil
}

func generateCaPki() (*x509.Certificate, *rsa.PrivateKey, error) {
	// Create the cert to act as the CA
	ca := &x509.Certificate{
		SerialNumber: big.NewInt(time.Now().UnixNano()),
		Subject: pkix.Name{
			Organization:  []string{"Certificate Authority"},
			Country:       []string{"UK"},
			Province:      []string{""},
			Locality:      []string{"Halifax"},
			StreetAddress: []string{"1 Main Street"},
			PostalCode:    []string{"HX1 1AA"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	// Generate private key for the CA
	caPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, nil, err
	}

	return ca, caPrivKey, nil
}

func generateClientPki(organisation string) (*x509.Certificate, *rsa.PrivateKey, error) {
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(time.Now().UnixNano()),
		Subject: pkix.Name{
			Organization:  []string{organisation},
			Country:       []string{"UK"},
			Province:      []string{""},
			Locality:      []string{"Halifax"},
			StreetAddress: []string{"1 Main Street"},
			PostalCode:    []string{"HX1 1AA"},
		},
		IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		DNSNames:     []string{"localhost"},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	// Create a private key for the cert
	certPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, nil, err
	}

	return cert, certPrivKey, nil
}

func signCertificate(cert *x509.Certificate, certPrivKey *rsa.PrivateKey, caCert *x509.Certificate, caPrivKey *rsa.PrivateKey) ([]byte, error) {
	certBytes, err := x509.CreateCertificate(rand.Reader, cert, caCert, &certPrivKey.PublicKey, caPrivKey)
	if err != nil {
		return nil, err
	}

	return certBytes, nil
}

func encodeCertificateAndKey(privKey *rsa.PrivateKey, certBytes []byte) (*bytes.Buffer, *bytes.Buffer, error) {
	// Encode the certificate into a PEM
	certPEM := new(bytes.Buffer)
	pem.Encode(certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})

	// Encode the private key into a PEM
	privKeyPEM := new(bytes.Buffer)
	pem.Encode(privKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privKey),
	})

	return certPEM, privKeyPEM, nil
}

func writePki(prefix string, certPEM *bytes.Buffer, privKeyPEM *bytes.Buffer) error {
	// Write the cert and key to files
	if err := os.WriteFile(fmt.Sprintf("pki/%s.crt", prefix), certPEM.Bytes(), 0644); err != nil {
		return err
	}

	if err := os.WriteFile(fmt.Sprintf("pki/%s.key", prefix), privKeyPEM.Bytes(), 0644); err != nil {
		return err
	}

	return nil
}

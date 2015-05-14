package plugins

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"flag"
	"fmt"
	"math/big"
	"net"
	"os"
	"strings"
	"time"

	"github.com/gsdocker/gserrors"
	"github.com/gsdocker/gsmake"
)

// TaskGsweb implement task gsweb
func TaskGsweb(context *gsmake.Runner, args ...string) error {
	context.I("hello gsweb")
	return nil
}

// TaskGencert implement gswebcert task
func TaskGencert(context *gsmake.Runner, args ...string) error {

	flagset := flag.NewFlagSet("gencert", flag.ContinueOnError)

	var (
		host       = flagset.String("host", "", "Comma-separated hostnames and IPs to generate a certificate for")
		validFrom  = flagset.String("start-date", "", "Creation date formatted as Jan 1 15:04:05 2011")
		validFor   = flagset.Duration("duration", 365*24*time.Hour, "Duration that certificate is valid for")
		isCA       = flagset.Bool("ca", false, "whether this cert should be its own Certificate Authority")
		rsaBits    = flagset.Int("rsa-bits", 2048, "Size of RSA key to generate. Ignored if --ecdsa-curve is set")
		ecdsaCurve = flagset.String("ecdsa-curve", "", "ECDSA curve to use to generate a key. Valid values are P224, P256, P384, P521")
	)

	err := flagset.Parse(args)

	if err != nil {
		return gserrors.Newf(err, "invalid gencert task args")
	}

	if len(*host) == 0 {
		return gserrors.Newf(nil, "Missing required --host parameter")
	}

	var priv interface{}
	switch *ecdsaCurve {
	case "":
		priv, err = rsa.GenerateKey(rand.Reader, *rsaBits)
	case "P224":
		priv, err = ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
	case "P256":
		priv, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	case "P384":
		priv, err = ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	case "P521":
		priv, err = ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	default:
		fmt.Fprintf(os.Stderr, "Unrecognized elliptic curve: %q", *ecdsaCurve)
		os.Exit(1)
	}
	if err != nil {
		return gserrors.Newf(err, "failed to generate private key")
	}

	var notBefore time.Time
	if len(*validFrom) == 0 {
		notBefore = time.Now()
	} else {
		notBefore, err = time.Parse("Jan 2 15:04:05 2006", *validFrom)
		if err != nil {
			return gserrors.Newf(err, "Failed to parse creation date:")
		}
	}

	notAfter := notBefore.Add(*validFor)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return gserrors.Newf(err, "failed to generate serial number")
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"gsweb Co"},
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	hosts := strings.Split(*host, ",")
	for _, h := range hosts {
		if ip := net.ParseIP(h); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, h)
		}
	}

	if *isCA {
		template.IsCA = true
		template.KeyUsage |= x509.KeyUsageCertSign
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, publicKey(priv), priv)
	if err != nil {
		return gserrors.Newf(err, "Failed to create certificate")
	}

	certOut, err := os.Create("cert.pem")
	if err != nil {
		return gserrors.Newf(err, "failed to open cert.pem for writing")
	}
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	certOut.Close()
	context.I("written cert.pem")

	keyOut, err := os.OpenFile("key.pem", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return gserrors.Newf(err, "failed to open key.pem for writing")
	}
	pem.Encode(keyOut, pemBlockForKey(priv))
	keyOut.Close()
	context.I("written key.pem\n")

	return nil
}
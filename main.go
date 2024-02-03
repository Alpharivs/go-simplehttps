package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"flag"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)
// Arguments
var (
	dir    = flag.String("d", ".", "The directory to serve files from. (Default: current dir)")
	secure = flag.Bool("s", false, "Use TLS.")
	port   = flag.String("p", "", "Listening port. (Default 80 or 443 is using TLS")
)
// Not a graceful Server Shutdown, may improve later.
func shutdown() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\r[!] Shutting down server...")
		fmt.Printf("\nLVX SIT - ALPHARIVS - MMXXI")
		os.Exit(1)
	}()
}
// Generate the .crt and .key file
func GenKeyAndCert() ([]byte, []byte, error) {
	bits := 4096
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, errors.Wrap(err, "rsa.GenerateKey")
	}

	template := x509.Certificate{
		SerialNumber:          big.NewInt(0),
		Subject:               pkix.Name{CommonName: "https-go"},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(1, 0, 0),
		BasicConstraintsValid: true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
	}

	derCert, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return nil, nil, errors.Wrap(err, "x509.CreateCertificate")
	}

	buf := &bytes.Buffer{}
	err = pem.Encode(buf, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: derCert,
	})
	if err != nil {
		return nil, nil, errors.Wrap(err, "pem.Encode")
	}

	pemCert := buf.Bytes()

	buf = &bytes.Buffer{}
	err = pem.Encode(buf, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})
	if err != nil {
		return nil, nil, errors.Wrap(err, "pem.Encode")
	}

	pemKey := buf.Bytes()

	return pemCert, pemKey, nil
}

// Configure and run Https server.
func HttpsServer(port string, handler http.Handler) error {
	rawCert, rawKey, err := GenKeyAndCert()
	if err != nil {
		return errors.Wrap(err, "GenKeyAndCert")
	}

	cert, err := tls.X509KeyPair(rawCert, rawKey)
	if err != nil {
		return errors.Wrap(err, "x509KeyPair")
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	server := http.Server{
		Addr:      ":" + port,
		Handler:   handler,
		TLSConfig: tlsConfig,
	}

	fmt.Printf("[!] Started Https server on port %s\n", port)
	return server.ListenAndServeTLS("", "")
}

func main() {
	flag.Parse()
	shutdown()

	r := mux.NewRouter()
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(*dir)))
	l := handlers.LoggingHandler(os.Stdout, r)

	var err error

	switch {
	case *secure:
		if *port == "" {
			err = HttpsServer("443", l)
		} else {
			err = HttpsServer(*port, l)
		}
	case !*secure:
		if *port == "" {
			fmt.Printf("[!] Started Http server on port 80\n")
			err = http.ListenAndServe(":80", l)
		} else {
			fmt.Printf("[!] Started Http server on port %s\n", *port)
			err = http.ListenAndServe(":"+*port, l)
		}
	}
	if err != nil {
		log.Fatal(err)
	}
}

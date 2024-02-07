package main

import (
	"bytes"
	"context"
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
)

// Arguments
var (
	dir    = flag.String("d", ".", "The directory to serve files from. (Default: current dir)")
	secure = flag.Bool("s", false, "Use TLS.")
	port   = flag.String("p", "", "Listening port. (Default 80 or 443 is using TLS")
)

// Graceful Server Shutdown.
func shutdown(server *http.Server) {
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		// create context for server shutdown
		shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownRelease()
		// clear terminal and print banner
		fmt.Print("\r\033[K")
		fmt.Println("\r\nLVX SIT - ALPHARIVS - MMXXI")
		// graceful shutdown
		server.Shutdown(shutdownCtx)
	}()
}

// Generate the .crt and .key file
func GenKeyAndCert() ([]byte, []byte, error) {
	bits := 4096
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, fmt.Errorf("error in rsa.GenerateKey(): %w", err)
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
		return nil, nil, fmt.Errorf("error in x509.CreateCertificate(): %w", err)
	}

	buf := &bytes.Buffer{}
	err = pem.Encode(buf, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: derCert,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("error in cert pem.Encode(): %w", err)
	}

	pemCert := buf.Bytes()

	buf = &bytes.Buffer{}
	err = pem.Encode(buf, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})
	if err != nil {
		return nil, nil, fmt.Errorf("error in key pem.Encode(): %w", err)
	}

	pemKey := buf.Bytes()

	return pemCert, pemKey, nil
}

// Configure and run HTTPS server.
func HttpsServer(port string, handler http.Handler) error {
	rawCert, rawKey, err := GenKeyAndCert()
	if err != nil {
		return fmt.Errorf("error in GenKeyAndCert(): %w", err)
	}

	cert, err := tls.X509KeyPair(rawCert, rawKey)
	if err != nil {
		return fmt.Errorf("error in x509KeyPair(): %w", err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	if port == "" {
		port = "443"
	}

	server := http.Server{
		Addr:      ":" + port,
		Handler:   handler,
		TLSConfig: tlsConfig,
	}
	shutdown(&server)

	fmt.Printf("[!] Started Https server on port %s\n", port)
	return server.ListenAndServeTLS("", "")
}

// Configure and run HTTP server.
func HttpServer(port string, handler http.Handler) error {
	if port == "" {
		port = "80"
	}

	server := http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}
	shutdown(&server)

	fmt.Printf("[!] Started Http server on port %s\n", port)
	return server.ListenAndServe()
}

func main() {
	flag.Parse()
	// Configure routing and middleware.
	r := mux.NewRouter()
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(*dir)))
	l := handlers.LoggingHandler(os.Stdout, r)

	var err error
	// Start HTTPS or HTTP server according to user flag
	switch {
	case *secure:
		err = HttpsServer(*port, l)
	case !*secure:
		err = HttpServer(*port, l)
	}

	if err != nil {
		log.Fatal(err)
	}
}

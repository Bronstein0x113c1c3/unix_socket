package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"math/big"
	"net"
	"net/http"

	"github.com/quic-go/quic-go/http3"
)

func dummy_mux() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "<html><body>hello</body></html>")
	})
	mux.HandleFunc("/demo/tile", func(w http.ResponseWriter, r *http.Request) {
		// Small 40x40 png
		w.Write([]byte{
			0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a, 0x00, 0x00, 0x00, 0x0d,
			0x49, 0x48, 0x44, 0x52, 0x00, 0x00, 0x00, 0x28, 0x00, 0x00, 0x00, 0x28,
			0x01, 0x03, 0x00, 0x00, 0x00, 0xb6, 0x30, 0x2a, 0x2e, 0x00, 0x00, 0x00,
			0x03, 0x50, 0x4c, 0x54, 0x45, 0x5a, 0xc3, 0x5a, 0xad, 0x38, 0xaa, 0xdb,
			0x00, 0x00, 0x00, 0x0b, 0x49, 0x44, 0x41, 0x54, 0x78, 0x01, 0x63, 0x18,
			0x61, 0x00, 0x00, 0x00, 0xf0, 0x00, 0x01, 0xe2, 0xb8, 0x75, 0x22, 0x00,
			0x00, 0x00, 0x00, 0x49, 0x45, 0x4e, 0x44, 0xae, 0x42, 0x60, 0x82,
		})
	})

	mux.HandleFunc("/demo/tiles", func(w http.ResponseWriter, r *http.Request) {
		log.Println("received!!!")
		io.WriteString(w, "<html><head><style>img{width:40px;height:40px;}</style></head><body>")
		for i := 0; i < 200; i++ {
			fmt.Fprintf(w, `<img src="/demo/tile?cachebust=%d">`, i)
		}
		io.WriteString(w, "</body></html>")
	})

	mux.HandleFunc("/demo/echo", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Printf("error reading body while handling /echo: %s\n", err.Error())
		}
		w.Write(body)
	})
	return mux
}
func main() {
	mux := dummy_mux()
	// TLS configuration for HTTP/3
	// tlsConfig := &tls.Config{
	// 	// Certificates: getCertificates(), // Load TLS certificates
	// 	NextProtos:         []string{"h3"}, // ALPN for HTTP/3
	// 	InsecureSkipVerify: true,
	// }
	tls_config, _ := GenerateTLSConfig("h3")

	// Create an HTTP/3 server
	server := http3.Server{
		// Addr:      "../tmp/d3.sock", // Listen on port 443
		TLSConfig: tls_config,
		Handler:   mux,
	}

	// Start the server
	log.Println("Starting HTTP/3 server on Unix domain socket")
	// if err := server.ListenAndServe(); err != nil {
	// 	log.Fatalf("Failed to start HTTP/3 server: %v", err)
	// }
	lis, _ := net.ListenUnixgram("unixgram", &net.UnixAddr{Name: "../tmp/d3.sock", Net: "unixgram"})
	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to start HTTP/3 server: %v", err)
	}
}

// Helper function to load TLS certificates
func getCertificates() []tls.Certificate {
	cert, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		log.Fatalf("Failed to load TLS certificates: %v", err)
	}
	return []tls.Certificate{cert}
}
func GenerateTLSConfig(p string) (*tls.Config, error) {
	// Generate a new RSA key
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Printf("failed to generate RSA key: %s", err)
		return nil, err
	}

	// Create a certificate template
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		// Set other fields of the certificate as required
	}

	// Create a certificate using the template
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		log.Printf("failed to create certificate: %s", err)
		return nil, err
	}

	// Encode the certificate and key to PEM format
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	// Load the X509 key pair from PEM blocks
	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		log.Printf("failed to load X509 key pair from PEM: %s", err)
		return nil, err
	}
	protos := []string{}
	protos = append(protos, p)
	return &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		NextProtos:   protos,
	}, nil
}

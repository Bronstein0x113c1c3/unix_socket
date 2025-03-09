package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
)

func main() {
	// Create an HTTP/3 client
	client := &http.Client{
		Transport: &http3.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // Skip TLS verification for testing
				NextProtos:         []string{"h3"},
			},
			Dial: func(_ context.Context, _ string, _ *tls.Config, _ *quic.Config) (quic.EarlyConnection, error) {
				conn, err := net.ListenPacket("unixgram", "")
				if err != nil {
					return nil, err
				}

				// Dial QUIC over the Unix domain socket
				quicConn, err := quic.DialEarly(context.Background(), conn, &net.UnixAddr{Name: "../tmp/d3.sock", Net: "unixgram"}, &tls.Config{
					InsecureSkipVerify: true,
					NextProtos:         []string{"h3"},
				}, &quic.Config{})
				return quicConn, err
			},
		},
	}

	// Send a GET request to the server
	resp, err := client.Get("https://unix")
	if err != nil {
		log.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	// Read and print the response
	body := make([]byte, 1024)
	n, _ := resp.Body.Read(body)
	fmt.Println("Response:", string(body[:n]))
}

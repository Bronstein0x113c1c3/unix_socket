// package main

// import (
// 	"context"
// 	"crypto/tls"
// 	"fmt"
// 	"log"
// 	"net"
// 	"net/http"

// 	"github.com/quic-go/quic-go"
// 	"github.com/quic-go/quic-go/http3"
// )

// func main() {
// 	client_transport := &http3.Transport{
// 		TLSClientConfig: &tls.Config{
// 			InsecureSkipVerify: true,
// 			NextProtos:         []string{"h3-29"},
// 		},
// 		Dial: func(_ context.Context, _ string, _ *tls.Config, _ *quic.Config) (quic.EarlyConnection, error) {
// 			conn, err := net.DialUnix("unixgram", nil, &net.UnixAddr{Name: "../tmp/d2.sock", Net: "????"})
// 			log.Printf("err at conn: %v \n", err)
// 			if err != nil {
// 				return nil, err
// 			}
// 			tls_config := &tls.Config{
// 				InsecureSkipVerify: true,
// 				NextProtos:         []string{"h3-29"},
// 			}
// 			net, err := quic.DialEarly(context.Background(), conn, nil, tls_config, &quic.Config{})
// 			log.Printf("err at dial: %v \n", err)
// 			return net, err
// 		},
// 		// quic.Connection
// 	}
// 	client := http.Client{
// 		Transport: client_transport}
// 	_, err := client.Get("https://unix/")
// 	fmt.Println(err)
// }

// // func setup_client_conn() {

// // }
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
	// Path to the Unix domain socket
	socketPath := "../tmp/d2.sock"

	// Create a transport for HTTP/3 over Unix domain socket
	clientTransport := &http3.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // Skip TLS verification for testing
			NextProtos:         []string{"something"},
		},
		Dial: func(_ context.Context, _ string, _ *tls.Config, _ *quic.Config) (quic.EarlyConnection, error) {
			// Create an unconnected Unix domain socket
			conn, err := net.ListenPacket("unixgram", "")
			if err != nil {
				return nil, err
			}
			quic.Dial()
			// Dial QUIC over the Unix domain socket
			quicConn, err := quic.DialEarly(context.Background(), conn, &net.UnixAddr{Name: socketPath, Net: "unixgram"}, &tls.Config{
				InsecureSkipVerify: true,
				NextProtos:         []string{"something"},
			}, &quic.Config{})
			return quicConn, err
		},
	}

	// Create an HTTP client
	client := http.Client{
		Transport: clientTransport,
	}

	// Send a GET request to the server
	resp, err := client.Get("https://unix/") // The scheme doesn't matter for Unix sockets
	if err != nil {
		log.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	// Read and print the response
	body := make([]byte, 1024)
	n, _ := resp.Body.Read(body)
	fmt.Println("Response:", string(body[:n]))
}

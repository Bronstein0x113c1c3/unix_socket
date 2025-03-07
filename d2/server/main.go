package main

//http3 through unix socket plan
import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"

	"github.com/quic-go/quic-go/http3"
)

func dummy_mux() http.Handler {
	mux := http.NewServeMux()
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
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	tls, _ := setup_tls()
	mux := dummy_mux()
	server := http3.Server{
		TLSConfig: tls,
		Handler:   mux,
	}
	conn, err := conn_setup("../tmp/d2.sock")
	if err != nil {
		fmt.Printf("Error setting up connection: %s\n", err.Error())
		return
	}
	log.Println("?????")
	go func() {
		server.Serve(conn)
		log.Println("shutdown called...")
		conn.Close()
	}()
	<-ctx.Done()
	server.Shutdown(ctx)
	// defer stop()

}
func conn_setup(socket_path string) (net.PacketConn, error) {
	return net.ListenUnixgram("unixgram", &net.UnixAddr{Name: socket_path, Net: "????"})
}
func setup_tls() (*tls.Config, error) {
	return &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"h3-29"},
	}, nil
}

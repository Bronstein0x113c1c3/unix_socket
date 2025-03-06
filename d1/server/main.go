package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
)

// add graceful shutdown with ctrl+c

// must close after use!!!
func main() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request received from %v \n", r.RemoteAddr)
		w.Write([]byte("Say hello to my little friend!!!!!"))
	})

	lis, err := net.ListenUnix("unix", &net.UnixAddr{Name: "../tmp/d1.sock", Net: "host"})
	defer lis.Close()
	if err != nil {
		log.Printf("Something wrong!!!: %v \n", err)
	}
	go func() {
		log.Printf("Server is running on %s\n", lis.Addr().String())
		http.Serve(lis, mux)
		log.Println("Server is closed")

	}()
	<-sig
	log.Println("Shutting down the server")
}

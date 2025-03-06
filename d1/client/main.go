package main

import (
	"context"
	"io/ioutil"
	"log"
	"net"
	"net/http"
)

func main() {
	client := clientsetup("unix", "../tmp/d1.sock")
	resp, err := client.Get("http://unix/hello")
	// defer resp.Body.Close()
	if err != nil {
		log.Printf("Something wrong!!!: %v \n", err)
		return
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Something wrong!!!: %v \n", err)
		return
	}
	log.Printf("Response: %s\n", string(data))
}

func clientsetup(network, addr string) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial(network, addr)
			},
		},
	}
}

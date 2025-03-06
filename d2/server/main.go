package main

//http3 through unix socket plan
import (
	"net"

	"github.com/quic-go/quic-go/http3"
)

func main() {
	server := http3.Server{}

}
func conn_setup(socket_path string) (net.PacketConn, error) {
	return net.ListenPacket("unixgram", socket_path)
}

package shared

import (
	"log"
	"net"
)

var localIPAddress net.IP

func init() {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Println("Outbound IP error: ", err)
		panic(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	localIPAddress = localAddr.IP
}

func GetOutboundIP() net.IP {
	return localIPAddress
}

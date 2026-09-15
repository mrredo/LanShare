package main

import (
	"net"

	"github.com/gin-gonic/gin"
	"github.com/pion/mdns/v2"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

func main() {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.String(200, "Hello")
	})

	go func() {
		if err := router.Run(":80"); err != nil {
			panic(err)
		}
	}()

	addr4, err := net.ResolveUDPAddr("udp4", mdns.DefaultAddressIPv4)
	if err != nil {
		panic(err)
	}

	addr6, err := net.ResolveUDPAddr("udp6", mdns.DefaultAddressIPv6)
	if err != nil {
		panic(err)
	}

	l4, err := net.ListenUDP("udp4", addr4)
	if err != nil {
		panic(err)
	}

	l6, err := net.ListenUDP("udp6", addr6)
	if err != nil {
		panic(err)
	}

	_, err = mdns.NewServer(
		ipv4.NewPacketConn(l4),
		ipv6.NewPacketConn(l6),
		mdns.WithLocalNames("lanshare.local"),
	)
	if err != nil {
		panic(err)
	}

	select {}
}

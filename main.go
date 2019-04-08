package main

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

func main() {

	// get the network interface
	eth0, err := net.InterfaceByName("eth0")
	if err != nil {
		panic(err)
	}

	laddr, err := getLocalAddressFromInterface(eth0)
	if err != nil {
		panic(err)
	}

	multicast, err := getMulticastAddrFromInterface(eth0)
	if err != nil {
		panic(err)
	}

	multicastPingEveryTwoSeconds(laddr, multicast)

	listen(eth0, laddr, multicast)
}

func multicastPingEveryTwoSeconds(laddr *net.UDPAddr, mcaddr net.Addr) {
	go func() {
		// resolve the remote address from multicast address to prepare for multicast
		raddr, err := net.ResolveUDPAddr("udp", mcaddr.String()+":9999")
		if err != nil {
			panic(err)
		}
		for {
			conn, err := net.DialUDP("udp", laddr, raddr)
			if err != nil {
				log.Println(err)
			}
			conn.Write([]byte(fmt.Sprintf("hello from %s", laddr)))
			conn.Close()
			log.Println(fmt.Sprintf("multicasted - %s. sleeping for 2 seconds", raddr))
			time.Sleep(time.Second * 2)
		}
	}()
}

func listen(netIf *net.Interface, laddr *net.UDPAddr, mcaddr net.Addr) {
	// resolve the group address
	gaddr, err := net.ResolveUDPAddr("udp", mcaddr.String()+":9999")
	if err != nil {
		panic(err)
	}
	for {
		conn, err := net.ListenMulticastUDP("udp", netIf, gaddr)
		if err != nil {
			log.Println("error", err)
			continue
		}
		buf := make([]byte, 1024)
		_, incAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Println("error", err)
			continue
		}
		if incAddr.IP.Equal(laddr.IP) {
			continue
		}
		log.Println(fmt.Sprintf("%s: %s", incAddr.IP, buf))
		conn.Close()
	}
}

func getLocalAddressFromInterface(netIf *net.Interface) (laddr *net.UDPAddr, err error) {
	laddrs, err := netIf.Addrs()
	if err != nil {
		return
	}
	if len(laddrs) == 0 {
		err = fmt.Errorf("no available local address")
		return
	}
	laddr, err = net.ResolveUDPAddr("udp", strings.Split(laddrs[0].String(), "/")[0]+":8000")
	return
}

func getMulticastAddrFromInterface(netIf *net.Interface) (mcaddr net.Addr, err error) {
	multicasts, err := netIf.MulticastAddrs()
	if err != nil {
		return
	}
	if len(multicasts) == 0 {
		err = fmt.Errorf("no multicast address available")
	}
	return multicasts[0], nil
}

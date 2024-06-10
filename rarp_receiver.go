package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"syscall"
)

func main() {
	verbose := flag.Bool("v", false, "Enable verbose output")
	flag.Parse()

	if len(flag.Args()) == 0 {
		log.Fatal("You must specify a network interface.")
	}

	ifaceName := flag.Arg(0)

	iface, err := net.InterfaceByName(ifaceName)
	if err != nil {
		log.Fatalf("Error getting interface: %v", err)
	}

	// create a raw socket for RARP ethertype
	fd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, int(htons(syscall.ETH_P_RARP)))
	if err != nil {
		log.Fatalf("Error creating socket: %v", err)
	}
	defer syscall.Close(fd)

	// bind it to the interface
	addr := syscall.SockaddrLinklayer{
		Ifindex:  iface.Index,
		Protocol: htons(syscall.ETH_P_RARP),
	}

	if err := syscall.Bind(fd, &addr); err != nil {
		log.Fatalf("Error binding socket: %v", err)
	}

	// no jumbos expected
	buffer := make([]byte, 1500)

	for {
		n, _, err := syscall.Recvfrom(fd, buffer, 0)
		if err != nil {
			log.Fatalf("Error reading packet: %v", err)
		}

		// print packet contents as HEX
		if *verbose {
			fmt.Printf("Received packet:\n")
			fmt.Printf("%x\n", buffer[:n])
		}
	}
}

// host byte order to network byte order. Little Endian to Big Endian
func htons(i uint16) uint16 {
	return (i<<8)&0xff00 | i>>8
}

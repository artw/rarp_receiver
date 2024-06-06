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

	fd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, int(htons(syscall.ETH_P_RARP)))
	if err != nil {
		log.Fatalf("Error creating socket: %v", err)
	}
	defer syscall.Close(fd)

	addr := syscall.SockaddrLinklayer{
		Ifindex:  iface.Index,
		Protocol: htons(syscall.ETH_P_RARP),
	}

	if err := syscall.Bind(fd, &addr); err != nil {
		log.Fatalf("Error binding socket: %v", err)
	}

	buffer := make([]byte, 1500)

	for {
		n, _, err := syscall.Recvfrom(fd, buffer, 0)
		if err != nil {
			log.Fatalf("Error reading packet: %v", err)
		}

		if *verbose {
			fmt.Printf("Received packet:\n")
			fmt.Printf("%x\n", buffer[:n])
		}
	}
}

func htons(i uint16) uint16 {
	return (i<<8)&0xff00 | i>>8
}

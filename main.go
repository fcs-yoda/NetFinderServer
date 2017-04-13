package main

import (
	"log"
	"net"
)

const UDP_PORT_NO = "12000"
const TCP_PORT_NO = "12001"

func main() {

	go connectionWait()

	listener, err := net.Listen("tcp", "localhost:"+TCP_PORT_NO)

	if err != nil {                                                                
		log.Fatalln(err)
	}    

	log.Print("Sever start")

	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Printf("Accept Error: %v\n", err)
			continue
		}

		log.Printf("Accept [%v]\n", conn.RemoteAddr())

		go doTcpConnection(conn)
	}
}

func doTcpConnection(conn net.Conn) {

	defer conn.Close()

}

func connectionWait() {

	addr, err := net.ResolveUDPAddr("udp", ":"+UDP_PORT_NO)

	if err != nil {
		log.Fatalln(err)
	}

	conn, err := net.ListenUDP("udp", addr)

	if err != nil {
		log.Fatalln(err)
	}

	defer conn.Close()

	buf := make([]byte, 1024)

	log.Print("Wait Conection...")

	for {
		rlen, remote, err := conn.ReadFromUDP(buf)

		if err != nil {
			log.Fatalf("Error: %v\n", err)
		}

		s := string(buf[:rlen])

		log.Printf("Receive [%v]: %v\n", remote, s)

		if s == "CON_REQ" {

			myIp := getIP()

			rlen, err = conn.WriteToUDP([]byte(myIp), remote)

			if err != nil {
				log.Printf("Receive Error [%v]: %v\n", remote, s)
			}

			log.Printf("Send [%v]: %v\n", remote, myIp)
		}
	}

}

func getIP() string {

	addrs, err := net.InterfaceAddrs()

	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ip := ipnet.IP.To4(); ip != nil {
				return a.String()
			}
		}
	}
	
	return ""
}

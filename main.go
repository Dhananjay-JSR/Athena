package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	//	Step 1 :- Create a Intermediate Server That Initialize
	//	Connection between Client and Actual Resource
	//	Step 2 :- When any External Client Request to Server
	//	The Server Creates a Connection to Actual Server and Process
	//	The Actual Resource is Behind a Firewall
	//	3rd party client cannot make actual call to Resource

	//	Server

	// 6000 Running Webserber
	// 80 is Server Running
	// 3000 Client Running

	// Chrome Access -> 3400 -> Server Checking if there's a Client on Port 3400
	// means chrome is trying to access tcp connection from 3rd party hosted on port 3400
	// Server Check if there's a Client on 3400
	// Server Sets the Payload to Client 3400
	// Client Response Back with that ID

	isServer := true
	if isServer {
		con, err := net.Listen("tcp", "127.0.0.1:2002")
		if err != nil {
			log.Fatal(err)
		}
		conn, err := con.Accept()

		if err != nil {
			log.Fatal(err)
		}
		byteWritten, err := conn.Write([]byte("Hello World"))
		if err != nil {
			fmt.Println(err)
		}

		fmt.Printf("The Bytes Written Are %d", byteWritten)
	} else {
		_, err := net.Dial("tcp", "127.0.0.1:2002")

		if err != nil {
			log.Fatal(err)
		}

		//bytesData := make([]byte, 1024)
		//for {
		//	n, err := conn.Read(bytesData)
		//	if err != nil {
		//		fmt.Println("Error:", err)
		//		return
		//	}
		//	fmt.Println(string(bytesData[:n]))
		//}
	}

}

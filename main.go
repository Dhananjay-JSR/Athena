package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
)

func main() {

	connectType := flag.String("type", "server", "defines the type application is started")
	flag.Parse()
	ToClientChan := make(chan []byte)
	FromClientChan := make(chan []byte)
	InterProcessChan := make(chan []byte)

	if *connectType == "client" || *connectType == "NULL" {

		ConnectionedClient := make(map[string]chan []byte)

		ServerDialer, ServerDialErr := net.Dial("tcp", "127.0.0.1:2001")
		if ServerDialErr != nil {
			log.Println(ServerDialErr)
		}

		var wg sync.WaitGroup

		wg.Add(1)

		go func() {
			for {
				readBuffer := make([]byte, 4098)
				readCount, ReadErr := ServerDialer.Read(readBuffer)
				if ReadErr != nil {
					log.Println(ReadErr)
				}
				fmt.Println("Data Received from Middle Server")
				InterProcessChan <- readBuffer[:readCount]
			}
		}()

		go func() {
			for {
				RecvData := <-FromClientChan
				write, err := ServerDialer.Write(RecvData)
				if err != nil {
					return
				}
				fmt.Println("-- Data Sent to Middle Server -- \n", string(RecvData), "With Write Count", write)
			}
		}()

		go func() {

			for {
				RecvData := <-InterProcessChan
				identifier, Data := DecodeData(RecvData)
				//fmt.Println("Incoming Connection", identifier, "With Data", string(Data))
				AssignedChann, ConExist := ConnectionedClient[identifier]
				if ConExist {
					fmt.Println("Identifier", identifier, "Exists")
					AssignedChann <- Data
				} else {
					fmt.Println("-- Connection Identiefier ", identifier, "Does Not Exists Before")
					//fmt.Println("Creating New Connection")
					dialerCon, dialerr := net.Dial("tcp", "127.0.0.1:27017")
					AllocateChan := make(chan []byte)
					ConnectionedClient[identifier] = AllocateChan
					println("Sending Data to Channel Allocated to ", identifier)
					go func(Identifier string, DialerCon net.Conn, AllocatedChannel chan []byte) {
						for {
							RecvData := <-AllocatedChannel
							fmt.Println("-- Received Data from Channel Allocated to ", Identifier, " --- Data Content -- \n", string(RecvData))
							fmt.Println("Writing Data to Actual Server", Identifier)
							_, WriteErr := DialerCon.Write(RecvData)
							if WriteErr != nil {
								log.Println(WriteErr)
							}
						}
					}(identifier, dialerCon, AllocateChan)

					go func(Identifier string, DialerCon net.Conn, FromClientChan chan []byte) {
						for {
							readBuffer := make([]byte, 4098)
							readCount, ReadErr := DialerCon.Read(readBuffer)
							if ReadErr != nil {
								log.Println(ReadErr)
							}
							fmt.Println("-- Received Data from Actual Server ", Identifier, " --- Data Content -- \n", string(readBuffer[:readCount]))
							FromClientChan <- EncodeData(Identifier, readBuffer[:readCount])
						}
					}(identifier, dialerCon, FromClientChan)
					//FUCKK !! Go Routine Blocks untill There is Something to Read
					AllocateChan <- Data

					if dialerr != nil {
						log.Println(dialerr)
					}

				}
			}
		}()

		wg.Wait()

	} else {

		// MY ISSUE
		// Each Client Was Getting Resource From Same Connection
		// There a Single FLow of Data , Not Continuous

		// TODO: List\
		// MAIN :- BIDUPLEX CONTINUE DATA
		//  :- EACH CLIENT TO NEW SERVER CONNECTION
		// 1 For Server Implementation
		// 2 as Soon as Any Request Comes in , Generate a Hash from the Request
		// and Send it a a GRoutine
		// and in GRoutine , Send it to Server
		//  like FORMAT {{{{HASH}}}}{Data} the Hash is used to identify the client connected to Server
		//  on Client Side , the Hash is used to identify the client connected to Server
		//  After the Hash is Identified , the Data is sent to the Client
		BindCon, BindErr := net.Listen("tcp", "127.0.0.1:2001")
		if BindErr != nil {
			log.Println(BindErr)
		}
		// var AthenaConnection net.Conn
		isAthenaConneted := false
		for {

			if !isAthenaConneted {
				fmt.Println("Waiting for Athena")
			} else {
				fmt.Println("Waiting for Someone")
			}
			AcceptCon, AcceptErr := BindCon.Accept()
			if AcceptErr != nil {

				log.Println(AcceptErr)
			}
			if !isAthenaConneted {
				go func(AthenaCon net.Conn) {
					for {
						RecvData := <-ToClientChan
						//println("Data Sent to Athena")
						_, WriteErr := AthenaCon.Write(RecvData)
						//fmt.Println("Data Sent To Athena ", string(RecvData))
						if WriteErr != nil {
							log.Println(WriteErr)
						}

					}
				}(AcceptCon)

				go func(AthenaCon net.Conn) {
					for {
						readBuffer := make([]byte, 4098)
						readCount, ReadErr := AcceptCon.Read(readBuffer)
						//fmt.Println("Data Received from Athena ", string(readBuffer[:readCount]))
						if ReadErr != nil {
							log.Println(ReadErr)
						}
						FromClientChan <- readBuffer[:readCount]
					}
				}(AcceptCon)
				isAthenaConneted = true
				fmt.Println("Athena Connected")
			} else {

				fmt.Println("Someon Connected")

				ConnIdentifier := AcceptCon.RemoteAddr().String()

				go func(ConnectingClient net.Conn, Itenifier string) {
					for {
						RecvResponse := <-FromClientChan
						identifier, Data := DecodeData(RecvResponse)
						if identifier == Itenifier {
							fmt.Println("-- Sending Data back to Connected Client ID \n", (identifier), "With Data Content -- \n", string(Data))
							_, WriteErr := ConnectingClient.Write(Data)
							if WriteErr != nil {
								log.Println(WriteErr)
							}
						}

					}
				}(AcceptCon, ConnIdentifier)

				go func(ConnectingClient net.Conn, Itenifier string) {
					for {
						ReadBuffer := make([]byte, 4098)
						ReadCount, ReadErr := ConnectingClient.Read(ReadBuffer)
						if ReadErr != nil {
							log.Println(ReadErr)
						}
						fmt.Println("-- Received This Data Coming From Connected Client ID \n", Itenifier, " With Data Content -- \n", string(ReadBuffer[:ReadCount]))
						ToClientChan <- EncodeData(Itenifier, ReadBuffer[:ReadCount])
					}

				}(AcceptCon, ConnIdentifier)
			}
		}
	}
}

func asasa() {
	listenner, listenErr := net.Listen("tcp", "127.0.0.1:2001")
	if listenErr != nil {
		return

	}
	for {
		acceptedCon, acceperr := listenner.Accept()
		if acceperr != nil {
			return

		}

		go func(acceptedCon net.Conn) {
			fmt.Println("New Connection Accepted ", acceptedCon.RemoteAddr())

			dialer, dialerrr := net.Dial("tcp", "127.0.0.1:27017")
			if dialerrr != nil {
				return
			}

			readByte := make([]byte, 4098)
			readCount, _ := acceptedCon.Read(readByte)

			fmt.Println(acceptedCon.RemoteAddr(), " Data Received from Client \n", string(readByte[:readCount]))
			dialer.Write(readByte[:readCount])
			readCount, _ = dialer.Read(readByte)
			fmt.Println(acceptedCon.RemoteAddr(), " Data Received from Server \n", string(readByte[:readCount]))
			acceptedCon.Write(readByte[:readCount])
			readCount, _ = acceptedCon.Read(readByte)
			fmt.Println(acceptedCon.RemoteAddr(), " Data Received from Client \n", string(readByte[:readCount]))
			dialer.Write(readByte[:readCount])
			readCount, _ = dialer.Read(readByte)
			fmt.Println(acceptedCon.RemoteAddr(), " Data Received from Server \n", string(readByte[:readCount]))
			acceptedCon.Write(readByte[:readCount])

		}(acceptedCon)

	}

}

const DELIMITER = "|||||"

func EncodeData(UniqueIdentifier string, Data []byte) []byte {
	// Encode the Data with the UniqueIdentifier
	// Return the Encoded Data
	return append([]byte(UniqueIdentifier+DELIMITER), Data...)
}

func DecodeData(Data []byte) (string, []byte) {

	// Decode the Data
	// Return the UniqueIdentifier and the Data
	Identifier := string(Data[:strings.Index(string(Data), DELIMITER)])
	return Identifier, Data[strings.Index(string(Data), DELIMITER)+5:]
}

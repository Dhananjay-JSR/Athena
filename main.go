package main

import (
	"flag"
	"fmt"
	logManager "github.com/Dhananjay-JSR/Athena.git/cli"
	"github.com/Dhananjay-JSR/Athena.git/internal"
	"net"
	"strconv"
	"strings"
	"sync"
)

const secret_key = "Athena"

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

	//EXACT IMPLEMENTATION
	//BROWSER -> SERVER -> CLIENT -> Actual Resource
	internal.InitWindowsEscape()

	// flags
	//1) default secret string used for encryption
	//2) type , server , ip , forwarder
	//3) server-range
	//4) localhost range
	//5) help
	// 6) server url -> url client will connect to

	secretFlag := flag.String("secret", secret_key, "Secret key is used to provide an encryption layer over client server communication to prevent anyone from hijacking connection")
	connectType := flag.String("type", "NULL", "defines the type application is started")
	serverRange := flag.String("server-range", "2001", "defines port/s to setting up Middleware Server")
	localhostRange := flag.String("local-range", "3000", "defines port/s to which connects are needs to be forwarded")
	serverUrl := flag.String("url", "127.0.0.1:2001", "defines url to which client should connect to")
	flag.Parse()
	fmt.Printf("Flag Parsed %s %s %s %s %s \n", *secretFlag, *connectType, *serverRange, *localhostRange, *serverUrl)
	//fmt.Fprintf(flag.NewFlagSet(os.Args[0], flag.ExitOnError).Output(), "Usage of %s:\n", os.Args[0])
	ToClientChan := make(chan string)
	FromClientChan := make(chan string)
	if *connectType == "client" || *connectType == "NULL" {
		if *connectType == "NULL" {
			logManager.Warn("No Type Selected Starting Application in Client Mode")
		} else {
			logManager.Info("Application Started in Client Mode")
		}
		var wg sync.WaitGroup
		logManager.Info("Starting Start Client Mode")
		logManager.Info("Request from " + *localhostRange + " will be forwarded")
		ServerConnection, dialErr := net.Dial("tcp", *serverUrl)
		if dialErr != nil {
			logManager.Error(dialErr.Error(), 23)
		}
		logManager.Info("Connection to Server Established")
		wg.Add(2)
		go func() {
			//	 Handles Communication with Server and Client
			defer wg.Done()
			ReadBuffer := make([]byte, 1024)
			_, WriteErr := ServerConnection.Write([]byte("ATHENA_CONNECTION_READY"))
			if WriteErr != nil {
				logManager.Error(WriteErr.Error(), 12)
			}
			readCount, readErr := ServerConnection.Read(ReadBuffer)
			if readErr != nil {
				logManager.Error(readErr.Error(), 1)
			}
			if (string(ReadBuffer[:readCount])) == "ATHENA_CONNECTION_ACCEPT" {
				_, writeErr := ServerConnection.Write([]byte("ATHENA_CONNECTION_" + *localhostRange))
				if writeErr != nil {
					logManager.Error(writeErr.Error(), 1)
				}
				logManager.Info("Server Handshake Successful")
				for {
					readCount, readErr := ServerConnection.Read(ReadBuffer)
					if readErr != nil {
						logManager.Error(readErr.Error(), 1)
					}
					FromClientChan <- string(ReadBuffer[:readCount])
					go func() {
						for {
							select {
							case ResponseData := <-ToClientChan:
								_, writeErr := ServerConnection.Write([]byte(ResponseData))
								if writeErr != nil {
									logManager.Error(writeErr.Error(), 1)
								}
								break
								//	End Routine if Receive a Response
							}
						}
					}()
				}
			} else {
				logManager.Error("Received Unknown Handshake Request from Server , Exiting", 1)
			}
		}()

		go func() {
			defer wg.Done()
			//	Handles Communication with Resource Behind Firewall
			resourceConn, resourceConnErr := net.Dial("tcp", "127.0.0.1:"+*localhostRange)
			if resourceConnErr != nil {
				logManager.Error(resourceConnErr.Error(), 10)
			}
			logManager.Info("Connection to Resource Established :" + *localhostRange)
			for {
				select {
				case ReceivingData := <-FromClientChan:
					//	Write Data and wait for Response
					readBuffer := make([]byte, 1024)
					_, writeErr := resourceConn.Write([]byte(ReceivingData))
					if writeErr != nil {
						logManager.Error(writeErr.Error(), 1)
					}
					readCount, readErr := resourceConn.Read(readBuffer)
					if readErr != nil {
						logManager.Error(readErr.Error(), 1)
					}
					ToClientChan <- string(readBuffer[:readCount])
				}
			}

		}()
		wg.Wait()
		logManager.Info("Process lifecycle Ended, Exiting Program")

	} else {

		var wg sync.WaitGroup
		logManager.Info("Application Started in Server Mode")
		portRange := strings.Split(*serverRange, "-")

		startRange, connError := strconv.Atoi(portRange[0])
		var endRange int
		if len(portRange) == 2 {
			endRange, connError = strconv.Atoi(portRange[1])
		}
		if connError != nil {
			logManager.Error(connError.Error(), 23)
		}

		if len(portRange) == 2 {
			fmt.Printf("%d %d \n", startRange, endRange)
		} else if len(portRange) == 1 {

			listenerConn, listenerErr := net.Listen("tcp", "127.0.0.1:"+strconv.Itoa(startRange))
			if listenerErr != nil {
				logManager.Error(listenerErr.Error(), 1)
			}
			logManager.Info("Server listening on Port " + strconv.Itoa(startRange))
			acceptConn, acceptErr := listenerConn.Accept() //Sync Mode
			if acceptErr != nil {
				logManager.Error(acceptErr.Error(), 1)
			}
			wg.Add(1)
			go func() {
				logManager.Info("New Client Connection Accepted")
				readBuffer := make([]byte, 1024)
				readCount, readErr := acceptConn.Read(readBuffer)
				if readErr != nil {
					logManager.Error(readErr.Error(), 1)
				}
				// The Server Received Request from Athena Client
				if string(readBuffer[:readCount]) == "ATHENA_CONNECTION_READY" {
					logManager.Info("ATHENA CLIENT CONNECTING")
					_, writeErr := acceptConn.Write([]byte("ATHENA_CONNECTION_ACCEPT"))
					if writeErr != nil {
						logManager.Error(writeErr.Error(), 1)
					}
					readCount, readErr := acceptConn.Read(readBuffer)
					if readErr != nil {
						logManager.Error(readErr.Error(), 1)
					}
					portNumber := strings.Split(string(readBuffer[:readCount]), "ATHENA_CONNECTION_")
					logManager.Info("Client-Server Handshake Complete :" + portNumber[1])
					//	Allocate a Separate routine for client
					wg.Add(1)
					go func() {
						defer wg.Done()
						for {
							select {
							case IncomingData := <-ToClientChan:
								_, writeErr := acceptConn.Write([]byte(IncomingData))
								if writeErr != nil {
									logManager.Error(writeErr.Error(), 1)
								}
								go func() {
									readCount, readErr := acceptConn.Read(readBuffer)
									if readErr != nil {
										logManager.Error(readErr.Error(), 1)
									}
									FromClientChan <- string(readBuffer[:readCount])
								}()
							}
						}
					}()

				} else {
					fmt.Println("Receved External Client Req")
					fmt.Println(string(readBuffer[:readCount]))
					// Handle Connection from External Client
					wg.Add(1)
					go func() {
						ToClientChan <- string(readBuffer[:readCount])
						for {
							select {
							case receivedResponse := <-FromClientChan:
								_, writtenErr := acceptConn.Write([]byte(receivedResponse))
								if writtenErr != nil {
									logManager.Error(writtenErr.Error(), 1)
								}

							}
						}
					}()
				}
			}()

			wg.Wait()
		} else {
			logManager.Error("Could not parse range for server port Exiting", 44)
		}
		logManager.Info("Server LifeCycle Ended , Exiting Program")

	}

	//isServer := true
	//if isServer {
	//	con, err := net.Listen("tcp", "127.0.0.1:2002")
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	conn, err := con.Accept()
	//
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	bytess := make([]byte, 1024)
	//	bytes, erro := conn.Read(bytess)
	//	if erro != nil {
	//		log.Fatal(erro)
	//	}
	//
	//	fmt.Printf("Value Recievd \n%s \n\n", bytess[:bytes])
	//	fmt.Println("Connecting to Resource Server")
	//	connectio, err := net.Dial("tcp", "127.0.0.1:3000")
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//
	//	fmt.Println("Connecting Successful")
	//	fmt.Println("Now Writing Data to Client")
	//	Writentytes, err := connectio.Write(bytess[:bytes])
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	fmt.Printf("Written Bytes %d", Writentytes)
	//
	//	StorageBuffer := make([]byte, 1024)
	//	n, err := connectio.Read(StorageBuffer)
	//	if err != nil {
	//		fmt.Println("Error:", err)
	//		return
	//	}
	//
	//	// Process and use the data (here, we'll just print it)
	//	//fmt.Printf("Received: %s\n", StorageBuffer[:n])
	//	conn.Write(StorageBuffer[:n])
	//	//for {
	//	//	// Read data from the client
	//	//
	//	//}
	//
	//	//
	//	//byteWritten, err := conn.Write([]byte("Hello World"))
	//	//if err != nil {
	//	//	fmt.Println(err)
	//	//}
	//	for {
	//
	//	}
	//
	//	//fmt.Printf("The Bytes Written Are %d", byteWritten)
	//} else {
	//	_, err := net.Dial("tcp", "127.0.0.1:2002")
	//
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//
	//	//bytesData := make([]byte, 1024)
	//	//for {
	//	//	n, err := conn.Read(bytesData)
	//	//	if err != nil {
	//	//		fmt.Println("Error:", err)
	//	//		return
	//	//	}
	//	//	fmt.Println(string(bytesData[:n]))
	//	//}
	//}

}

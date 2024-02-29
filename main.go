package main

import (
	"flag"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"

	logManager "github.com/Dhananjay-JSR/Athena.git/cli"
	"github.com/Dhananjay-JSR/Athena.git/internal"
)

const secret_key = "Athena"

func main() {

	internal.InitWindowsEscape()
	secretFlag := flag.String("secret", secret_key, "Secret key is used to provide an encryption layer over client server communication to prevent anyone from hijacking connection")
	connectType := flag.String("type", "server", "defines the type application is started")
	serverRange := flag.String("server-range", "2001", "defines port/s to setting up Middleware Server")
	localhostRange := flag.String("local-range", "3000", "defines port/s to which connects are needs to be forwarded")
	serverUrl := flag.String("url", "athena.dhananjaay.dev:2001", "defines url to which client should connect to")
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

		ServerConnection, dialErr := net.Dial("tcp", *serverUrl)
		if dialErr != nil {
			logManager.Error(dialErr.Error(), 23)
		}
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
				logManager.Info("Request from " + *localhostRange + " will be forwarded")
				for {
					readCount, readErr := ServerConnection.Read(ReadBuffer)
					logManager.DEBUG("Received Connection From Server")
					if readErr != nil {
						logManager.Error(readErr.Error(), 1)
					}
					logManager.DEBUG("Reading Data from Server Successful")
					FromClientChan <- string(ReadBuffer[:readCount])
					logManager.DEBUG("Send to Worker Thread SUCCESS")
					go func() {
						for {
							select {
							case ResponseData := <-ToClientChan:
								logManager.DEBUG("Recieved Data from Worker Thread")
								_, writeErr := ServerConnection.Write([]byte(ResponseData))

								if writeErr != nil {
									logManager.Error(writeErr.Error(), 1)
								}
								logManager.DEBUG("Writing Data Successful ")

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
					logManager.DEBUG("Received Data to Worker Node")
					//	Write Data and wait for Response
					readBuffer := make([]byte, 1024)
					_, writeErr := resourceConn.Write([]byte(ReceivingData))

					if writeErr != nil {
						logManager.DEBUG("Worker Node Dead , Making Reconnecting Attempt")
						resourceConn, resourceConnErr = net.Dial("tcp", "127.0.0.1:"+*localhostRange)
						if resourceConnErr != nil {
							logManager.Error(resourceConnErr.Error(), 10)
						}
						logManager.Info("Connection to Resource Established :" + *localhostRange)
						_, _ = resourceConn.Write([]byte(ReceivingData))
						//logManager.Error(writeErr.Error(), 1)

					}
					logManager.DEBUG("Data Sent to Resource")
					readCount, readErr := resourceConn.Read(readBuffer)
					if readErr != nil {
						logManager.Error(readErr.Error(), 1)
					}
					logManager.DEBUG("Data Read From Resource")
					ToClientChan <- string(readBuffer[:readCount])
					logManager.DEBUG("Data Sent to Client Thread")
				}
			}

		}()
		wg.Wait()
		logManager.Info("Client lifecycle Ended, Exiting Program")

	} else {

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
			logManager.DEBUG("Connected to PORT listening for Request")
			if listenerErr != nil {
				logManager.Error(listenerErr.Error(), 1)
			}
			logManager.Info("Server listening on Port " + strconv.Itoa(startRange))
			readBuffer := make([]byte, 1024) // read and Store Buffer

			for {
				acceptConn, acceptErr := listenerConn.Accept() //Sync Mode
				logManager.DEBUG("Connection Accepted")
				if acceptErr != nil {
					logManager.Error(acceptErr.Error(), 1)
				}

				readCount, readErr := acceptConn.Read(readBuffer)
				logManager.DEBUG("Reading Request Successful")
				if readErr != nil {
					logManager.Error(readErr.Error(), 1)
				}

				// The Server Received Request from Athena Client
				fmt.Println(string(readBuffer[:readCount]))
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
					go func() {
						logManager.DEBUG("Athena Client Connection READY")
						for {
							select {
							case IncomingData := <-ToClientChan:
								logManager.DEBUG("Sending Response to Athena Client")
								_, writeErr := acceptConn.Write([]byte(IncomingData))
								if writeErr != nil {
									logManager.Error(writeErr.Error(), 1)
								}
								readCount, readErr = acceptConn.Read(readBuffer)
								if readErr != nil {
									logManager.Error(readErr.Error(), 1)
								}
								logManager.DEBUG("Received Response from ATHENA Client ")
								FromClientChan <- string(readBuffer[:readCount])
								logManager.DEBUG("Response Shared With Conection")

							}
						}
					}()

				} else {

					//3rd party client request
					logManager.Info("Connection Request")

					go func() {
						defer acceptConn.Close()
						//Message Passing between 3rd party and client routine

						ToClientChan <- string(readBuffer[:readCount])
						logManager.DEBUG("Data Passed to Client Manager")
					connectionPoint:
						for {
							select {
							case receivedResponse := <-FromClientChan:
								logManager.DEBUG("Received Response from Client Manager")
								_, _ = acceptConn.Write([]byte(receivedResponse))
								logManager.DEBUG("Killing Connection")
								break connectionPoint
							}

						}

					}()
				}
			}
		} else {
			logManager.Error("Could not parse range for server port Exiting", 44)
		}
		logManager.Info("Server LifeCycle Ended , Exiting Program")
	}
}

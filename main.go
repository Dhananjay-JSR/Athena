package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"

	"github.com/Dhananjay-JSR/Athena.git/internal/Server"
	"github.com/Dhananjay-JSR/Athena.git/internal/client"

	logManager "github.com/Dhananjay-JSR/Athena.git/cli"
	"github.com/Dhananjay-JSR/Athena.git/internal"
)

const secret_key = "Athena"
const BUFFER_SIZE = 1024

func main() {

	internal.InitWindowsEscape()
	secretFlag := flag.String("secret", secret_key, "Secret key is used to provide an encryption layer over client server communication to prevent anyone from hijacking connection")
	connectType := flag.String("type", "server", "defines the type application is started")
	serverRange := flag.String("server-range", "2001", "defines port/s to setting up Middleware Server")
	localhostRange := flag.String("local-range", "27017", "defines port/s to which connects are needs to be forwarded")
	serverUrl := flag.String("url", "127.0.0.1:2001", "defines url to which client should connect to")
	flag.Parse()

	fmt.Printf("Flag Parsed %s %s %s %s %s \n", *secretFlag, *connectType, *serverRange, *localhostRange, *serverUrl)
	//fmt.Fprintf(flag.NewFlagSet(os.Args[0], flag.ExitOnError).Output(), "Usage of %s:\n", os.Args[0])
	ToClientChan := make(chan []byte)
	FromClientChan := make(chan []byte)

	if *connectType == "client" || *connectType == "NULL" {
		if *connectType == "NULL" {
			logManager.Warn("No Type Selected Starting Application in Client Mode")
		} else {
			logManager.Info("Application Started in Client Mode")
		}
		var wg sync.WaitGroup
		logManager.Info("Initializing Server Handshake")

		ServerConnection, dialErr := net.Dial("tcp", *serverUrl)
		//Connecting to Server
		if dialErr != nil {
			logManager.Error(dialErr.Error(), 23)
		}

		wg.Add(2)
		//Added Wait Group of 2
		//Start Client Thread
		go client.ClientHandler(&wg, ServerConnection, connectType, FromClientChan, ToClientChan)

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

			//listenerConn, listenErr := net.Listen("tcp", "127.0.0.1:2001")
			//if listenErr != nil {
			//	log.Fatal(listenErr)
			//}
			//
			//defer listenerConn.Close()
			//
			//for {
			//	acceptedCon, acceptedErr := listenerConn.Accept()
			//	if acceptedErr != nil {
			//		log.Println("Accept error:", acceptedErr)
			//		continue
			//	}
			//
			//	go func(conn net.Conn) {
			//		defer conn.Close()
			//
			//		dialerConn, dialerErr := net.Dial("tcp", "127.0.0.1:3000")
			//		if dialerErr != nil {
			//			log.Println("Dialer error:", dialerErr)
			//			return
			//		}
			//		defer dialerConn.Close()
			//
			//		go func() {
			//			defer dialerConn.Close()
			//			io.Copy(dialerConn, conn)
			//		}()
			//
			//		io.Copy(conn, dialerConn)
			//	}(acceptedCon)
			//}

			listenerConn, listenerErr := net.Listen("tcp", "127.0.0.1:"+strconv.Itoa(startRange))
			logManager.DEBUG("Connected to PORT listening for Request")
			if listenerErr != nil {
				logManager.Error(listenerErr.Error(), 1)
			}
			logManager.Info("Server listening on Port " + strconv.Itoa(startRange))

			for {
				readBuffer := make([]byte, 4076) // read and Store Buffer

				acceptConn, acceptErr := listenerConn.Accept() //Sync Mode

				logManager.DEBUG("Connection Accepted")
				if acceptErr != nil {

					logManager.Error(acceptErr.Error(), 1)
				}
				readCount, readErr := acceptConn.Read(readBuffer)
				logManager.DEBUG("Reading Request Successful")
				if readErr != nil {
					if readErr == io.EOF {
						logManager.DEBUG("CONNECTION CLOSEDDD !!!")
					}
				}
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
					logManager.DEBUG("Athena Client Connection READY")
					go Server.MasterClientSynchronizer(ToClientChan, FromClientChan, &acceptConn)
				} else {
					fmt.Println("SERVER RECEUVED NEW CONNECTION REQUEST")
					logManager.Info("NEW Connection Request")
					go Server.ExternalConnectionHandler(&acceptConn, ToClientChan, readBuffer, readCount, FromClientChan)
				}
			}
		} else {
			logManager.Error("Could not parse range for server port Exiting", 44)
		}
		logManager.Info("Server LifeCycle Ended , Exiting Program")
	}
}

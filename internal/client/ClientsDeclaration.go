package client

import (
	logManager "github.com/Dhananjay-JSR/Athena.git/cli"
	"net"
	"sync"
)

func ClientHandler(wg *sync.WaitGroup, ServerConnection net.Conn, localhostRange *string, FromClientChan chan string, ToClientChan chan string) {
	defer wg.Done()
	ReadBuffer := make([]byte, 1024)
	_, WriteErr := ServerConnection.Write([]byte("ATHENA_CONNECTION_READY"))
	logManager.Info("Waiting for Server Response")
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
			logManager.DEBUG("Received Connection From Server")
			if readErr != nil {
				logManager.Error(readErr.Error(), 1)
			}
			logManager.DEBUG("Reading Data from Server Successful")
			FromClientChan <- string(ReadBuffer[:readCount])
			logManager.DEBUG("Sending to Resource Thread SUCCESS")
			go func() {
				for {
					select {
					case ResponseData := <-ToClientChan:
						logManager.DEBUG("Recieved Data from Resource Thread")
						_, writeErr := ServerConnection.Write([]byte(ResponseData))

						if writeErr != nil {
							logManager.Error(writeErr.Error(), 1)
						}
						logManager.DEBUG("Writing Data Successful Socket Successful")
					}
				}
			}()
		}
	} else {
		logManager.Error("Received Unknown Handshake Request from Server , Exiting", 1)
	}
}

func ResourceConnector(wg *sync.WaitGroup, localhostRange *string, FromClientChan chan string, ToClientChan chan string) {

	defer wg.Done()
	resourceConn, resourceConnErr := net.Dial("tcp", "127.0.0.1:"+*localhostRange)
	if resourceConnErr != nil {
		logManager.Error(resourceConnErr.Error(), 10)
	}
	logManager.Info("Connection to Resource Established :" + *localhostRange)
ConnectionLoop:
	for {
		select {
		case ReceivingData := <-FromClientChan:
			logManager.DEBUG("Received Data to Worker Node")
			//	Write Data and wait for Response
			readBuffer := make([]byte, 1024)
			_, writeErr := resourceConn.Write([]byte(ReceivingData))

			if writeErr != nil {
				logManager.DEBUG("Worker Node Dead , Making Reconnecting Attempt")
				go ResourceConnector(wg, localhostRange, FromClientChan, ToClientChan)
				break ConnectionLoop

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

}

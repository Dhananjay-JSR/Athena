package client

import (
	"fmt"
	"log"
	"net"
	"sync"

	logManager "github.com/Dhananjay-JSR/Athena.git/cli"
)

const BUFFER_SIZE = 1024

//go client.ResourceConnector(&wg, localhostRange, FromClientChan, ToClientChan)

func ClientHandler(wg *sync.WaitGroup, ServerConnection net.Conn, localhostRange *string, FromClientChan chan []byte, ToClientChan chan []byte) {

	defer wg.Done()

	ReadBuffer := make([]byte, BUFFER_SIZE)
	//Allocator := make([]byte, BUFFER_SIZE)

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
			//for readCount, readErr := ServerConnection.Read(ReadBuffer); readCount != 0; {
			//	if readErr != nil {
			//		log.Println(readErr)
			//	}
			//	Allocator = append(Allocator, ReadBuffer[:readCount]...)
			//}

			readCount, readErr := ServerConnection.Read(ReadBuffer)
			{
				if readErr != nil {
					log.Println(readErr)
				}

			}
			fmt.Println("Read Response ftom server , Sending it to resource allocator")
			//FromClientChan <- Allocator
			//ResponseData := <-ToClientChan
			ByteFromResourve := ResourceConnector(localhostRange, ReadBuffer[:readCount])
			fmt.Println(string(ByteFromResourve))
			WriteCount, writeErr := ServerConnection.Write(ByteFromResourve)
			fmt.Println("DATA SENT BACK TO MIDD SERVER ", WriteCount)
			if writeErr != nil {
				log.Println("Write Error")

			}

		}

		//for {

		//readCount, readErr := ServerConnection.Read(ReadBuffer)
		//logManager.DEBUG("Received Connection From Server")
		//if readErr != nil {
		//	logManager.Error(readErr.Error(), 1)
		//}
		//logManager.DEBUG("Reading Data from Server Successful")
		//// if no error send TO channel "FROMClientChan"
		//FromClientChan <- ReadBuffer[:readCount]
		//logManager.DEBUG("Sending to Resource Thread SUCCESS .. WAITINGG FOR RES")
		//
		//ResponseData := <-ToClientChan
		//logManager.DEBUG("Recieved Data from Resource Thread")
		//_, writeErr := ServerConnection.Write([]byte(ResponseData))
		//if writeErr != nil {
		//	logManager.Error(writeErr.Error(), 1)
		//}
		//logManager.DEBUG("Writing Data Successful Socket Successful")

		//}
	} else {
		logManager.Error("Received Unknown Handshake Request from Server , Exiting", 1)
	}
}

func ResourceConnector(localhostRange *string, WrittenData []byte) []byte {

	//
	//defer func() {
	//	fmt.Println("Resource Connection Operation Ending Releasing Resource")
	//	wg.Done()
	//
	//}()
	//if resourceConnErr != nil {
	//	logManager.Error(resourceConnErr.Error(), 10)
	//}
	//logManager.Info("Connection to Resource Established :" + *localhostRange)
	readBuffer := make([]byte, BUFFER_SIZE)
	// Allocator := make([]byte, BUFFER_SIZE)
	//for {
	//RescBuffer := <-FromClientChan
	resourceConn, resourceConnErr := net.Dial("tcp", "127.0.0.1:"+*localhostRange)
	if resourceConnErr != nil {
		log.Println(resourceConnErr)

	}
	defer resourceConn.Close()
	resourceConn.Write(WrittenData)
	readCount, readErr := resourceConn.Read(readBuffer)
	if readErr != nil {
		log.Println(readErr)
	}
	// DEBUG: Slice Error
	// FUCK YOU BELOW CODE
	// Allocator = append(Allocator, readBuffer[:readCount]...)

	//ToClientChan <- Allocator
	// Send Back Actual Response with their ReadCount ,  Now Extra Data Which was being Sent
	return readBuffer[:readCount]

	//}

}

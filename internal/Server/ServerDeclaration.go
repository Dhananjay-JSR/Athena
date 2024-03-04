package Server

import (
	"log"
	"net"
)

const BUFFER_SIZE = 1024

func MasterClientSynchronizer(ToClientChan chan []byte, FromClientChan chan []byte, AthenaClient *net.Conn) {
	//readBuffer := make([]byte, BUFFER_SIZE)

	defer (*AthenaClient).Close()
	// Read from the Channel
	// Write to the connection
	go func(AthenaClient *net.Conn) {
		//
		for {
			RecvData := <-ToClientChan
			(*AthenaClient).Write(RecvData)
		}
	}(AthenaClient)
	// Read from the connection
	// Write to the Channel
	go func(AthenaClient *net.Conn) {

		for {
			readBuffer := make([]byte, BUFFER_SIZE)
			readCount, readErr := (*AthenaClient).Read(readBuffer)
			if readErr != nil {
				if readErr == nil {
					return
				}
			}
			FromClientChan <- readBuffer[:readCount]
		}
	}(AthenaClient)
}

func ExternalConnectionHandler(AcceptedConnection *net.Conn, ToClientChan chan []byte, readBuffer []byte, readCount int, FromClientChan chan []byte) {
	defer (*AcceptedConnection).Close() // Ensure connection is closed when function exits
	// Read from the connection
	// Write to the Channel
	go func(AcceptedConnection *net.Conn) {
		for {
			readCount, readErr := (*AcceptedConnection).Read(readBuffer)
			if readErr != nil {
				log.Println(readErr)
			}
			ToClientChan <- readBuffer[:readCount]
		}

	}(AcceptedConnection)

	// Read from the Channel
	// Write to the connection
	go func(AcceptedConnection *net.Conn) {
		for {
			RecvData := <-FromClientChan
			(*AcceptedConnection).Write(RecvData)
		}
	}(AcceptedConnection)

}

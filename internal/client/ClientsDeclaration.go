package client

import (
	"net"
	"sync"
)

const BUFFER_SIZE = 4098

//go client.ResourceConnector(&wg, localhostRange, FromClientChan, ToClientChan)

func ClientHandler(wg *sync.WaitGroup, ServerConnection net.Conn, localhostRange *string, FromClientChan chan []byte, ToClientChan chan []byte) {
	go ResourceConnector(localhostRange, FromClientChan, ToClientChan)
	defer wg.Done()
	defer ServerConnection.Close()
	// Reaf from the Channel
	// Write to the connection
	go func(ServerConnection net.Conn) {
		for {
			RecvData := <-ToClientChan
			ServerConnection.Write(RecvData)
		}
	}(ServerConnection)

	// Read from the connection
	// Write to the Channel
	go func(ServerConnection net.Conn) {
		readBuffer := make([]byte, BUFFER_SIZE)
		for {
			readCount, readErr := ServerConnection.Read(readBuffer)
			if readErr != nil {
				if readErr == nil {
					return
				}
			}
			FromClientChan <- readBuffer[:readCount]
		}
	}(ServerConnection)

}

func ResourceConnector(localhostRange *string, FromClientChan chan []byte, ToClientChan chan []byte) {
	dialerConn, dialerErr := net.Dial("tcp", "127.0.0.1:3000")
	if dialerErr != nil {
		return
	}
	defer dialerConn.Close()

	// Read from the Channel
	// Write to the connection
	go func() {
		for {
			RecvData := <-FromClientChan
			dialerConn.Write(RecvData)
		}
	}()

	// Read from the connection
	// Write to the Channel
	go func(dial net.Conn) {
		readBuffer := make([]byte, BUFFER_SIZE)
		for {
			readCount, readErr := dial.Read(readBuffer)
			if readErr != nil {
				if readErr == nil {
					return
				}
			}
			ToClientChan <- readBuffer[:readCount]
		}
	}(dialerConn)

}

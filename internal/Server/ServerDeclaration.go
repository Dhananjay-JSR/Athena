package Server

import (
	"fmt"
	"log"
	"net"

	logManager "github.com/Dhananjay-JSR/Athena.git/cli"
)

const BUFFER_SIZE = 1024

func MasterClientSynchronizer(ToClientChan chan []byte, FromClientChan chan []byte, acceptConn *net.Conn) {
	//readBuffer := make([]byte, BUFFER_SIZE)
	Allocator := make([]byte, BUFFER_SIZE*4)
	defer (*acceptConn).Close()
	for {
		IncomingData := <-ToClientChan
		logManager.DEBUG("Sending Response to Athena Client")

		_, writeErr := (*acceptConn).Write(IncomingData)
		if writeErr != nil {
			logManager.Error(writeErr.Error(), 1)
		}
		// logManager.DEBUG("Written " + strconv.Itoa(writeCount) + " Number ")
		//for readCount, readErr := (*acceptConn).Read(readBuffer); readCount != 0; {
		//	if readErr != nil {
		//		log.Println(readErr)
		//		Allocator = append(Allocator, readBuffer[:readCount]...)
		//	}
		//}

		readCount, readErr := (*acceptConn).Read(Allocator)
		if readErr != nil {
			log.Println(readErr)
			//Allocator = append(Allocator, readBuffer[:readCount]...)
		}
		fmt.Println(readCount)
		//FromClientChan <- Allocator[:readCount]
		FromClientChan <- Allocator[:readCount]

	}
}

func ExternalConnectionHandler(acceptConn *net.Conn, ToClientChan chan []byte, readBuffer []byte, readCount int, FromClientChan chan []byte) {
	defer (*acceptConn).Close() // Ensure connection is closed when function exits

	Allocator := make([]byte, readCount)

	copy(Allocator, readBuffer[:readCount])

	log.Println("Passing Data to Client Manager")

	ToClientChan <- Allocator

	response := <-FromClientChan
	_, writeErr := (*acceptConn).Write(response)
	if writeErr != nil {
		log.Println("Error writing HTTP response:", writeErr)
	}

}

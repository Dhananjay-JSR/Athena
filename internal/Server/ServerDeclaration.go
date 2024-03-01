package Server

import (
	logManager "github.com/Dhananjay-JSR/Athena.git/cli"
	"net"
)

func MasterClientSynchronizer(ToClientChan chan string, FromClientChan chan string, acceptConn net.Conn, readBuffer []byte) {
	defer acceptConn.Close()
	for {
		select {
		case IncomingData := <-ToClientChan:
			logManager.DEBUG("Sending Response to Athena Client")
			_, writeErr := acceptConn.Write([]byte(IncomingData))
			if writeErr != nil {
				logManager.Error(writeErr.Error(), 1)
			}
			readCount, readErr := acceptConn.Read(readBuffer)
			if readErr != nil {
				logManager.Error(readErr.Error(), 1)
			}
			logManager.DEBUG("Received Response from ATHENA Client ")
			FromClientChan <- string(readBuffer[:readCount])
			logManager.DEBUG("Response Shared With Connection")

		}
	}
}

func ExternalConnectionHandler(acceptConn *net.Conn, ToClientChan chan string, readBuffer []byte, readCount int, FromClientChan chan string) {
	defer (*acceptConn).Close()
	//Message Passing between 3rd party and client routine

	ToClientChan <- string(readBuffer[:readCount])
	readBufferLoc := make([]byte, 1024)

	logManager.DEBUG("Data Passed to Client Manager " + string(readBuffer[:readCount]))
connectionPoint:

	for {
		select {
		case receivedResponse := <-FromClientChan:
			logManager.DEBUG("Received Response from Client Manager")
			_, _ = (*acceptConn).Write([]byte(receivedResponse))
			readCountLoc, readErrLoc := (*acceptConn).Read(readBufferLoc)
			if readErrLoc != nil {
				if readErrLoc.Error() == "EOF" {
					logManager.Info("CLIENT ENDED CONNECTION KILL GOROUTINE")
					break connectionPoint
				}

				//logManager.Error(readErrLoc.Error(), 1)
				//logManager.DEBUG(readErrLoc.Error())
			}
			ToClientChan <- string(readBufferLoc[:readCountLoc])
		}

	}

}

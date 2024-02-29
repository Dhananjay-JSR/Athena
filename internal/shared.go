package internal

import (
	"net"
	"runtime"
	"sync"
	"syscall"
)

func InitWindowsEscape() {
	if runtime.GOOS == "windows" {
		//this enables ANSI Escape on Windows Console
		EnableVirtualTerminalProcessing(syscall.Stdout, true)
	}
}

func HandleConnection(acceptConn net.Conn, wg *sync.WaitGroup, ToClientChan chan string, FromClientChan chan string) {

}

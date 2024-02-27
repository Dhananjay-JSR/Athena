package logManager

import (
	"fmt"
	"log"
	"os"
)

const (
	_FG_RED     = 31
	_FG_GREEN   = 32
	_FG_BLUE    = 34
	_FG_DEFAULT = 39
	_BG_RED     = 41
	_BG_GREEN   = 42
	_FG_YELLOW  = 33
	_BG_BLUE    = 44
	_BG_DEFAULT = 49
	_BG_WHITE   = 47
	_BOLD       = 1
)

//ESCAPE  [

func Info(msg string) {

	fmt.Printf("\u001B[%d;%d;%dmINFO \u001B[0m| %s \n", _FG_GREEN, _BG_DEFAULT, _BOLD, msg)
}

func Warn(msg string) {
	l := log.New(os.Stderr, "", 0)
	l.Printf("\u001B[%d;%d;%dmWARN \u001B[0m| %s \n", _FG_YELLOW, _BG_DEFAULT, _BOLD, msg)
}

func Error(msg string, exitCode int) {
	log.Printf("\u001B[%d;%d;%dmERROR \u001B[0m| %s \n", _FG_RED, _BG_DEFAULT, _BOLD, msg)
	os.Exit(exitCode)

}

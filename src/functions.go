package main

import (
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

// logger implements the paho.Logger interface
type logger struct {
	prefix string
}

// Println is the library provided NOOPLogger's
// implementation of the required interface function()
func (l logger) Println(v ...interface{}) {
	//fmt.Println(append([]interface{}{l.prefix + ":"}, v...)...)
	log.Debug(append([]interface{}{l.prefix + ":"}, v...)...)

}

// Printf is the library provided NOOPLogger's
// implementation of the required interface function(){}
func (l logger) Printf(format string, v ...interface{}) {
	//	if len(format) > 0 && format[len(format)-1] != '\n' {
	//		format = format + "\n" // some log calls in paho do not add \n
	//	}
	//	fmt.Printf(l.prefix+":"+format, v...)
	log.Debug(append([]interface{}{l.prefix + ":"}, v...)...)

}

func sleep_sec(seconds int, wg *sync.WaitGroup) {
	defer wg.Done()

	time.Sleep(time.Second * time.Duration(seconds))
}

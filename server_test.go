package rlog

import (
	"testing"

	"bytes"
	"compress/gzip"
	"io/ioutil"
	"log"
)

func TestNewServer(t *testing.T) {
	NewServer("0.0.0.0:1234", 100, 100, onPack)

	<-make(chan int)
}

func onPack(key string, value []byte) {
	gz, _ := gzip.NewReader(bytes.NewBuffer(value))
	c, _ := ioutil.ReadAll(gz)
	log.Println("onPack:", key, string(c), len(value), len(c))
}

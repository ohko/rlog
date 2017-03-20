package rlog

import (
	"fmt"
	"time"

	"github.com/ohko/omsg"
)

// Client ...
type Client struct {
	cache     chan []byte // 缓冲日志
	cacheSize int         // 缓冲大小
	omsg      *omsg.Client
}

// NewClient ...
func NewClient(addr string, cacheSize int) *Client {
	o := &Client{cacheSize: cacheSize}
	o.cache = make(chan []byte, cacheSize)
	o.omsg = omsg.NewClient(nil, nil)
	go o.omsg.Connect(addr, true, 5, 5)

	go func() {
		for {
			o.do()
		}
	}()

	return o
}

// Send 发送日志
func (o *Client) Send(key, value string) {
	select {
	case o.cache <- enData(key, value):
	default:
	}
}

// Status 发送日志
func (o *Client) Status() string {
	h := fmt.Sprintf("日志缓冲：%v/%v\n", len(o.cache), o.cacheSize)
	return h
}

func (o *Client) do() {
	defer func() { recover() }()
	for {
		_d := <-o.cache

		// 直到发送成功
		for {
			if 0 != o.omsg.Send(_d) {
				break
			}
			time.Sleep(time.Second)
		}
	}
}

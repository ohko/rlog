package rlog

import "github.com/ohko/omsg"

import "time"

// Client ...
type Client struct {
	cache chan []byte // 缓冲日志
	omsg  *omsg.Client
}

// NewClient ...
func NewClient(addr string, cacheSize int) *Client {
	o := &Client{}
	o.cache = make(chan []byte, cacheSize)
	o.omsg = omsg.NewClient(nil, nil)
	o.omsg.Connect(addr, true, 5, 5)

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

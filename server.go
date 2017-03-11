package rlog

import (
	"bytes"
	"compress/gzip"
	"net"

	"github.com/ohko/omsg"
)

// Server ...
type Server struct {
	logs     map[string]chan string
	cache    chan string // 缓冲日志
	packSize int         // 打包处理数量
	omsg     *omsg.Server
	onPack   func(string, []byte)
}

// NewServer ...
func NewServer(addr string, cacheSize, packSize int, onPack func(string, []byte)) *Server {
	o := &Server{packSize: packSize, onPack: onPack, logs: make(map[string]chan string)}
	o.cache = make(chan string, cacheSize)
	o.omsg = omsg.NewServer(o.onData, nil, nil)
	o.omsg.StartServer(addr)

	go func() {
		for {
			o.do()
		}
	}()

	return o
}

// 收到数据
func (o *Server) onData(conn net.Conn, data []byte) {
	select {
	case o.cache <- string(data):
	default:
	}
}

func (o *Server) do() {
	defer func() { recover() }()
	for {
		key, value := deData(<-o.cache)

		// 没有此KEY，创建一个
		if _, ok := o.logs[key]; !ok {
			o.logs[key] = make(chan string, o.packSize)
		}
		o.logs[key] <- value

		// 达到压缩尺寸
		if len(o.logs[key]) >= o.packSize {
			bs := bytes.NewBuffer(nil)
			gz, _ := gzip.NewWriterLevel(bs, gzip.BestCompression)
			for i := 0; i < o.packSize; i++ {
				gz.Write([]byte(<-o.logs[key]))
			}
			gz.Close()

			// 回调处理
			if o.onPack != nil {
				o.onPack(key, bs.Bytes())
			}
		}
	}
}

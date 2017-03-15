package rlog

import (
	"bytes"
	"compress/gzip"
	"net"

	"fmt"

	"log"

	"github.com/ohko/omsg"
)

// Server ...
type Server struct {
	logs      map[string]chan string
	cache     chan string // 缓冲日志
	cacheSize int         // 缓冲大小
	packSize  int         // 打包处理数量
	omsg      *omsg.Server
	onPack    func(string, []byte)
}

// NewServer ...
func NewServer(addr string, cacheSize, packSize int, onPack func(string, []byte)) (*Server, error) {
	o := &Server{cacheSize: cacheSize, packSize: packSize, onPack: onPack, logs: make(map[string]chan string)}
	o.cache = make(chan string, cacheSize)
	o.omsg = omsg.NewServer(o.onData, nil, nil)
	if err := o.omsg.StartServer(addr); err != nil {
		log.Println("ERROR:", err)
		return nil, err
	}

	go func() {
		for {
			o.do()
		}
	}()

	return o, nil
}

// Status 获取当前状态
func (o *Server) Status() string {
	h := fmt.Sprintf("日志缓冲：%v/%v\n", len(o.cache), o.cacheSize)
	for k, v := range o.logs {
		h += fmt.Sprintf("%v：%v/%v\n", k, len(v), o.packSize)
	}
	return h
}

// StoreNow 立刻处理数据
func (o *Server) StoreNow() {
	for k := range o.logs {
		o.compress(k)
	}
}

// SetPackSize 设置pack大小
func (o *Server) SetPackSize(packsize int) int {
	if packsize > 100 {
		o.packSize = packsize
	}
	return o.packSize
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
			o.compress(key)
		}
	}
}

func (o *Server) compress(key string) {
	length := len(o.logs[key])
	if length == 0 {
		return
	}
	bs := bytes.NewBuffer(nil)
	gz, _ := gzip.NewWriterLevel(bs, gzip.BestCompression)
	for i := 0; i < length; i++ {
		gz.Write([]byte(<-o.logs[key]))
		gz.Write([]byte("\n"))
	}
	gz.Close()

	// 回调处理
	if o.onPack != nil {
		o.onPack(key, bs.Bytes())
	}
}

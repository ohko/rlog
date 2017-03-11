# rlog
远程日志解决方案，生产日志和处理日志均带有缓冲。

# 使用
```
$ go get -u github.com/ohko/rlog
```

# Server
```
NewServer("0.0.0.0:1234", 100000, 100000, onPack)

func onPack(key string, value []byte) {
	// 日志达到处理阀值
}
```

# Client
```
// 本地缓冲10000万条日志
l := NewClient("0.0.0.0:1234", 10000)

// 发送日志
l.Send(...)
```
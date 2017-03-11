package rlog

import "testing"

import "fmt"

func TestNewClient(t *testing.T) {
	l := NewClient("0.0.0.0:1234", 100)

	// 模拟放松10个不同的key，每个key有1000条数据
	for i := 0; i < 10000; i++ {
		for j := 0; j < 10; j++ {
			l.Send(fmt.Sprintf("k-%v", j), fmt.Sprintf("k-%v-%v,", j, i))
		}
	}

	<-make(chan int)
}

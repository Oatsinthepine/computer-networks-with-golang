package main

import "time"

func main() {
	sender := NewTCPSender(1000, 3*time.Second)

	go sender.Start()

	sender.appDataChan <- []byte("Hello")
	sender.appDataChan <- []byte("World")

	time.Sleep(500 * time.Millisecond)
	sender.ackChan <- 1005

	// 超过超时时间，观察全部确认后是否还会重传。
	time.Sleep(4 * time.Second)
}

package main

import (
	"fmt"
	"time"
)

// Segment 表示 TCP 报文段
type Segment struct {
	Seq  uint32
	Data []byte
}

// TCPSender 实现简化版 TCP 发送方逻辑
type TCPSender struct {
	SendBase   uint32
	NextSeqNum uint32
	// 事件通道
	appDataChan chan []byte // 事件 1: 来自应用层的数据
	ackChan     chan uint32 // 事件 3: 收到 ACK
	// 模拟未确认的报文缓冲区
	unackedSegments []Segment

	// 定时器
	rttTimeout  time.Duration
	timer       *time.Timer
	timerActive bool
}

func NewTCPSender(initialSeqNum uint32, timeout time.Duration) *TCPSender {
	return &TCPSender{
		SendBase:        initialSeqNum,
		NextSeqNum:      initialSeqNum,
		appDataChan:     make(chan []byte),
		ackChan:         make(chan uint32),
		unackedSegments: make([]Segment, 0),
		rttTimeout:      timeout,
	}
}

func (s *TCPSender) Start() {
	var timerC <-chan time.Time // 用于 select 监听的 Channel，未启动时为 nil

	fmt.Printf("[TCP] 初始化完成，初始 SendBase: %d, NextSeqNum: %d\n", s.SendBase, s.NextSeqNum)

	for {
		select {
		// -------------------------------------------------------------
		// 事件 1: 上层应用有数据到达 (data received from application above)
		// -------------------------------------------------------------
		case data := <-s.appDataChan:
			// 创建 TCP 报文段并发给 IP 层
			seg := Segment{
				Seq:  s.NextSeqNum,
				Data: data,
			}
			s.passToIp(seg)
			s.unackedSegments = append(s.unackedSegments, seg)
			// if (timer currently not running) -> start timer
			if !s.timerActive {
				s.startTimer(&timerC)
			}
			// NextSeqNum = NextSeqNum + length(data)
			s.NextSeqNum += uint32(len(data))

		// -------------------------------------------------------------
		// 事件 2: 定时器超时 (timer timeout)
		// -------------------------------------------------------------
		case <-timerC:
			if len(s.unackedSegments) > 0 {
				oldestSeg := s.unackedSegments[0]
				fmt.Printf("[Timeout] 超时！重传最早未确认报文: Seq=%d\n", oldestSeg.Seq)
				s.passToIp(oldestSeg)

				// start timer
				s.startTimer(&timerC)
			}

		// -------------------------------------------------------------
		// 事件 3: 收到 ACK (ACK received, with ACK field value of y)
		// -------------------------------------------------------------
		case y := <-s.ackChan:
			fmt.Printf("[ACK] 收到 ACK=%d\n", y)
			// if (y > SendBase)
			if y > s.SendBase {
				s.SendBase = y
				// 从缓冲区移除已确认的报文 (Seq + Len <= y)
				newUnacked := make([]Segment, 0)
				for _, seg := range s.unackedSegments {
					if seg.Seq+uint32(len(seg.Data)) > y {
						newUnacked = append(newUnacked, seg)
					}
				}
				s.unackedSegments = newUnacked

				// if (there are currently not-yet-acknowledged segments)
				if len(s.unackedSegments) > 0 {
					s.startTimer(&timerC) // re-start timer
				} else {
					s.stopTimer(&timerC) // stop timer
				}
			}
		}
	}
}

// 模拟发送到 IP 层
func (s *TCPSender) passToIp(seg Segment) {
	fmt.Printf(" -> [IP Layer] 发送报文 Seq=%d, Length=%d\n", seg.Seq, len(seg.Data))
}

// 启动 / 重启定时器
func (s *TCPSender) startTimer(timerC *<-chan time.Time) {
	s.stopTimer(timerC)
	s.timer = time.NewTimer(s.rttTimeout)
	*timerC = s.timer.C
	s.timerActive = true
}

// 停止定时器
func (s *TCPSender) stopTimer(timerC *<-chan time.Time) {
	if s.timer != nil {
		s.timer.Stop()
	}
	*timerC = nil // 设置为 nil 后，select 会自动忽略该 channel
	s.timerActive = false
}

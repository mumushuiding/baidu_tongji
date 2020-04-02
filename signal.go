package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
)

// shutdownRequestChannel 用于关闭初始化
var shutdownRequestChannel = make(chan struct{})

// interruptSignals 定义默认的关闭触发信号
var interruptSignals = []os.Signal{os.Interrupt}

// interruptListener 监听关闭信号(Ctrl+C)
func interruptListener(s *http.Server) <-chan struct{} {
	c := make(chan struct{})
	go func() {
		interruptChannel := make(chan os.Signal, 1)
		signal.Notify(interruptChannel, interruptSignals...)
		select {
		case sig := <-interruptChannel:
			log.Printf("收到关闭信号 (%s). 关闭...\n", sig)
			s.Close()
		case <-shutdownRequestChannel:
			log.Println("关闭请求.关闭...")
			s.Close()
		}
		close(c)
		// 重复关闭信号处理
		for {
			select {
			case sig := <-interruptChannel:
				log.Printf("Received signal (%s).  Already "+
					"shutting down..\n", sig)
			case <-shutdownRequestChannel:
				log.Println("Shutdown requested.  Already " +
					"shutting down...")
			}
		}
	}()
	return c
}
func interruptRequested(interrupted <-chan struct{}) bool {
	select {
	case <-interrupted:
		return true
	default:
	}
	return false
}

package loji

import (
	"fmt"
	"sync"
	"time"
)

type LoadingEmoji struct {
	clocks      []rune
	signal      chan string
	isStart     bool
	loadingText string
	mx          sync.Mutex
}

func New() *LoadingEmoji {
	return NewLoading("🕐🕑🕒🕓🕔🕕🕖🕗🕘🕙🕚🕛")
}

func NewLoading(clocks string) *LoadingEmoji {
	if clocks == "" {
		clocks = "🕐🕑🕒🕓🕔🕕🕖🕗🕘🕙🕚🕛"
	}
	return &LoadingEmoji{
		clocks:  []rune(clocks),
		signal:  make(chan string),
		isStart: false,
	}
}

func (l *LoadingEmoji) Loading(msg string) {
	l.loadingText = msg
	if l.isStart == true {
		return
	}
	l.mx.Lock()
	l.isStart = true
	go func() {
		l.loading()
	}()
}

func (l *LoadingEmoji) Stop() {
	l.signal <- "stop"
}

func (l *LoadingEmoji) loading() {
	index := 0
	go func() {
		for true {
			fmt.Printf("\r %s%s", string(l.clocks[index]), l.loadingText)
			index++
			if index >= len(l.clocks) {
				index = 0
			}
			time.Sleep(100 * time.Millisecond)
			if l.isStart == false {
				fmt.Println("")
				l.mx.Unlock()
				break
			}
		}
	}()

	select {
	case _ = <-l.signal:
		l.isStart = false
	}

}

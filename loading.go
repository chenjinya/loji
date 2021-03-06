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
	tickFunc func()
	countDownIndex int
	sleepDuration time.Duration // ms
	mx          sync.Mutex
}

const DefaultClocks = "🕐🕑🕒🕓🕔🕕🕖🕗🕘🕙🕚🕛"

func New() *LoadingEmoji {
	return NewLoading("")
}

func NewLoading(clocks string) *LoadingEmoji {
	if clocks == "" {
		clocks = DefaultClocks
	}
	return &LoadingEmoji{
		clocks:  []rune(clocks),
		signal:  make(chan string),
		isStart: false,
		sleepDuration: 100 * time.Millisecond,
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
	if !l.isStart {
		return
	}
	l.signal <- "stop"
	l.mx.Unlock()
}

func (l *LoadingEmoji) NextTick(f func()) {
	l.tickFunc = f
}

func (l *LoadingEmoji) CountDownIndex() int {
	return l.countDownIndex
}

func (l *LoadingEmoji) SleepDuration() time.Duration {
	return l.sleepDuration
}

func (l *LoadingEmoji) loading() {
	index := 0
	l.countDownIndex = 0
	go func() {
		for true {
			if l.tickFunc != nil {
				l.tickFunc()
			}
			fmt.Printf("\r%s%s", string(l.clocks[index]), l.loadingText)
			index++
			if index >= len(l.clocks) {
				index = 0
			}
			l.countDownIndex ++
			time.Sleep(l.sleepDuration)
			if l.isStart == false {
				break
			}
		}
	}()

	select {
	case _ = <-l.signal:
		l.isStart = false
	}

}

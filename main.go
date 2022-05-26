package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Planner struct {
	mu         sync.Mutex
	wg         sync.WaitGroup
	once       sync.Once
	stopFn     func()
	signalChan chan os.Signal
}

func NewPlanner() *Planner {
	signalChan := make(chan os.Signal, 1)
	signal.Ignore(syscall.SIGHUP, syscall.SIGPIPE)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	return &Planner{
		signalChan: signalChan,
	}
}

func (p *Planner) Start(ctx context.Context, fn func()) {
	go func() {
		p.once.Do(func() {
			exit := func() {
				log.Println("exit2")
				fn()
			}

			p.wg.Add(1)
			ctxC, cancel := context.WithCancel(ctx)
			p.stopFn = cancel
		LOOP:
			for {
				select {
				case <-p.signalChan:
					exit()
					break LOOP
				case <-ctxC.Done():
					exit()
					break LOOP
				default:
					if p.mu.TryLock() {
						go func(ctx context.Context) {
							select {
							case <-p.signalChan:
								return
							case <-ctx.Done():
								return
							default:
								fn()
								p.mu.Unlock()
							}
						}(ctxC)
					}
				}
			}
			p.wg.Done()
		},
		)
	}()
}

func (p *Planner) Stop() {
	p.stopFn()
}

func (p *Planner) Wait() {
	p.wg.Wait()
}

func main() {
	p := NewPlanner()
	p.Start(context.Background(), fn)
	time.Sleep(time.Millisecond * 10530)
	p.Stop()
	p.Wait()
}

func fn() {
	time.Sleep(time.Second) // Имитация 1 секунды длительности работы
	log.Println("fn")       // Вывод в лог факта завершения выполнения
}

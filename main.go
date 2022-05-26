package main

import (
	"context"
	"log"
	"sync"
	"time"
)

type Planner struct {
	mu     sync.Mutex
	wg     sync.WaitGroup
	once   sync.Once
	stopFn func()
}

func NewPlanner() *Planner {
	return &Planner{}
}

func (p *Planner) Start(ctx context.Context, fn func()) {
	go func() {
		p.once.Do(func() {
			p.wg.Add(1)
			ctxC, cancel := context.WithCancel(ctx)
			p.stopFn = cancel
		LOOP:
			for {
				select {
				case <-ctxC.Done():
					log.Println("exit")
					fn()
					break LOOP
				default:
					if p.mu.TryLock() {
						go func(ctx context.Context) {
							select {
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
	time.Sleep(time.Millisecond * 1530)
	p.Stop()
	p.Wait()
}

func fn() {
	time.Sleep(time.Second) // Имитация 1 секунды длительности работы
	log.Println("fn")       // Вывод в лог факта завершения выполнения
}

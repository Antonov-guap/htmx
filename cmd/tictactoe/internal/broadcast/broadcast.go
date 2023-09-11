package broadcast

import (
	"log"
	"sync"
	"time"
)

// Broadcast реализует простую широковещательную шину
type Broadcast[T comparable] struct {
	mu          sync.RWMutex
	subscribers []chan T
	closed      bool

	Logger *log.Logger
}

// Subscribe добавляет нового подписчика и возвращает канал для получения сообщений
func (b *Broadcast[T]) Subscribe() <-chan T {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.closed {
		b.log("cannot subscribe, broadcast closed")
		return nil
	}

	ch := make(chan T)
	b.subscribers = append(b.subscribers, ch)

	b.log("new subscriber connected")
	return ch
}

// Send отправляет сообщение всем подписчикам
func (b *Broadcast[T]) Send(v T) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if b.closed {
		b.log("bus is already closed")
		return
	}

	for i, ch := range b.subscribers {
		select {
		case ch <- v:
			b.log("message sent to subscriber", i)
		case <-time.After(10 * time.Second):
			b.log("cannot send message, timout 10s")
			b.unsubscribeUnsafe(ch)
		}
	}

	if len(b.subscribers) == 0 {
		b.log("no subscribers to send message")
	}
	return
}

// Unsubscribe отписывает подписчика
func (b *Broadcast[T]) Unsubscribe(ch <-chan T) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.unsubscribeUnsafe(ch)
}

func (b *Broadcast[T]) unsubscribeUnsafe(ch <-chan T) {
	for i, subscriber := range b.subscribers {
		if subscriber == ch {
			close(subscriber)
			b.subscribers = append(b.subscribers[:i], b.subscribers[i+1:]...)
			b.log("subscriber", i, "disconnected")
			return
		}
	}
	b.log("cannot disconnect, subscriber not found")
}

// Close закрывает шину и все каналы подписчиков
func (b *Broadcast[T]) Close() {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.closed {
		b.log("cannot close, already closed")
		return
	}

	b.closed = true

	for i, ch := range b.subscribers {
		b.log("subscriber", i, "disconnected")
		close(ch)
	}
	if len(b.subscribers) == 0 {
		b.log("no subscribers to close")
	}
}

func (b *Broadcast[T]) log(message ...any) {
	if b.Logger == nil {
		return
	}
	b.Logger.Println(message...)
}

package queue

import (
	"context"
	"errors"
)

var (
	ErrQueueEmpty = errors.New("queue is empty")
	ErrQueueFull  = errors.New("queue is full")
)

// Queue 定义了队列接口
type Queue interface {
	// Push 将数据推送到队列
	Push(ctx context.Context, data []byte) error

	// Pop 从队列中获取数据
	Pop(ctx context.Context) ([]byte, error)

	// Close 关闭队列
	Close() error
}

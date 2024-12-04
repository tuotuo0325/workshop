package queue

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type QueueItem struct {
	ID        string    `json:"id"`
	Data      []byte    `json:"data"`
	Timestamp time.Time `json:"timestamp"`
}

type FileQueue struct {
	path      string
	mu        sync.Mutex
	processed map[string]bool
}

func NewFileQueue(path string) (*FileQueue, error) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	return &FileQueue{
		path:      path,
		processed: make(map[string]bool),
	}, nil
}

func (q *FileQueue) Push(ctx context.Context, data []byte) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	item := QueueItem{
		ID:        time.Now().Format("20060102150405.000000000"),
		Data:      data,
		Timestamp: time.Now(),
	}

	f, err := os.OpenFile(q.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	encoder := json.NewEncoder(f)
	return encoder.Encode(item)
}

func (q *FileQueue) Pop(ctx context.Context) ([]byte, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	f, err := os.OpenFile(q.path, os.O_RDWR, 0644)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrQueueEmpty
		}
		return nil, err
	}
	defer f.Close()

	decoder := json.NewDecoder(f)
	for {
		var item QueueItem
		err := decoder.Decode(&item)
		if err != nil {
			return nil, ErrQueueEmpty
		}

		if !q.processed[item.ID] {
			q.processed[item.ID] = true
			return item.Data, nil
		}
	}
}

func (q *FileQueue) Close() error {
	q.mu.Lock()
	defer q.mu.Unlock()

	// 清理已处理的记录
	q.processed = make(map[string]bool)
	return nil
}

package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"airbnb-cli/internal/crawler"
	"airbnb-cli/internal/model"
	"airbnb-cli/internal/queue"
	"airbnb-cli/internal/storage"
)

// Consumer 定义了消费者接口
type Consumer interface {
	// Start 启动消费者
	Start(ctx context.Context) error
	// Close 关闭消费者
	Close() error
}

// TaskConsumer 实现了任务消费者
type TaskConsumer struct {
	workers int
	queue   queue.Queue
	storage storage.Storage
	crawler crawler.Crawler
	wg      sync.WaitGroup
}

// NewTaskConsumer 创建任务消费者
func NewTaskConsumer(workers int, queue queue.Queue, storage storage.Storage) Consumer {
	return &TaskConsumer{
		workers: workers,
		queue:   queue,
		storage: storage,
		crawler: crawler.NewAirbnbCrawler(),
	}
}

func (c *TaskConsumer) Start(ctx context.Context) error {
	// 启动工作协程
	for i := 0; i < c.workers; i++ {
		c.wg.Add(1)
		go c.worker(ctx)
	}

	// 等待所有工作协程完成
	c.wg.Wait()
	return nil
}

func (c *TaskConsumer) worker(ctx context.Context) {
	defer c.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			// 从队列获取任务
			data, err := c.queue.Pop(ctx)
			if err == queue.ErrQueueEmpty {
				return
			}
			if err != nil {
				fmt.Printf("pop task from queue failed: %v\n", err)
				continue
			}

			var task model.Task
			if err := json.Unmarshal(data, &task); err != nil {
				fmt.Printf("unmarshal task failed: %v\n", err)
				continue
			}

			// 爬取数据
			hotel, err := c.crawler.Crawl(ctx, &task)
			if err != nil {
				fmt.Printf("crawl task failed: %v\n", err)
				continue
			}

			// 保存数据
			if err := c.storage.SaveHotel(ctx, hotel); err != nil {
				fmt.Printf("save hotel failed: %v\n", err)
				continue
			}
		}
	}
}

func (c *TaskConsumer) Close() error {
	if err := c.queue.Close(); err != nil {
		return fmt.Errorf("close queue failed: %w", err)
	}
	if err := c.storage.Close(); err != nil {
		return fmt.Errorf("close storage failed: %w", err)
	}
	return nil
}

package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"airbnb-cli/internal/consumer"
	"airbnb-cli/internal/queue"
	"airbnb-cli/internal/storage"
)

// ConsumerCommand 表示消费者命令
type ConsumerCommand struct {
	fs          *flag.FlagSet
	workers     int
	queuePath   string
	storagePath string
}

// NewConsumerCommand 创建消费者命令
func NewConsumerCommand() *ConsumerCommand {
	c := &ConsumerCommand{
		fs: flag.NewFlagSet("consumer", flag.ExitOnError),
	}

	c.fs.IntVar(&c.workers, "workers", 10, "Number of worker threads")
	c.fs.StringVar(&c.queuePath, "queue", "file://queue.data", "Queue server address")
	c.fs.StringVar(&c.storagePath, "storage", "data/hotels.json", "Storage file path")

	return c
}

// Run 运行消费者命令
func (c *ConsumerCommand) Run(args []string) error {
	if err := c.fs.Parse(args); err != nil {
		return err
	}

	fmt.Printf("Starting consumer with workers: %d, queue: %s, storage: %s\n",
		c.workers, c.queuePath, c.storagePath)

	// 创建存储目录
	if err := os.MkdirAll(filepath.Dir(c.storagePath), 0755); err != nil {
		return fmt.Errorf("create storage directory failed: %w", err)
	}

	fmt.Printf("Creating storage at: %s\n", c.storagePath)

	// 创建存储
	store, err := storage.NewFileStorage(c.storagePath)
	if err != nil {
		return fmt.Errorf("create storage failed: %w", err)
	}

	// 创建队列
	queueFile := c.queuePath[7:]
	fmt.Printf("Opening queue at: %s\n", queueFile)

	q, err := queue.NewFileQueue(queueFile)
	if err != nil {
		return fmt.Errorf("create queue failed: %w", err)
	}

	// 创建消费者
	cons := consumer.NewTaskConsumer(c.workers, q, store)

	fmt.Printf("Starting to process tasks with %d workers...\n", c.workers)
	if err := cons.Start(context.Background()); err != nil {
		return fmt.Errorf("consumer failed: %w", err)
	}

	if err := cons.Close(); err != nil {
		return fmt.Errorf("close consumer failed: %w", err)
	}

	fmt.Println("Consumer completed successfully")
	return nil
}

// Name 返回命令名称
func (c *ConsumerCommand) Name() string {
	return c.fs.Name()
}

// Usage 返回命令用法
func (c *ConsumerCommand) Usage() string {
	return "Start the consumer service"
}

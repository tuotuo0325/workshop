package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"airbnb-cli/internal/producer"
	"airbnb-cli/internal/queue"
)

// ProducerCommand 表示生产者命令
type ProducerCommand struct {
	fs        *flag.FlagSet
	dataFile  string
	queuePath string
}

// NewProducerCommand 创建生产者命令
func NewProducerCommand() *ProducerCommand {
	c := &ProducerCommand{
		fs: flag.NewFlagSet("producer", flag.ExitOnError),
	}

	c.fs.StringVar(&c.dataFile, "data", "tasks.json", "Tasks data file")
	c.fs.StringVar(&c.queuePath, "queue", "data/queue.data", "Queue server address")

	return c
}

// Run 运行生产者命令
func (c *ProducerCommand) Run(args []string) error {
	if err := c.fs.Parse(args); err != nil {
		return err
	}

	fmt.Printf("Starting producer with data file: %s, queue: %s\n", c.dataFile, c.queuePath)

	// 检查任务文件是否存在
	if _, err := os.Stat(c.dataFile); os.IsNotExist(err) {
		return fmt.Errorf("task file not found: %s", c.dataFile)
	}

	// 创建队列目录
	queueFile := c.queuePath
	if err := os.MkdirAll(filepath.Dir(queueFile), 0755); err != nil {
		return fmt.Errorf("create queue directory failed: %w", err)
	}

	fmt.Printf("Creating queue at: %s\n", queueFile)

	// 创建队列
	q, err := queue.NewFileQueue(queueFile)
	if err != nil {
		return fmt.Errorf("create queue failed: %w", err)
	}

	// 创建生产者
	p := producer.NewTaskProducer(c.dataFile, q)

	fmt.Printf("Starting to process tasks...\n")
	if err := p.Start(context.Background()); err != nil {
		return fmt.Errorf("producer failed: %w", err)
	}

	if err := p.Close(); err != nil {
		return fmt.Errorf("close producer failed: %w", err)
	}

	fmt.Println("Producer completed successfully")
	return nil
}

// Name 返回命令名称
func (c *ProducerCommand) Name() string {
	return c.fs.Name()
}

// Usage 返回命令用法
func (c *ProducerCommand) Usage() string {
	return "Start the producer service"
}

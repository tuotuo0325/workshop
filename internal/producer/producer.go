package producer

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"airbnb-cli/internal/model"
	"airbnb-cli/internal/queue"
)

// Producer 定义了生产者接口
type Producer interface {
	// Start 启动生产者
	Start(ctx context.Context) error
	// Close 关闭生产者
	Close() error
}

// TaskProducer 实现了任务生产者
type TaskProducer struct {
	dataFile string
	queue    queue.Queue
}

// NewTaskProducer 创建任务生产者
func NewTaskProducer(dataFile string, queue queue.Queue) Producer {
	return &TaskProducer{
		dataFile: dataFile,
		queue:    queue,
	}
}

func (p *TaskProducer) Start(ctx context.Context) error {
	// 读取任务文件
	data, err := os.ReadFile(p.dataFile)
	if err != nil {
		return fmt.Errorf("read task file failed: %w", err)
	}

	var taskList model.TaskList
	if err := json.Unmarshal(data, &taskList); err != nil {
		return fmt.Errorf("parse task file failed: %w", err)
	}

	// 将任务推送到队列
	for _, task := range taskList.Tasks {
		taskData, err := json.Marshal(task)
		if err != nil {
			return fmt.Errorf("marshal task failed: %w", err)
		}

		if err := p.queue.Push(ctx, taskData); err != nil {
			return fmt.Errorf("push task to queue failed: %w", err)
		}
	}

	return nil
}

func (p *TaskProducer) Close() error {
	return p.queue.Close()
}

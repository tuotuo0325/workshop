package storage

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"sync"

	"airbnb-cli/internal/model"
)

// Storage 定义了存储接口
type Storage interface {
	// SaveHotel 保存酒店信息
	SaveHotel(ctx context.Context, hotel *model.Hotel) error

	// Close 关闭存储
	Close() error
}

// FileStorage 实现了基于文件的存储
type FileStorage struct {
	path string
	mu   sync.Mutex
}

// NewFileStorage 创建基于文件的存储
func NewFileStorage(path string) (Storage, error) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	return &FileStorage{
		path: path,
	}, nil
}

func (s *FileStorage) SaveHotel(ctx context.Context, hotel *model.Hotel) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	f, err := os.OpenFile(s.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	encoder := json.NewEncoder(f)
	return encoder.Encode(hotel)
}

func (s *FileStorage) Close() error {
	return nil
}

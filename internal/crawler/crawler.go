package crawler

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"airbnb-cli/internal/model"
)

// Crawler 定义了爬虫接口
type Crawler interface {
	// Crawl 爬取指定URL的数据
	Crawl(ctx context.Context, task *model.Task) (*model.Hotel, error)
}

// AirbnbCrawler 实现了Airbnb爬虫
type AirbnbCrawler struct {
	client *http.Client
}

// NewAirbnbCrawler 创建Airbnb爬虫实例
func NewAirbnbCrawler() Crawler {
	return &AirbnbCrawler{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *AirbnbCrawler) Crawl(ctx context.Context, task *model.Task) (*model.Hotel, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", task.URL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	// 设置请求头
	for k, v := range task.Headers {
		req.Header.Set(k, v)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body failed: %w", err)
	}

	// TODO: 解析HTML内容，提取所需信息
	// 这里需要根据实际的Airbnb页面结构来实现
	// 可以使用goquery等库来解析HTML

	// 临时返回一个示例数据
	hotel := &model.Hotel{
		Name:           "Sample Hotel",
		Stars:          4.5,
		Price:          199.99,
		PriceBeforeTax: 180.00,
		CheckIn:        time.Now(),
		CheckOut:       time.Now().Add(24 * time.Hour),
		Guests:         2,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	return hotel, nil
}

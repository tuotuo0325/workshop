package crawler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"airbnb-cli/internal/model"

	"github.com/PuerkitoBio/goquery"
)

// Crawler 定义了爬虫接口
type Crawler interface {
	// Crawl 爬取指定URL的数据
	Crawl(ctx context.Context, task *model.Task) ([]*model.Hotel, error)
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

func (c *AirbnbCrawler) Crawl(ctx context.Context, task *model.Task) ([]*model.Hotel, error) {
	hotels, err := c.crawlListPage(ctx, task)
	if err != nil {
		return nil, fmt.Errorf("crawl list page failed: %w", err)
	}

	if len(hotels) == 0 {
		return nil, fmt.Errorf("no hotels found")
	}

	return hotels, nil
}

func (c *AirbnbCrawler) crawlListPage(ctx context.Context, task *model.Task) ([]*model.Hotel, error) {
	var hotels []*model.Hotel
	currentURL := task.URL

	for {
		// 检查上下文是否已取消
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		// 发送请求
		req, err := http.NewRequestWithContext(ctx, "GET", currentURL, nil)
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

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}

		// 解析HTML
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			resp.Body.Close()
			return nil, fmt.Errorf("parse HTML failed: %w", err)
		}
		resp.Body.Close()

		// 查找所有房源列表
		doc.Find("div.gsgwcjk div.itemListElement").Each(func(i int, s *goquery.Selection) {
			// 获取详情页URL
			detailURL, exists := s.Find("meta[itemprop='url']").Attr("content")
			if !exists {
				return
			}

			// 爬取详情页
			hotel, err := c.crawlDetailPage(ctx, detailURL, task.Headers)
			if err != nil {
				fmt.Printf("crawl detail page failed: %v\n", err)
				return
			}

			hotels = append(hotels, hotel)
		})

		// 查找下一页链接
		nextLink := doc.Find("a[aria-label='下一个']")
		if nextLink.Length() == 0 {
			break
		}

		nextURL, exists := nextLink.Attr("href")
		if !exists {
			break
		}

		currentURL = nextURL
	}

	return hotels, nil
}

func (c *AirbnbCrawler) crawlDetailPage(ctx context.Context, url string, headers map[string]string) (*model.Hotel, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	for k, v := range headers {
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

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("parse HTML failed: %w", err)
	}

	hotel := &model.Hotel{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 1. 获取酒店名称
	hotel.Name = doc.Find("div[data-plugin-in-point-id='OVERVIEW_DEFAULT_V2'] div.hpipapi h2").Text()

	// 2. 获取星级
	stars := doc.Find("div[data-plugin-in-point-id='GUEST_FAVORITE_BANNER'] div.a8jhwcl div[aria-hidden='true']").Text()
	if stars == "" {
		stars = doc.Find("div.r1lutz1s").Text()
	}
	if s, err := strconv.ParseFloat(strings.TrimSpace(stars), 32); err == nil {
		hotel.Stars = float32(s)
	}

	// 3. 获取价格
	if price := doc.Find("span._11jcbg2").Text(); price != "" {
		if p, err := strconv.ParseFloat(strings.TrimPrefix(price, "¥"), 64); err == nil {
			hotel.Price = p
		}
	}

	// 4. 获取税前价格
	if priceBeforeTax := doc.Find("div._1avmy66 span._j1kt73").Text(); priceBeforeTax != "" {
		if p, err := strconv.ParseFloat(strings.TrimPrefix(priceBeforeTax, "¥"), 64); err == nil {
			hotel.PriceBeforeTax = p
		}
	}

	// 5. 获取入住日期
	if checkIn := doc.Find("div[data-testid='change-dates-checkIn']").Text(); checkIn != "" {
		if t, err := time.Parse("2006-01-02", checkIn); err == nil {
			hotel.CheckIn = t
		}
	}

	// 6. 获取退房日期
	if checkOut := doc.Find("div[data-testid='change-dates-checkOut']").Text(); checkOut != "" {
		if t, err := time.Parse("2006-01-02", checkOut); err == nil {
			hotel.CheckOut = t
		}
	}

	// 7. 获取客人数量
	if guests := doc.Find("div._7pspom span._j1kt73").Text(); guests != "" {
		if g, err := strconv.Atoi(strings.TrimSuffix(guests, "位房客")); err == nil {
			hotel.Guests = g
		}
	}

	return hotel, nil
}

package model

import "time"

// Hotel 表示一个酒店的预订信息
type Hotel struct {
	ID             int64     `json:"id"`
	Name           string    `json:"name"`           // 酒店名称
	Stars          float32   `json:"stars"`          // 星级
	Price          float64   `json:"price"`          // 价格
	PriceBeforeTax float64   `json:"priceBeforeTax"` // 税前价格
	CheckIn        time.Time `json:"checkIn"`        // 入住日期
	CheckOut       time.Time `json:"checkOut"`       // 退房日期
	Guests         int       `json:"guests"`         // 客人数量
	CreatedAt      time.Time `json:"createdAt"`      // 记录创建时间
	UpdatedAt      time.Time `json:"updatedAt"`      // 记录更新时间
}

// Task 表示一个爬取任务
type Task struct {
	Name    string            `json:"name"`    // 任务名称
	URL     string            `json:"url"`     // 爬取URL
	Headers map[string]string `json:"headers"` // 请求头
}

// TaskList 表示任务列表
type TaskList struct {
	Tasks []Task `json:"tasks"`
}

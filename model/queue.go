package model

import (
	"time"
)

// Queue 消息队列
type Queue struct {
	ID         uint64    `gorm:"column:id;primary_key;AUTO_INCREMENT"`                 // 主键
	Data       string    `gorm:"column:data"`                                          // 消息
	LockID     string    `gorm:"column:lock_id;NOT NULL"`                              // 锁id
	RetryTimes uint      `gorm:"column:retry_times;NOT NULL"`                          // 重试次数
	CreatedAt  time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP;NOT NULL"` // 创建时间
	UpdatedAt  time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP;NOT NULL"` // 修改时间
}

func (m *Queue) TableName() string {
	return "queue"
}

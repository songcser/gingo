package model

import (
	"github.com/songcser/gingo/utils"
)

type Model interface {
	Get() int64
}

type BaseModel struct {
	ID        int64          `json:"id" gorm:"primarykey" admin:"disable"`                // 主键ID
	CreatedAt utils.JsonTime `json:"createdAt" gorm:"index;comment:创建时间" admin:"disable"` // 创建时间
	UpdatedAt utils.JsonTime `json:"updatedAt" gorm:"index;comment:更新时间" admin:"disable"` // 更新时间
}

func (m BaseModel) Get() int64 {
	return m.ID
}

package model

import (
	"github.com/songcser/gingo/utils"
)

type Model interface {
	Get() int64
}

type BaseModel struct {
	ID        int64          `json:"id" gorm:"primarykey"`                // 主键ID
	CreatedAt utils.JsonTime `gorm:"index;comment:创建时间" json:"createdAt"` // 创建时间
	UpdatedAt utils.JsonTime `gorm:"index;comment:更新时间" json:"updatedAt"` // 更新时间
}

func (m BaseModel) Get() int64 {
	return m.ID
}

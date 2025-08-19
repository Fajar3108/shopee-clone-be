package model

import "database/sql"

type Category struct {
	ID        string       `json:"id" gorm:"primaryKey;not null"`
	ParentID  *string      `json:"parent_id" gorm:"index"`
	Name      string       `json:"name" gorm:"not null;size:255"`
	Slug      string       `json:"slug" gorm:"unique;not null"`
	CreatedAt sql.NullTime `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt sql.NullTime `json:"updated_at" gorm:"autoUpdateTime"`
	Parent    *Category    `json:"parent" gorm:"foreignKey:ParentID"`
	Children  []*Category  `json:"children" gorm:"foreignKey:ParentID"`
}

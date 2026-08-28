package role

import (
	"go-boilerplate/internal/audit"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Name      string         `gorm:"size:100;not null" json:"name"`
	Slug      string         `gorm:"size:100;not null" json:"slug"`
	IsActive  bool           `gorm:"not null;default:true" json:"isActive"`
	audit.BaseModel
}

func (r *Role) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		id, err := uuid.NewV7()
		if err != nil {
			return err
		}
		r.ID = id
	}
	return r.BaseModel.BeforeCreate(tx)
}

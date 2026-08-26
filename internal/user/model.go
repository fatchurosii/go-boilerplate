package user

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"go-boilerplate/internal/role"
)

type User struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	RoleID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"roleId"`
	Username  string         `gorm:"size:100;not null" json:"username"`
	Password  string         `gorm:"not null" json:"-"`
	IsActive  bool           `gorm:"not null;default:true" json:"isActive"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Role      role.Role      `gorm:"foreignKey:RoleID" json:"role"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID != uuid.Nil {
		return nil
	}

	id, err := uuid.NewV7()
	if err != nil {
		return err
	}

	u.ID = id
	return nil
}

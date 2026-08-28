package audit

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type contextKey string

const UserIDKey contextKey = "userId"

func UserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value(UserIDKey).(uuid.UUID)
	return userID, ok
}

type AuditableModel struct {
	CreatedBy *uuid.UUID `json:"createdBy,omitempty"`
	UpdatedBy *uuid.UUID `json:"updatedBy,omitempty"`
	DeletedBy *uuid.UUID `json:"deletedBy,omitempty"`
}

func (b *BaseModel) BeforeCreate(tx *gorm.DB) error {
	if userID, ok := UserIDFromContext(tx.Statement.Context); ok {
		b.CreatedBy = &userID
		b.UpdatedBy = &userID
	}
	return nil
}

func (b *BaseModel) BeforeUpdate(tx *gorm.DB) error {
	if userID, ok := UserIDFromContext(tx.Statement.Context); ok {
		b.UpdatedBy = &userID
	}
	return nil
}

func (b *BaseModel) BeforeDelete(tx *gorm.DB) error {
	if userID, ok := UserIDFromContext(tx.Statement.Context); ok {
		tx.Statement.SetColumn("deleted_by", userID)
	}
	return nil
}

type BaseModel struct {
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	AuditableModel
}

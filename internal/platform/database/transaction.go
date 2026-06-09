package database

import (
	"context"

	"gorm.io/gorm"
)

func WithinTransaction(ctx context.Context, db *gorm.DB, fn func(tx *gorm.DB) error) error {
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(tx)
	})
}

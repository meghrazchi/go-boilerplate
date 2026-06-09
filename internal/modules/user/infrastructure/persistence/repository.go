package persistence

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/your-org/go-ddd-boilerplate/internal/modules/user/domain"
)

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) Create(ctx context.Context, user *domain.User) error {
	model := toModel(user)
	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		if isUniqueViolation(err) {
			return domain.ErrUserAlreadyExists
		}
		return err
	}
	return nil
}

func (r *GormRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var model UserModel
	if err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return toDomain(model)
}

func (r *GormRepository) FindByEmail(ctx context.Context, email domain.Email) (*domain.User, error) {
	var model UserModel
	if err := r.db.WithContext(ctx).First(&model, "email = ?", email.String()).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return toDomain(model)
}

func (r *GormRepository) List(ctx context.Context, params domain.ListParams) ([]*domain.User, int64, error) {
	query := r.db.WithContext(ctx).Model(&UserModel{})

	if params.Search != "" {
		search := "%" + strings.ToLower(params.Search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(email) LIKE ?", search, search)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var models []UserModel
	orderBy := fmt.Sprintf("%s %s", sanitizeSort(params.Sort), sanitizeOrder(params.Order))
	if err := query.Order(orderBy).Limit(params.Limit).Offset(params.Offset()).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	users, err := toDomainSlice(models)
	if err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (r *GormRepository) Update(ctx context.Context, user *domain.User) error {
	model := toModel(user)
	result := r.db.WithContext(ctx).Model(&UserModel{}).Where("id = ?", model.ID).Updates(map[string]any{
		"name":       model.Name,
		"email":      model.Email,
		"updated_at": model.UpdatedAt,
	})
	if result.Error != nil {
		if isUniqueViolation(result.Error) {
			return domain.ErrUserAlreadyExists
		}
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrUserNotFound
	}
	return nil
}

func (r *GormRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&UserModel{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrUserNotFound
	}
	return nil
}

func isUniqueViolation(err error) bool {
	return errors.Is(err, gorm.ErrDuplicatedKey) ||
		strings.Contains(strings.ToLower(err.Error()), "duplicate key") ||
		strings.Contains(strings.ToLower(err.Error()), "unique constraint")
}

func sanitizeSort(sort string) string {
	switch sort {
	case "name", "email", "created_at", "updated_at":
		return sort
	default:
		return "created_at"
	}
}

func sanitizeOrder(order string) string {
	if strings.EqualFold(order, "asc") {
		return "asc"
	}
	return "desc"
}

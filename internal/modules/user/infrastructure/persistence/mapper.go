package persistence

import "github.com/your-org/go-ddd-boilerplate/internal/modules/user/domain"

func toModel(user *domain.User) UserModel {
	return UserModel{
		ID:        user.ID(),
		Name:      user.Name(),
		Email:     user.Email().String(),
		CreatedAt: user.CreatedAt(),
		UpdatedAt: user.UpdatedAt(),
	}
}

func toDomain(model UserModel) (*domain.User, error) {
	email, err := domain.NewEmail(model.Email)
	if err != nil {
		return nil, err
	}
	return domain.RehydrateUser(model.ID, model.Name, email, model.CreatedAt, model.UpdatedAt), nil
}

func toDomainSlice(models []UserModel) ([]*domain.User, error) {
	users := make([]*domain.User, 0, len(models))
	for _, model := range models {
		user, err := toDomain(model)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

package repository

import (
	"github.com/event-management/api-gateway/internal/config"
	"github.com/event-management/api-gateway/internal/domain"
	"github.com/event-management/api-gateway/internal/repository/user"
)

// Repository contains all repositories
type Repository struct {
	User domain.UserRepository
}

// New creates a new repository
func New(cfg *config.Config) (*Repository, error) {
	// В реальном приложении здесь будет инициализация подключения к базе данных
	// и создание экземпляров репозиториев
	// В качестве заглушки используем in-memory реализацию
	userRepo := user.NewInMemoryRepository()

	return &Repository{
		User: userRepo,
	}, nil
}

// Close closes the repository connections
func (r *Repository) Close() error {
	// Закрытие соединений с базой данных
	return nil
}

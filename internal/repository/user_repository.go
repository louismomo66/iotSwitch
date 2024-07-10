package repository

import (
	"iot_switch/internal/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	GetUserByEmail(email string) (models.User, error)
	CreateUser(user *models.User) error
	GetUserEmail(email string) (string, error)
	GetAllUsers() ([]models.User, error)
	UpdatePasswordByEmail(email, hashedPassword string) error
}
type gormUserRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *gormUserRepo {
	return &gormUserRepo{db}
}

func (repo *gormUserRepo) GetUserByEmail(email string) (models.User, error) {
	var user models.User
	err := repo.db.Where("email = ?", email).First(&user).Error
	return user, err
}

// In your repository file
func (repo *gormUserRepo) GetAllUsers() ([]models.User, error) {
	var users []models.User
	err := repo.db.Find(&users).Error
	return users, err
}

func (repo *gormUserRepo) GetUserEmail(email string) (string, error) {
	var user models.User
	err := repo.db.Where("email=?", email).First(&user).Error
	if err != nil {
		return "", err
	}
	return user.Email, err
}
func (repo *gormUserRepo) CreateUser(user *models.User) error {
	return repo.db.Create(user).Error
}

func (repo *gormUserRepo) UpdatePasswordByEmail(email, hashedPassword string) error {
	return repo.db.Model(&models.User{}).Where("email = ?", email).Update("password", hashedPassword).Error
}

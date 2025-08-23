package userService

import (
	"POSTnGETtrain/internal/models"

	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetAll() ([]models.User, error) {
	args := m.Called()
	if res := args.Get(0); res != nil {
		return res.([]models.User), args.Error(1)
	}
	return []models.User{}, args.Error(1)
}
func (m *MockUserRepository) GetByID(id string) (*models.User, error) {
	args := m.Called(id)
	var user *models.User
	if res := args.Get(0); res != nil {
		user = res.(*models.User)
	}
	return user, args.Error(1)
}

func (m *MockUserRepository) Create(user *models.User) (*models.User, error) {
	args := m.Called(user)
	var u *models.User
	if res := args.Get(0); res != nil {
		u = res.(*models.User)
	}
	return u, args.Error(1)
}

func (m *MockUserRepository) Update(user *models.User) (*models.User, error) {
	args := m.Called(user)
	var u *models.User
	if res := args.Get(0); res != nil {
		u = res.(*models.User)
	}
	return u, args.Error(1)
}

func (m *MockUserRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepository) EmailExists(email string) (bool, error) {
	args := m.Called(email)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) GetTasksForUser(userID string) ([]models.Task, error) {
	args := m.Called(userID)
	if res := args.Get(0); res != nil {
		return res.([]models.Task), args.Error(1)
	}
	return []models.Task{}, args.Error(1)
}

func (m *MockUserRepository) GetUserWithTasks(userID string) (*models.User, error) {
	args := m.Called(userID)
	var user *models.User
	if res := args.Get(0); res != nil {
		user = res.(*models.User)
	}
	return user, args.Error(1)
}

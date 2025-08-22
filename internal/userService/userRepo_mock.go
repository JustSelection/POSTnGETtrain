package userService

import (
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetAll() ([]User, error) {
	args := m.Called()
	if res := args.Get(0); res != nil {
		return res.([]User), args.Error(1)
	}
	return []User{}, args.Error(1)
}
func (m *MockUserRepository) GetByID(id string) (*User, error) {
	args := m.Called(id)
	var user *User
	if res := args.Get(0); res != nil {
		user = res.(*User)
	}
	return user, args.Error(1)
}

func (m *MockUserRepository) Create(user *User) (*User, error) {
	args := m.Called(user)
	var u *User
	if res := args.Get(0); res != nil {
		u = res.(*User)
	}
	return u, args.Error(1)
}

func (m *MockUserRepository) Update(user *User) (*User, error) {
	args := m.Called(user)
	var u *User
	if res := args.Get(0); res != nil {
		u = res.(*User)
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

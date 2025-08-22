package userService

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateUser(t *testing.T) {
	tests := []struct {
		name      string
		email     string
		password  string
		mockSetup func(m *MockUserRepository, email, password string)
		wantErr   bool
		errMsg    string
	}{
		{
			name:     "успешное создание",
			email:    "test@mail.ru",
			password: "password123",
			mockSetup: func(m *MockUserRepository, email, password string) {
				m.On("Create", mock.MatchedBy(func(user *User) bool {
					return user.Email == email && user.Password == password
				})).Return(&User{
					ID:       "test-id",
					Email:    email,
					Password: password,
				}, nil)
			},
			wantErr: false,
		},
		{
			name:     "email уже существует",
			email:    "duplicate@mail.ru",
			password: "pass1",
			mockSetup: func(m *MockUserRepository, email, password string) {
				m.On("Create", mock.MatchedBy(func(user *User) bool {
					return user.Email == email && user.Password == password
				})).Return(nil, ErrEmailExists)
			},
			wantErr: true,
			errMsg:  "email already exists",
		},
		{
			name:     "ошибка создания пользователя",
			email:    "create-error@mail.ru",
			password: "pass1",
			mockSetup: func(m *MockUserRepository, email, password string) {
				m.On("Create", mock.AnythingOfType("*userService.User")).Return(nil, errors.New("create error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			tt.mockSetup(mockRepo, tt.email, tt.password)

			service := NewUserService(mockRepo)
			user, err := service.CreateUser(tt.email, tt.password)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.email, user.Email)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetUserByID(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		mockSetup func(m *MockUserRepository, id string)
		want      *User
		wantErr   bool
	}{
		{
			name: "успешное получение",
			id:   "user-id",
			mockSetup: func(m *MockUserRepository, id string) {
				m.On("GetByID", id).Return(&User{
					ID:       id,
					Email:    "user@mail.ru",
					Password: "secret",
				}, nil)
			},
			want: &User{
				ID:       "user-id",
				Email:    "user@mail.ru",
				Password: "secret",
			},
			wantErr: false,
		},
		{
			name: "ошибка получения",
			id:   "nonexistent-id",
			mockSetup: func(m *MockUserRepository, id string) {
				m.On("GetByID", id).Return(nil, ErrUserNotFound)
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			tt.mockSetup(mockRepo, tt.id)

			service := NewUserService(mockRepo)
			user, err := service.GetUserByID(tt.id)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, user)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetAllUsers(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func(m *MockUserRepository)
		want      []User
		wantErr   bool
	}{
		{
			name: "успешное получение всех юзеров",
			mockSetup: func(m *MockUserRepository) {
				m.On("GetAll").Return([]User{
					{ID: "1", Email: "alabay@gmail.com", Password: "111"},
					{ID: "2", Email: "barista@mail.ru", Password: "222"},
				}, nil)
			},
			want: []User{
				{ID: "1", Email: "alabay@gmail.com", Password: "111"},
				{ID: "2", Email: "barista@mail.ru", Password: "222"},
			},
			wantErr: false,
		},
		{
			name: "ошибка репозитория",
			mockSetup: func(m *MockUserRepository) {
				m.On("GetAll").Return(nil, errors.New("db error"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			tt.mockSetup(mockRepo)

			service := NewUserService(mockRepo)

			result, err := service.GetAllUsers()

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, result)
			}

			mockRepo.AssertExpectations(t)
		})

	}
}

func TestUpdateUser(t *testing.T) {
	newEmail := "new@mail.ru"
	newPass := "newpass"

	tests := []struct {
		name      string
		id        string
		email     *string
		password  *string
		mockSetup func(m *MockUserRepository, id string)
		want      *User
		wantErr   bool
	}{
		{
			name:     "успешное обновление",
			id:       "user-id",
			email:    &newEmail,
			password: &newPass,
			mockSetup: func(m *MockUserRepository, id string) {
				m.On("GetByID", id).Return(&User{
					ID:       id,
					Email:    "old@mail.ru",
					Password: "oldpass",
				}, nil)
				m.On("Update", mock.MatchedBy(func(user *User) bool {
					return user.ID == id && user.Email == newEmail && user.Password == newPass
				})).Return(&User{
					ID:       id,
					Email:    newEmail,
					Password: newPass,
				}, nil)
			},
			want: &User{
				ID:       "user-id",
				Email:    "new@mail.ru",
				Password: "newpass",
			},
			wantErr: false,
		},
		{
			name:     "ошибка получения пользователя",
			id:       "not_found",
			email:    &newEmail,
			password: &newPass,
			mockSetup: func(m *MockUserRepository, id string) {
				m.On("GetByID", id).Return(nil, ErrUserNotFound)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:     "обновление только email",
			id:       "user-id",
			email:    &newEmail,
			password: nil,
			mockSetup: func(m *MockUserRepository, id string) {
				m.On("GetByID", id).Return(&User{
					ID:       id,
					Email:    "old@mail.ru",
					Password: "oldpass",
				}, nil)
				m.On("Update", mock.MatchedBy(func(user *User) bool {
					return user.ID == id && user.Email == newEmail && user.Password == "oldpass"
				})).Return(&User{
					ID:       id,
					Email:    newEmail,
					Password: "oldpass",
				}, nil)
			},
			want: &User{
				ID:       "user-id",
				Email:    "new@mail.ru",
				Password: "oldpass",
			},
			wantErr: false,
		},
		{
			name:     "обновление только пароля",
			id:       "user-id",
			password: &newPass,
			mockSetup: func(m *MockUserRepository, id string) {
				m.On("GetByID", id).Return(&User{
					ID:       id,
					Email:    "old@mail.ru",
					Password: "oldpass",
				}, nil)
				m.On("Update", mock.MatchedBy(func(user *User) bool {
					return user.ID == id && user.Email == "old@mail.ru" && user.Password == newPass
				})).Return(&User{
					ID:       id,
					Email:    "old@mail.ru",
					Password: newPass,
				}, nil)
			},
			want: &User{
				ID:       "user-id",
				Email:    "old@mail.ru",
				Password: "newpass",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			tt.mockSetup(mockRepo, tt.id)

			service := NewUserService(mockRepo)
			result, err := service.UpdateUser(tt.id, tt.email, tt.password)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, result)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestDeleteUser(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		mockSetup func(m *MockUserRepository, id string)
		wantErr   bool
	}{
		{
			name: "успешное удаление",
			id:   "user-id",
			mockSetup: func(m *MockUserRepository, id string) {
				m.On("Delete", id).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "ошибка удаления",
			id:   "not_found",
			mockSetup: func(m *MockUserRepository, id string) {
				m.On("Delete", id).Return(ErrUserNotFound)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			tt.mockSetup(mockRepo, tt.id)

			service := NewUserService(mockRepo)

			err := service.DeleteUser(tt.id)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestEmailExists(t *testing.T) {
	tests := []struct {
		name      string
		email     string
		mockSetup func(m *MockUserRepository, email string)
		want      bool
		wantErr   bool
	}{
		{
			name:  "email существует",
			email: "exists@mail.ru",
			mockSetup: func(m *MockUserRepository, email string) {
				m.On("EmailExists", email).Return(true, nil)
			},
			want:    true,
			wantErr: false,
		},
		{
			name:  "email не существует",
			email: "notexists@mail.ru",
			mockSetup: func(m *MockUserRepository, email string) {
				m.On("EmailExists", email).Return(false, nil)
			},
			want:    false,
			wantErr: false,
		},
		{
			name:  "ошибка проверки email",
			email: "error@mail.ru",
			mockSetup: func(m *MockUserRepository, email string) {
				m.On("EmailExists", email).Return(false, errors.New("db error"))
			},
			want:    false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			tt.mockSetup(mockRepo, tt.email)

			result, err := mockRepo.EmailExists(tt.email)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, result)

				mockRepo.AssertExpectations(t)
			}
		})
	}
}

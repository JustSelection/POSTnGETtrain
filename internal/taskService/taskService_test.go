package taskService

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateTask(t *testing.T) {
	tests := []struct {
		name      string
		input     Task
		mockSetup func(m *MockTaskRepository, input Task)
		wantErr   bool
	}{
		{
			name:  "успешное создание",
			input: Task{Name: "Test Task", IsDone: false},
			mockSetup: func(m *MockTaskRepository, input Task) {
				m.On("Create", mock.MatchedBy(func(t Task) bool {
					return t.Name == input.Name && t.IsDone == input.IsDone
				})).Return(input, nil)
			},
			wantErr: false,
		},
		{
			name:  "ошибка создания",
			input: Task{Name: "Bad Task", IsDone: false},
			mockSetup: func(m *MockTaskRepository, input Task) {
				m.On("Create", mock.AnythingOfType("taskService.Task")).Return(Task{},
					errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockTaskRepository)
			tt.mockSetup(mockRepo, tt.input)

			service := NewTaskService(mockRepo)
			_, err := service.CreateTask(tt.input.Name, tt.input.IsDone)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetAllTasks(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func(m *MockTaskRepository)
		want      []Task
		wantErr   bool
	}{
		{
			name: "успешное получение всех задач",
			mockSetup: func(m *MockTaskRepository) {
				m.On("GetAll").Return([]Task{
					{ID: "1", Name: "Task 1", IsDone: false},
					{ID: "2", Name: "Task 2", IsDone: true},
				}, nil)
			},
			want: []Task{
				{ID: "1", Name: "Task 1", IsDone: false},
				{ID: "2", Name: "Task 2", IsDone: true},
			},
			wantErr: false,
		},
		{
			name: "ошибка репозитория",
			mockSetup: func(m *MockTaskRepository) {
				m.On("GetAll").Return(nil, errors.New("db error"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockTaskRepository)
			tt.mockSetup(mockRepo)

			service := NewTaskService(mockRepo)
			result, err := service.GetAllTasks()

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

func TestGetTaskByID(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		mockSetup func(m *MockTaskRepository, id string)
		want      Task
		wantErr   bool
	}{
		{
			name: "успешное получение",
			id:   "1",
			mockSetup: func(m *MockTaskRepository, id string) {
				m.On("GetByID", id).Return(Task{
					ID: id, Name: "Test Task", IsDone: false}, nil)

			},
			want:    Task{ID: "1", Name: "Test Task", IsDone: false},
			wantErr: false,
		},
		{
			name: "ошибка получения",
			id:   "99",
			mockSetup: func(m *MockTaskRepository, id string) {
				m.On("GetByID", id).Return(Task{}, errors.New("not found"))
			},
			want:    Task{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockTaskRepository)
			tt.mockSetup(mockRepo, tt.id)

			service := NewTaskService(mockRepo)
			result, err := service.GetTaskByID(tt.id)

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

func TestUpdateTask(t *testing.T) {
	name := "Updated"
	isDone := true

	tests := []struct {
		name      string
		id        string
		newName   *string
		newDone   *bool
		mockSetup func(m *MockTaskRepository, id string, existing Task, updated Task)
		want      Task
		wantErr   bool
	}{
		{
			name:    "успешное обновление",
			id:      "1",
			newName: &name,
			newDone: &isDone,
			mockSetup: func(m *MockTaskRepository, id string, existing Task, updated Task) {
				m.On("GetByID", id).Return(existing, nil)
				m.On("Update", updated).Return(updated, nil)
			},
			want:    Task{ID: "1", Name: "Updated", IsDone: true},
			wantErr: false,
		},
		{
			name: "ошибка получения задачи",
			id:   "99",
			mockSetup: func(m *MockTaskRepository, id string, existing Task, updated Task) {
				m.On("GetByID", id).Return(Task{}, errors.New("not found"))
			},
			want:    Task{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockTaskRepository)

			existing := Task{ID: tt.id, Name: "Old", IsDone: false}
			updated := existing
			if tt.newName != nil {
				updated.Name = *tt.newName
			}
			if tt.newDone != nil {
				updated.IsDone = *tt.newDone
			}

			tt.mockSetup(mockRepo, tt.id, existing, updated)

			service := NewTaskService(mockRepo)
			result, err := service.UpdateTask(tt.id, tt.newName, tt.newDone)

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

func TestDeleteTask(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		mockSetup func(m *MockTaskRepository, id string)
		wantErr   bool
	}{
		{
			name: "успешное удаление",
			id:   "1",
			mockSetup: func(m *MockTaskRepository, id string) {
				m.On("Delete", id).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "ошибка удаления",
			id:   "2",
			mockSetup: func(m *MockTaskRepository, id string) {
				m.On("Delete", id).Return(errors.New("delete error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockTaskRepository)
			tt.mockSetup(mockRepo, tt.id)

			service := NewTaskService(mockRepo)
			err := service.DeleteTask(tt.id)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

package taskService

import (
	"POSTnGETtrain/internal/models"

	"github.com/stretchr/testify/mock"
)

type MockTaskRepository struct {
	mock.Mock
}

func (m *MockTaskRepository) Create(task models.Task) (models.Task, error) {
	args := m.Called(task)              // записываем вызов метода с аргументом task
	var t models.Task                   // переменная для возврата
	if res := args.Get(0); res != nil { // получим первый элемент с его проверкой
		t = res.(models.Task) // res интерфейс{} преобразуется в тип Task
	}
	return t, args.Error(1)
}

func (m *MockTaskRepository) GetAll() ([]models.Task, error) {
	args := m.Called()                  // Фиксируем вызов без аргументов
	if res := args.Get(0); res != nil { // проверяем первый возвращаемый аргумент
		return res.([]models.Task), args.Error(1) // res интерфейс{} преобразуется в тип Task
	}
	return []models.Task{}, args.Error(1) // если nil, возвращаем пустой слайс
}

func (m *MockTaskRepository) GetByID(id string) (models.Task, error) {
	args := m.Called(id)                // вызываем метод с аргументом айди
	var t models.Task                   // создаем переменную для результата
	if res := args.Get(0); res != nil { // проверяем первый возвращаемый аргумент
		t = res.(models.Task) // res интерфейс{} преобразуется в тип Task
	}
	return t, args.Error(1)
}

func (m *MockTaskRepository) GetByUserID(userID string) ([]models.Task, error) {
	args := m.Called(userID)
	if res := args.Get(0); res != nil {
		return res.([]models.Task), args.Error(1)
	}
	return []models.Task{}, args.Error(1)
}

func (m *MockTaskRepository) Update(task models.Task) (models.Task, error) {
	args := m.Called(task) // вызов с аргументом task
	var t models.Task
	if res := args.Get(0); res != nil {
		t = res.(models.Task) // res интерфейс{} преобразуется в тип Task
	}
	return t, args.Error(1)
}

func (m *MockTaskRepository) Delete(id string) error {
	args := m.Called(id) // фиксируем вызов с аргументом id
	return args.Error(0)
}

package taskService

import "github.com/stretchr/testify/mock"

type MockTaskRepository struct {
	mock.Mock
}

func (m *MockTaskRepository) Create(task Task) (Task, error) {
	args := m.Called(task)              // записываем вызов метода с аргументом task
	var t Task                          // переменная для возврата
	if res := args.Get(0); res != nil { // получим первый элемент с его проверкой
		t = res.(Task) // res интерфейс{} преобразуется в тип Task
	}
	return t, args.Error(1)
}

func (m *MockTaskRepository) GetAll() ([]Task, error) {
	args := m.Called()                  // Фиксируем вызов без аргументов
	if res := args.Get(0); res != nil { // проверяем первый возвращаемый аргумент
		return res.([]Task), args.Error(1) // res интерфейс{} преобразуется в тип Task
	}
	return []Task{}, args.Error(1) // если nil, возвращаем пустой слайс
}

func (m *MockTaskRepository) GetByID(id string) (Task, error) {
	args := m.Called(id)                // вызываем метод с аргументом айди
	var t Task                          // создаем переменную для результата
	if res := args.Get(0); res != nil { // проверяем первый возвращаемый аргумент
		t = res.(Task) // res интерфейс{} преобразуется в тип Task
	}
	return t, args.Error(1)
}

func (m *MockTaskRepository) Update(task Task) (Task, error) {
	args := m.Called(task) // вызов с аргументом task
	var t Task
	if res := args.Get(0); res != nil {
		t = res.(Task) // res интерфейс{} преобразуется в тип Task
	}
	return t, args.Error(1)
}

func (m *MockTaskRepository) Delete(id string) error {
	args := m.Called(id) // фиксируем вызов с аргументом id
	return args.Error(0)
}

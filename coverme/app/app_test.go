package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/slon/shad-go/coverme/models"
)

type mockStorage struct {
	todos []*models.Todo
}

func (m *mockStorage) GetAll() ([]*models.Todo, error) {
	return m.todos, nil
}

func (m *mockStorage) AddTodo(title, content string) (*models.Todo, error) {
	todo := models.Todo{ID: models.ID(len(m.todos) + 1), Title: title, Content: content, Finished: false}
	m.todos = append(m.todos, &todo)
	return &todo, nil
}

func (m *mockStorage) GetTodo(id models.ID) (*models.Todo, error) {
	for _, todo := range m.todos {
		if todo.ID == id {
			return todo, nil
		}
	}
	return &models.Todo{}, fmt.Errorf("not found")
}

func (m *mockStorage) FinishTodo(id models.ID) error {
	for i, todo := range m.todos {
		if todo.ID == id {
			m.todos[i].Finished = true
			return nil
		}
	}
	return fmt.Errorf("not found")
}

type brokenMockStorage struct {
	todos []*models.Todo
}

func (m *brokenMockStorage) GetAll() ([]*models.Todo, error) {
	return m.todos, fmt.Errorf("lol")
}

func (m *brokenMockStorage) AddTodo(title, content string) (*models.Todo, error) {
	return nil, fmt.Errorf("lol")
}

func (m *brokenMockStorage) GetTodo(id models.ID) (*models.Todo, error) {
	return &models.Todo{}, fmt.Errorf("not found")
}

func (m *brokenMockStorage) FinishTodo(id models.ID) error {
	return fmt.Errorf("not found")
}

func TestList(t *testing.T) {
	// Create a new mock storage with some test data
	storage := &mockStorage{
		todos: []*models.Todo{
			{ID: 1, Title: "Test Todo 1", Content: "Some content", Finished: false},
			{ID: 2, Title: "Test Todo 2", Content: "Some content", Finished: true},
			{ID: 3, Title: "Test Todo 3", Content: "Some content", Finished: false},
		},
	}
	brokenStorage := &brokenMockStorage{
		todos: []*models.Todo{
			{ID: 1, Title: "Test Todo 1", Content: "Some content", Finished: false},
			{ID: 2, Title: "Test Todo 2", Content: "Some content", Finished: true},
			{ID: 3, Title: "Test Todo 3", Content: "Some content", Finished: false},
		},
	}

	// Create a new app with the mock storage
	a := New(storage)
	b := New(brokenStorage)

	// Create a new request to the /todo endpoint
	req, err := http.NewRequest("GET", "/todo", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a new recorder to capture the response
	rr := httptest.NewRecorder()

	// Call the list function
	handler := http.HandlerFunc(a.list)
	handler.ServeHTTP(rr, req)

	// Check the response status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Unmarshal the response body into a slice of todos
	var todos []models.Todo
	err = json.Unmarshal(rr.Body.Bytes(), &todos)
	if err != nil {
		t.Fatal(err)
	}

	// Check that the number of todos in the response is equal to the number of todos in the storage
	assert.Equal(t, len(storage.todos), len(todos))

	// Check that the response todos match the storage todos
	for i, todo := range storage.todos {
		assert.Equal(t, todo.ID, todos[i].ID)
		assert.Equal(t, todo.Title, todos[i].Title)
		assert.Equal(t, todo.Content, todos[i].Content)
		assert.Equal(t, todo.Finished, todos[i].Finished)
	}

	br := httptest.NewRecorder()

	handler = b.list
	handler.ServeHTTP(br, req)

	assert.Equal(t, http.StatusInternalServerError, br.Code)
}

func TestStatus(t *testing.T) {
	storage := &mockStorage{
		todos: []*models.Todo{
			{ID: 1, Title: "Test Todo 1", Content: "Some content", Finished: false},
		},
	}
	a := New(storage)
	req, err := http.NewRequest("GET", "/todo", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(a.status)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	data, _ := io.ReadAll(rr.Body)
	assert.Equal(t, string(data), "\"API is up and working!\"")
}

func TestAddTodo(t *testing.T) {
	storage := &mockStorage{
		todos: []*models.Todo{
			{ID: 1, Title: "Test Todo 1", Content: "Some content", Finished: false},
			{ID: 2, Title: "", Content: "Some content", Finished: false},
		},
	}
	brokenStorage := &brokenMockStorage{
		todos: []*models.Todo{
			{ID: 1, Title: "Test Todo 1", Content: "Some content", Finished: false},
		},
	}
	a := New(storage)
	b := New(brokenStorage)
	mar, _ := json.Marshal(storage.todos[0])
	req, err := http.NewRequest("POST", "/todo/create", bytes.NewReader(mar))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(a.addTodo)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	mar, _ = json.Marshal(storage.todos[1])
	req, err = http.NewRequest("POST", "/todo/create", bytes.NewReader(mar))
	if err != nil {
		t.Fatal(err)
	}

	br := httptest.NewRecorder()

	handler = a.addTodo
	handler.ServeHTTP(br, req)

	assert.Equal(t, http.StatusBadRequest, br.Code)

	mar, _ = json.Marshal("lol")
	req, err = http.NewRequest("POST", "/todo/create", bytes.NewReader(mar))
	if err != nil {
		t.Fatal(err)
	}

	bbr := httptest.NewRecorder()

	handler = a.addTodo
	handler.ServeHTTP(bbr, req)

	assert.Equal(t, http.StatusBadRequest, bbr.Code)

	mar, _ = json.Marshal(brokenStorage.todos[0])
	req, err = http.NewRequest("POST", "/todo/create", bytes.NewReader(mar))
	if err != nil {
		t.Fatal(err)
	}

	bbbr := httptest.NewRecorder()

	handler = b.addTodo
	handler.ServeHTTP(bbbr, req)

	assert.Equal(t, http.StatusInternalServerError, bbbr.Code)
}

func TestGetTodo(t *testing.T) {
	todo := &models.Todo{ID: 1,
		Title:    "lol",
		Finished: false}

	var todos []*models.Todo
	todos = append(todos, todo)

	app := &App{
		db: &mockStorage{todos},
	}

	app.initRoutes()

	tests := []struct {
		name         string
		id           string
		expectStatus int
		expectBody   string
	}{
		{
			name:         "valid ID",
			id:           "1",
			expectStatus: http.StatusOK,
			expectBody:   `{"id":1,"title":"lol","content":"","finished":false}`,
		},
		{
			name:         "invalid ID",
			id:           "invalid",
			expectStatus: http.StatusNotFound,
			expectBody:   "404 page not found\n",
		},
		{
			name:         "todo not found",
			id:           "999",
			expectStatus: http.StatusInternalServerError,
			expectBody:   `Server encountered an error.`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, "/todo/"+tt.id, nil)
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}

			rr := httptest.NewRecorder()
			app.router.ServeHTTP(rr, req)

			if rr.Code != tt.expectStatus {
				t.Errorf("handler returned wrong status code: got %v, want %v",
					rr.Code, tt.expectStatus)
			}

			if body := rr.Body.String(); body != tt.expectBody {
				t.Errorf("handler returned unexpected body: got %v, want %v",
					body, tt.expectBody)
			}
		})
	}
}

func TestGetTodoNil(t *testing.T) {
	storage := &mockStorage{
		todos: []*models.Todo{
			{ID: 1, Title: "Test Todo 1", Content: "Some content", Finished: false},
			{ID: 2, Title: "", Content: "Some content", Finished: false},
		},
	}
	a := New(storage)
	req, err := http.NewRequest("GET", "/todo/lol", nil)
	if err != nil {
		t.Fatal(err)
	}
	g := httptest.NewRecorder()

	handler := http.HandlerFunc(a.getTodo)
	handler.ServeHTTP(g, req)

	assert.Equal(t, http.StatusBadRequest, g.Code)
}

func TestFinishTodoNil(t *testing.T) {
	storage := &mockStorage{
		todos: []*models.Todo{
			{ID: 1, Title: "Test Todo 1", Content: "Some content", Finished: false},
			{ID: 2, Title: "", Content: "Some content", Finished: false},
		},
	}
	a := New(storage)
	req, err := http.NewRequest("POST", "/todo/lol/finish", nil)
	if err != nil {
		t.Fatal(err)
	}
	g := httptest.NewRecorder()

	handler := http.HandlerFunc(a.finishTodo)
	handler.ServeHTTP(g, req)

	assert.Equal(t, http.StatusBadRequest, g.Code)

	a.initRoutes()

	req, err = http.NewRequest(http.MethodPost, "/todo/1/finish", nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	rr := httptest.NewRecorder()
	a.router.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, storage.todos[0].Finished, true)

	req, err = http.NewRequest(http.MethodPost, "/todo/3/finish", nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	rrr := httptest.NewRecorder()
	a.router.ServeHTTP(rrr, req)

	assert.Equal(t, rrr.Code, http.StatusInternalServerError)
}

func TestRun(t *testing.T) {
	storage := &mockStorage{
		todos: []*models.Todo{
			{ID: 1, Title: "Test Todo 1", Content: "Some content", Finished: false},
			{ID: 2, Title: "", Content: "Some content", Finished: false},
		},
	}
	a := New(storage)
	a.run("8081")
}

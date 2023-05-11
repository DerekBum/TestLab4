package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/slon/shad-go/coverme/models"
)

func TestAdd(t *testing.T) {
	client := New("8081")
	_, err := client.Add(&models.AddRequest{Title: "lol", Content: "kek"})
	assert.NotNil(t, err)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("unexpected method %s, expected POST", r.Method)
		}

		if r.URL.Path != "/todo/create" {
			t.Errorf("unexpected path %s, expected /todo/create", r.URL.Path)
		}

		expectedBody := &models.AddRequest{
			Title:   "lol",
			Content: "kek",
		}

		var actualBody models.AddRequest
		err = json.NewDecoder(r.Body).Decode(&actualBody)
		if err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}

		assert.Equal(t, actualBody, *expectedBody)

		resp := &models.Todo{
			ID:      1,
			Title:   "lol",
			Content: "kek",
		}
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client = New(server.URL)
	todo, err := client.Add(&models.AddRequest{Title: "lol", Content: "kek"})
	assert.Nil(t, err)
	assert.Equal(t, todo.Title, "lol")
	assert.Equal(t, todo.Content, "kek")

	serverError := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer serverError.Close()

	client = New(serverError.URL)
	_, err = client.Add(&models.AddRequest{Title: "lol", Content: "kek"})
	assert.NotNil(t, err)
}

func TestGet(t *testing.T) {
	client := New("8081")
	_, err := client.Get(1)
	assert.NotNil(t, err)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("unexpected method %s, expected GET", r.Method)
		}

		if r.URL.Path != "/todo/1" {
			t.Errorf("unexpected path %s, expected /todo/1", r.URL.Path)
		}

		resp := &models.Todo{
			ID:      1,
			Title:   "lol",
			Content: "kek",
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client = New(server.URL)
	todo, err := client.Get(1)
	assert.Nil(t, err)
	assert.Equal(t, todo.Title, "lol")
	assert.Equal(t, todo.Content, "kek")

	serverError := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer serverError.Close()

	client = New(serverError.URL)
	_, err = client.Get(1)
	assert.NotNil(t, err)
}

func TestList(t *testing.T) {
	client := New("8081")
	_, err := client.List()
	assert.NotNil(t, err)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("unexpected method %s, expected GET", r.Method)
		}

		if r.URL.Path != "/todo" {
			t.Errorf("unexpected path %s, expected /todo", r.URL.Path)
		}

		resp := make([]*models.Todo, 2)

		resp[0] = &models.Todo{
			ID:      1,
			Title:   "lol",
			Content: "kek",
		}
		resp[1] = &models.Todo{
			ID:      2,
			Title:   "kek",
			Content: "lol",
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client = New(server.URL)
	todo, err := client.List()
	assert.Nil(t, err)
	assert.Equal(t, len(todo), 2)
	assert.Equal(t, todo[0], &models.Todo{
		ID:      1,
		Title:   "lol",
		Content: "kek",
	})
	assert.Equal(t, todo[1], &models.Todo{
		ID:      2,
		Title:   "kek",
		Content: "lol",
	})

	serverError := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer serverError.Close()

	client = New(serverError.URL)
	_, err = client.List()
	assert.NotNil(t, err)
}

func TestFinish(t *testing.T) {
	client := New("8081")
	err := client.Finish(models.ID(1))
	assert.NotNil(t, err)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("unexpected method %s, expected POST", r.Method)
		}

		if r.URL.Path != "/todo/1/finish" {
			t.Errorf("unexpected path %s, expected /todo/1/finish", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client = New(server.URL)
	err = client.Finish(models.ID(1))
	assert.Nil(t, err)

	serverError := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer serverError.Close()

	client = New(serverError.URL)
	err = client.Finish(models.ID(1))
	assert.NotNil(t, err)
}

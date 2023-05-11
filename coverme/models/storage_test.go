package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddTodo(t *testing.T) {
	s := NewInMemoryStorage()
	todo, _ := s.AddTodo("lol", "kek")
	assert.Equal(t, todo.Title, "lol")
	assert.Equal(t, todo.Content, "kek")
	assert.Equal(t, todo.ID, ID(0))
}

func TestGetTodo(t *testing.T) {
	s := NewInMemoryStorage()
	todo, _ := s.AddTodo("lol", "kek")
	todoGet, err := s.GetTodo(ID(0))
	assert.Nil(t, err)
	assert.Equal(t, todo, todoGet)
	_, err = s.GetTodo(ID(1))
	assert.NotNil(t, err)
}

func TestGetAll(t *testing.T) {
	s := NewInMemoryStorage()
	todo1, _ := s.AddTodo("lol", "kek")
	todo2, _ := s.AddTodo("kek", "lol")
	todos, _ := s.GetAll()
	assert.Equal(t, len(todos), 2)
	assert.Equal(t, todos[0], todo1)
	assert.Equal(t, todos[1], todo2)
}

func TestFinishTodo(t *testing.T) {
	s := NewInMemoryStorage()
	todo, _ := s.AddTodo("lol", "kek")
	err := s.FinishTodo(ID(0))
	assert.Equal(t, todo.Finished, true)
	assert.Nil(t, err)
	err = s.FinishTodo(ID(1))
	assert.NotNil(t, err)
}

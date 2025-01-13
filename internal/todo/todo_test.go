package todo_test

import (
	"context"
	"github.com/iJosef/go-todo-api/internal/db"
	"github.com/iJosef/go-todo-api/internal/todo"
	"reflect"
	"testing"
)

type MockDB struct {
	items []db.Item
}

func (m *MockDB) DeleteItem(_ context.Context, id int) error {
	m.items = append(m.items[:id], m.items[id+1:]...)
	return nil
}

func (m *MockDB) InsertItem(_ context.Context, item db.Item) error {
	m.items = append(m.items, item)
	return nil
}

func (m *MockDB) GetAllItems(_ context.Context) ([]db.Item, error) {
	return m.items, nil
}

func TestService_Search(t *testing.T) {
	tests := []struct {
		name           string
		toDosToAdd     []string
		query          string
		expectedResult []string
	}{
		{
			name:           "given a todo of shop and a search of sh, i should get shop back",
			toDosToAdd:     []string{"shop"},
			query:          "sh",
			expectedResult: []string{"shop"},
		},
		{
			name:           "still returns shop, even if the case doesnt match",
			toDosToAdd:     []string{"Shopping"},
			query:          "sh",
			expectedResult: []string{"Shopping"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MockDB{}
			svc := todo.NewService(m)
			for _, toAdd := range tt.toDosToAdd {
				err := svc.Add(toAdd)
				if err != nil {
					t.Error(err)
				}
			}
			got, err := svc.Search(tt.query)
			if err != nil {
				t.Error(err)
			}
			if !reflect.DeepEqual(got, tt.expectedResult) {
				t.Errorf("Search() = %v, want %v", got, tt.expectedResult)
			}
		})
	}
}

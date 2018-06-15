package main_test

import (
	"testing"

	following "github.com/rakateja/di-go/manual-di"
)

type MockRepository struct{}

func NewMockRepository() following.Repository {
	return MockRepository{}
}

func (mock MockRepository) Store(entity following.Following) error {
	return nil
}

func TestInsertFollowing(t *testing.T) {
	followingRepository := NewMockRepository()
	followingService := following.NewService(followingRepository)
	if err := followingService.Insert("root", "I'm Root"); err != nil {
		t.Errorf("Got %v, expect nil  when inserting new following", err.Error())
	}
}

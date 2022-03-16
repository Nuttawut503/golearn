package store

import (
	"errors"
	"gographql/graph/model"
	"strconv"

	"github.com/google/uuid"
)

type UserStore struct {
	users []*model.User
}

func NewUserStore() *UserStore {
	return &UserStore{}
}

func (s *UserStore) Save(name string, age int) *model.User {
	id := uuid.New().ID()
	user := &model.User{
		ID:   strconv.FormatUint(uint64(id), 10),
		Name: name,
		Age:  age,
	}
	s.users = append(s.users, user)
	return user
}

func (s *UserStore) Delete(id string) (*model.User, error) {
	pos := -1
	for i, v := range s.users {
		if v.ID == id {
			pos = i
			break
		}
	}
	if pos == -1 {
		return nil, errors.New("user id not found")
	}
	user := s.users[pos]
	s.users = append(s.users[:pos], s.users[pos+1:]...)
	return user, nil
}

func (s *UserStore) FindByID(id string) (*model.User, error) {
	for _, v := range s.users {
		if v.ID == id {
			return v, nil
		}
	}
	return nil, errors.New("user id not found")
}

func (s *UserStore) GetAllUsers() []*model.User {
	return s.users
}

package graph

import (
	"gographql/graph/model"
	"gographql/store"
	"sync"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	mu              sync.Mutex
	UserSubscribers map[string]chan *model.UserUpdated
	UserStore       *store.UserStore
}

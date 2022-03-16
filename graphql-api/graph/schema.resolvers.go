package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"gographql/graph/generated"
	"gographql/graph/model"

	"github.com/google/uuid"
)

func (r *mutationResolver) CreateUser(ctx context.Context, input model.CreateUserInput) (*model.User, error) {
	user := r.UserStore.Save(input.Name, input.Age)
	r.mu.Lock()
	for _, v := range r.UserSubscribers {
		v <- &model.UserUpdated{
			Info: "created",
			User: user,
		}
	}
	r.mu.Unlock()
	return user, nil
}

func (r *mutationResolver) DeleteUser(ctx context.Context, id string) (*model.User, error) {
	user, err := r.UserStore.Delete(id)
	if err != nil {
		return nil, err
	}
	r.mu.Lock()
	for _, v := range r.UserSubscribers {
		v <- &model.UserUpdated{
			Info: "deleted",
			User: user,
		}
	}
	r.mu.Unlock()
	return user, nil
}

func (r *queryResolver) Users(ctx context.Context) ([]*model.User, error) {
	return r.UserStore.GetAllUsers(), nil
}

func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
	return r.UserStore.FindByID(id)
}

func (r *subscriptionResolver) UserUpdated(ctx context.Context) (<-chan *model.UserUpdated, error) {
	subscriberID := uuid.New().String()
	ch := make(chan *model.UserUpdated, 1)
	r.mu.Lock()
	r.UserSubscribers[subscriberID] = ch
	r.mu.Unlock()

	go func() {
		<-ctx.Done()
		r.mu.Lock()
		delete(r.UserSubscribers, subscriberID)
		r.mu.Unlock()
	}()

	return ch, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }

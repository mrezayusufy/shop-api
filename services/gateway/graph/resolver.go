package graph

import "github.com/mrezayusufy/shop-api/services/gateway"

type Resolver struct {
    *gateway.Resolver
}

type queryResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type userResolver struct{ *Resolver }
type orderResolver struct{ *Resolver }
type orderItemResolver struct{ *Resolver }

func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }
func (r *Resolver) User() UserResolver { return &userResolver{r} }
func (r *Resolver) Order() OrderResolver { return &orderResolver{r} }
func (r *Resolver) OrderItem() OrderItemResolver { return &orderItemResolver{r} }

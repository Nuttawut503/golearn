package server

import (
	"context"
	"grpc-api/customer/server/customerpb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Customer struct {
	ID   string
	Name string
}

var _customers = []Customer{{ID: "A310", Name: "There"}, {ID: "K423", Name: "Kian"}}

type CustomerServer struct {
	customerpb.UnimplementedCustomerServer
}

func (CustomerServer) GetCustomers(ctx context.Context, in *customerpb.GetCustomersRequest) (*customerpb.GetCustomersResponse, error) {
	res := &customerpb.GetCustomersResponse{}
	for _, customer := range _customers {
		res.Customers = append(res.Customers, &customerpb.CustomerInfo{
			CustomerId:   customer.ID,
			CustomerName: customer.Name,
		})
	}
	return res, nil
}
func (CustomerServer) GetCustomerById(ctx context.Context, in *customerpb.GetCustomerByIdRequest) (*customerpb.GetCustomerByIdResponse, error) {
	for _, customer := range _customers {
		if customer.ID == in.CustomerId {
			return &customerpb.GetCustomerByIdResponse{
				Customer: &customerpb.CustomerInfo{
					CustomerId:   customer.ID,
					CustomerName: customer.Name,
				},
			}, nil
		}
	}
	return nil, status.Errorf(codes.NotFound, "the customer does not exist")
}

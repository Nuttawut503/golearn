package main

import (
	"context"
	"grpc-api/client/server/customerpb"
	"log"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial(":9090", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	res, err := customerpb.NewCustomerClient(conn).GetCustomers(
		context.Background(),
		&customerpb.GetCustomersRequest{},
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(res.Customers)

	if len(res.Customers) > 0 {
		res, err := customerpb.NewCustomerClient(conn).GetCustomerById(
			context.Background(),
			&customerpb.GetCustomerByIdRequest{
				CustomerId: res.Customers[0].CustomerId,
			},
		)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(res.Customer.CustomerName)
	}
}

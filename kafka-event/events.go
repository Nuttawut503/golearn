package main

import "reflect"

var Topics = []string{
	reflect.TypeOf(RegisterAccountEvent{}).Name(),
	reflect.TypeOf(DeactivateAccountEvent{}).Name(),
}

type RegisterAccountEvent struct {
	TransactionID string `json:"transaction_id"`
	Email         string `json:"email"`
	Firstname     string `json:"firstname"`
	Lastname      string `json:"lastname"`
	Age           int    `json:"age"`
}

type DeactivateAccountEvent struct {
	TransactionID string `json:"transaction_id"`
	Email         string `json:"email"`
}

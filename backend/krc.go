package main

import (
	"fmt"

	"github.com/p2eengineering/kalp-sdk-public/kalpsdk"
)

type SmartContract struct {
	kalpsdk.Contract
}

func (s *SmartContract) Init(ctx kalpsdk.TransactionContextInterface) (bool, error) {
	s.Logger.Info("Greeting Smart Contract initialized")
	return true, nil
}

func (s *SmartContract) SetGreeting(ctx kalpsdk.TransactionContextInterface, greeting string) error {
	s.Logger.Info("Setting greeting")
	return ctx.PutStateWithoutKYC("greeting", []byte(greeting))
}

func (s *SmartContract) GetGreeting(ctx kalpsdk.TransactionContextInterface) (string, error) {
	s.Logger.Info("Getting greeting")
	greetingBytes, err := ctx.GetState("greeting")
	if err != nil {
		return "", fmt.Errorf("failed to read greeting: %v", err)
	}
	if greetingBytes == nil {
		return "", fmt.Errorf("greeting not found")
	}
	return string(greetingBytes), nil
}

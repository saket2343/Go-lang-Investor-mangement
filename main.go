package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

type Investor struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Balance float64 `json:"balance"`
}

func (s *SmartContract) RegisterInvestor(ctx contractapi.TransactionContextInterface, id string, name string, email string, balance float64) error {
	exists, err := s.InvestorExists(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to check if investor exists: %v", err)
	}
	if exists {
		return fmt.Errorf("investor with ID %s already exists", id)
	}

	investor := Investor{
		ID:      id,
		Name:    name,
		Email:   email,
		Balance: balance,
	}

	investorJSON, err := json.Marshal(investor)
	if err != nil {
		return fmt.Errorf("failed to marshal investor: %v", err)
	}

	return ctx.GetStub().PutState(id, investorJSON)
}

func (s *SmartContract) UpdateInvestor(ctx contractapi.TransactionContextInterface, id string, name string, email string, balance float64) error {
	exists, err := s.InvestorExists(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to check if investor exists: %v", err)
	}
	if !exists {
		return fmt.Errorf("investor with ID %s does not exist", id)
	}

	investor := Investor{
		ID:      id,
		Name:    name,
		Email:   email,
		Balance: balance,
	}

	investorJSON, err := json.Marshal(investor)
	if err != nil {
		return fmt.Errorf("failed to marshal investor: %v", err)
	}

	return ctx.GetStub().PutState(id, investorJSON)
}

func (s *SmartContract) ViewInvestor(ctx contractapi.TransactionContextInterface, id string) (*Investor, error) {
	investorJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if investorJSON == nil {
		return nil, fmt.Errorf("investor with ID %s does not exist", id)
	}

	var investor Investor
	err = json.Unmarshal(investorJSON, &investor)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal investor: %v", err)
	}

	return &investor, nil
}

func (s *SmartContract) InvestorExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	investorJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return investorJSON != nil, nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating chaincode: %v", err)
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting chaincode: %v", err)
	}
}

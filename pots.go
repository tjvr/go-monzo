package monzo

import (
	"strconv"
)

type Pot struct {
	ID       string `json:"id"`
	Deleted  bool   `json:"deleted"` // to be deprecated
	Name     string `json:"name"`
	Created  string `json:"created"`
	Updated  string `json:"updated"`
	Currency string `json:"currency"`
	Style    string `json:"style"`
	Balance  int64  `json:"balance"`
}

func (cl *Client) Pots() ([]*Pot, error) {
	rsp := &struct {
		Pots []*Pot `json:"pots"`
	}{}
	if err := cl.request("GET", "/pots", nil, rsp); err != nil {
		return nil, err
	}
	return rsp.Pots, nil
}

func (cl *Client) Pot(potID string) (*Pot, error) {
	pot := &Pot{}
	if err := cl.request("GET", "/pots/"+potID, nil, pot); err != nil {
		return nil, err
	}
	return pot, nil
}

type DepositRequest struct {
	PotID          string
	AccountID      string
	Amount         int64
	IdempotencyKey string
}

func (cl *Client) Deposit(op *DepositRequest) (*Pot, error) {
	args := map[string]string{
		"pot_id":            op.PotID,
		"source_account_id": op.AccountID,
		"amount":            strconv.FormatInt(op.Amount, 10),
		"dedupe_id":         op.IdempotencyKey,
	}
	pot := &Pot{}
	if err := cl.request("PUT", "/pots/"+op.PotID+"/deposit", args, pot); err != nil {
		return nil, err
	}
	return pot, nil
}

type WithdrawRequest struct {
	PotID          string
	AccountID      string
	Amount         int64
	IdempotencyKey string
}

func (cl *Client) Withdraw(op *WithdrawRequest) (*Pot, error) {
	args := map[string]string{
		"pot_id":                 op.PotID,
		"destination_account_id": op.AccountID,
		"amount":                 strconv.FormatInt(op.Amount, 10),
		"dedupe_id":              op.IdempotencyKey,
	}
	pot := &Pot{}
	if err := cl.request("PUT", "/pots/"+op.PotID+"/withdraw", args, pot); err != nil {
		return nil, err
	}
	return pot, nil
}

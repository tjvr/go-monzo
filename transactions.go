package monzo

import (
	"encoding/json"
	"fmt"
	"time"
)

/*
   "merchant": {
       "address": {
           "address": "98 Southgate Road",
           "city": "London",
           "country": "GB",
           "latitude": 51.54151,
           "longitude": -0.08482400000002599,
           "postcode": "N1 3JD",
           "region": "Greater London"
       },
       "created": "2015-08-22T12:20:18Z",
       "group_id": "grp_00008zIcpbBOaAr7TTP3sv",
       "id": "merch_00008zIcpbAKe8shBxXUtl",
       "logo": "https://pbs.twimg.com/profile_images/527043602623389696/68_SgUWJ.jpeg",
       "emoji": "🍞",
       "name": "The De Beauvoir Deli Co.",
       "category": "eating_out"
   }
*/

type TransactionEvent struct {
	Type        string       `json:"type"`
	Transaction *Transaction `json:"data"`
}

type Address struct {
	Address   string `json:"address"`
	City      string `json:"city"`
	Country   string `json:"country"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
	Postcode  string `json:"postcode"`
	Region    string `json:"region"`
}

type Merchant struct {
	ID              string            `json:"id"`
	Created         string            `json:"created"`
	Name            string            `json:"name"`
	Category        string            `json:"category"`
	IsOnline        bool              `json:"online"`
	Logo            string            `json:"logo"`
	GroupID         string            `json:"group_id"`
	DisableFeedback bool              `json:"disable_feedback"`
	Emoji           string            `json:"emoji"`
	IsATM           bool              `json:"atm"`
	Metadata        map[string]string `json:"metadata"`
	Address         *Address          `json:"address"`
}

type Transaction struct {
	ID                string            `json:"id"`
	Created           string            `json:"created"`
	AccountID         string            `json:"account_id"`
	Amount            int64             `json:"amount"`
	DeclineReason     string            `json:"decline_reason"`
	Scheme            string            `json:"scheme"`
	Currency          string            `json:"currency"`
	Description       string            `json:"description"`
	Category          string            `json:"category"`
	IsLoad            bool              `json:"is_load"` // is top up
	Settled           string            `json:"settled"`
	IncludeInSpending bool              `json:"include_in_spending"`
	LocalAmount       int64             `json:"local_amount"`
	LocalCurrency     string            `json:"local_currency"`
	IsOriginator      bool              `json:"originator"`
	DedupeID          string            `json:"dedupe_id"`
	Metadata          map[string]string `json:"metadata"`
	Notes             string            `json:"notes"`
	AccountBalance    int64             `json:"account_balance"` // not for CA
	Merchant          *Merchant
}

type RawTransaction struct {
	Transaction
	Merchant json.RawMessage `json:"merchant"`
}

type TransactionsInput struct {
	AccountId string
	Since     time.Time
	Before    time.Time
}

func (cl *Client) Transactions(input TransactionsInput, expandMerchant bool) ([]*Transaction, error) {
	args := map[string]string{
		"account_id": input.AccountId,
	}
	if !input.Since.IsZero() {
		args["since"] = input.Since.Format(time.RFC3339)
	}
	if !input.Before.IsZero() {
		args["before"] = input.Since.Format(time.RFC3339)
	}
	rsp := &struct {
		Transactions []*RawTransaction `json:"transactions"`
	}{}
	if err := cl.request("GET", "/transactions", args, rsp); err != nil {
		return nil, err
	}
	transactions, err := unmarshalTransactionList(rsp.Transactions, expandMerchant)
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func (cl *Client) Transaction(id string) (*Transaction, error) {
	args := map[string]string{
		"expand[]": "merchant",
	}
	rsp := &struct {
		Transaction *RawTransaction `json:"transaction"`
	}{}
	if err := cl.request("GET", "/transactions/"+id, args, rsp); err != nil {
		return nil, err
	}
	return unmarshalRawTransaction(rsp.Transaction, true)
}

func (cl *Client) AnnotateTransaction(id string, metadata map[string]string) (*Transaction, error) {
	args := map[string]string{}
	for k, v := range metadata {
		args[fmt.Sprintf("metadata[%s]", k)] = v
	}
	tx := &Transaction{}
	if err := cl.request("PATCH", "/transactions/"+id, args, tx); err != nil {
		return nil, err
	}
	return tx, nil
}

func unmarshalRawTransaction(raw *RawTransaction, expandMerchant bool) (*Transaction, error) {
	tx := &raw.Transaction
	tx.Merchant = &Merchant{}
	var target interface{}
	if expandMerchant {
		target = tx.Merchant
	} else {
		target = &tx.Merchant.ID
	}
	if err := json.Unmarshal(raw.Merchant, target); err != nil {
		return nil, err
	}
	return tx, nil
}

func unmarshalTransactionList(raw []*RawTransaction, expandMerchant bool) ([]*Transaction, error) {
	transactions := make([]*Transaction, len(raw))
	for i, raw := range raw {
		tx, err := unmarshalRawTransaction(raw, expandMerchant)
		if err != nil {
			return nil, err
		}
		transactions[i] = tx
	}
	return transactions, nil
}

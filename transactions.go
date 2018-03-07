package monzo

import "fmt"

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
	Metadata        map[string]string `json:"metadata"`
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
	IsATM             bool              `json:"atm"`
	Settled           string            `json:"settled"`
	IncludeInSpending bool              `json:"include_in_spending"`
	LocalAmount       int64             `json:"local_amount"`
	LocalCurrency     string            `json:"local_currency"`
	IsOriginator      bool              `json:"originator"`
	DedupeID          string            `json:"dedupe_id"`
	Metadata          map[string]string `json:"metadata"`
	Notes             string            `json:"notes"`
	AccountBalance    int64             `json:"account_balance"` // not for CA
	//Labels
	//Counterparty
	//Fees
	//Attachments

	Merchant   *Merchant `json:"merchant"`
	MerchantID string    `json:"merchant"`
}

func (cl *Client) Transactions(accountID string, expandMerchant bool) ([]*Transaction, error) {
	args := map[string]string{
		"account_id": accountID,
	}
	if expandMerchant {
		args["expand[]"] = "merchant"
	}
	rsp := &struct {
		Transactions []*Transaction `json:"transactions"`
	}{}
	if err := cl.request("GET", "/transactions", args, rsp); err != nil {
		return nil, err
	}
	return rsp.Transactions, nil
}

func (cl *Client) Transaction(id string) (*Transaction, error) {
	args := map[string]string{
		"expand[]": "merchant",
	}
	tx := &Transaction{}
	if err := cl.request("GET", "/transactions/"+id, args, tx); err != nil {
		return nil, err
	}
	return tx, nil
}

// TODO I'm not convinced this endpoint works.
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

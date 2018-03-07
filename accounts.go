package monzo

type Account struct {
	ID            string `json:"id"`
	Closed        bool   `json:"closed"`
	Created       string `json:"created"`
	Description   string `json:"description"`
	Type          string `json:"type"`
	SortCode      string `json:"sort_code"`
	AccountNumber string `json:"account_number"`
}

func (cl *Client) Accounts(accountType string) ([]*Account, error) {
	args := map[string]string{}
	if accountType != "" {
		args["account_type"] = accountType
	}
	rsp := &struct {
		Accounts []*Account `json:"accounts"`
	}{}
	if err := cl.request("GET", "/accounts", args, rsp); err != nil {
		return nil, err
	}
	return rsp.Accounts, nil
}

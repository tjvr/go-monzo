package monzo

type Balance struct {
	Balance    int64  `json:"balance"`
	Currency   string `json:"currency"`
	SpendToday int64  `json:"spend_today"`
}

func (cl *Client) Balance(accountID string) (*Balance, error) {
	rsp := &Balance{}
	if err := cl.request("GET", "/balance", map[string]string{
		"account_id": accountID,
	}, rsp); err != nil {
		return nil, err
	}
	return rsp, nil
}

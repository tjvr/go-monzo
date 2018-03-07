package monzo

type Webhook struct {
	ID        string `json:"id"`
	AccountID string `json:"account_id"`
	URL       string `json:"url"`
}

func (cl *Client) RegisterWebhook(hook *Webhook) (*Webhook, error) {
	args := map[string]string{
		"account_id": hook.AccountID,
		"url":        hook.URL,
	}
	body := &Webhook{}
	if err := cl.request("POST", "/webhooks", args, body); err != nil {
		return nil, err
	}
	return body, nil
}

func (cl *Client) Webhooks(accountID string) ([]*Webhook, error) {
	args := map[string]string{
		"account_id": accountID,
	}
	rsp := &struct {
		Webhooks []*Webhook `json:"webhooks"`
	}{}
	if err := cl.request("GET", "/webhooks", args, rsp); err != nil {
		return nil, err
	}
	return rsp.Webhooks, nil
}

func (cl *Client) DeleteWebhook(webhookID string) error {
	return cl.request("DELETE", "/webhooks/"+webhookID, nil, nil)
}

type Event struct {
	Type string `json:"type"`
}

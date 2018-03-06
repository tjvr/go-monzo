package monzo

type FeedItem struct {
	AccountID string
	Type      string
	URL       string
	Title     string
	Body      string
	ImageURL  string
}

func (cl *Client) CreateFeedItem(item *FeedItem) error {
	args := map[string]string{
		"account_id":        item.AccountID,
		"type":              item.Type,
		"url":               item.URL,
		"params[title]":     item.Title,
		"params[body]":      item.Body,
		"params[image_url]": item.ImageURL,
	}
	rsp := map[string]string{}
	return cl.request("POST", "/feed", args, rsp)
}

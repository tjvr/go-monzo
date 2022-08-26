package monzo

type FeedItem struct {
	AccountID       string
	Type            string
	URL             string
	Title           string
	Body            string
	ImageURL        string
	BackgroundColor string
	TitleColor      string
}

func (cl *Client) CreateFeedItem(item *FeedItem) error {
	args := map[string]string{
		"account_id":               item.AccountID,
		"type":                     item.Type,
		"url":                      item.URL,
		"params[title]":            item.Title,
		"params[body]":             item.Body,
		"params[image_url]":        item.ImageURL,
		"params[background_color]": item.BackgroundColor,
		"params[title_color]":      item.BackgroundColor,
	}
	rsp := &map[string]string{}
	return cl.request("POST", "/feed", args, rsp)
}

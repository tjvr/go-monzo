package monzo

type AccessToken struct {
	AccessToken  string `json:"access_token"`
	ClientID     string `json:"client_id"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	UserID       string `json:"user_id"`
}

func (cl *Client) tokenRequest(args map[string]string) (*AccessToken, error) {
	rsp := &AccessToken{}
	if err := cl.request("POST", "/oauth2/token", args, rsp); err != nil {
		return nil, err
	}
	return rsp, nil
}

func (auth *Authenticator) Authorize(authorizationCode string) (*Client, error) {
	cl := &Client{
		BaseURL: auth.BaseURL,
	}
	rsp, err := cl.tokenRequest(map[string]string{
		"grant_type":    "authorization_code",
		"client_id":     auth.ClientID,
		"client_secret": auth.ClientSecret,
		"redirect_uri":  auth.CallbackURI,
		"code":          authorizationCode,
	})
	if err != nil {
		return nil, err
	}
	cl.Init(rsp)
	return cl, nil
}

func (cl *Client) Init(rsp *AccessToken) {
	cl.AccessToken = rsp.AccessToken
	cl.RefreshToken = rsp.RefreshToken
	cl.UserID = rsp.UserID
}

func (auth *Authenticator) Refresh(refreshToken string) error {
	cl := &Client{
		BaseURL: auth.BaseURL,
	}
	rsp, err := cl.tokenRequest(map[string]string{
		"grant_type":    "refresh_token",
		"client_id":     auth.ClientID,
		"client_secret": auth.ClientSecret,
		"refresh_token": refreshToken,
	})
	if err != nil {
		return err
	}
	cl.Init(rsp)
	return nil
}

func (auth *Authenticator) RefreshClient(cl *Client) error {
	rsp, err := cl.tokenRequest(map[string]string{
		"grant_type":    "refresh_token",
		"client_id":     auth.ClientID,
		"client_secret": auth.ClientSecret,
		"refresh_token": cl.RefreshToken,
	})
	if err != nil {
		return err
	}
	cl.Init(rsp)
	return nil
}

// TODO whoami

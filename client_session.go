package monzo

import (
	"fmt"
	"net/http"
	"net/url"
)

var CookieName = "monzoSession"

type Session struct {
	Cookie string
	State  string
	Client *Client
}

func (s *Session) SetCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:  CookieName,
		Value: s.Cookie,
		// TODO more options?
	})
}

func (s *Session) IsAuthenticated() bool {
	return s.Client != nil
}

type Authenticator struct {
	BaseURL      string
	ClientID     string
	ClientSecret string
	CallbackURI  string

	Sessions map[string]*Session
}

func NewAuthenticator(clientID, clientSecret, callbackURI string) *Authenticator {
	switch {
	case clientID == "":
		panic("no clientID provided")
	case clientSecret == "":
		panic("no clientSecret provided")
	case callbackURI == "":
		panic("no callbackURI provided")
	}
	return &Authenticator{
		BaseURL:      "https://api.monzo.com",
		ClientID:     clientID,
		ClientSecret: clientSecret,
		CallbackURI:  callbackURI,
		Sessions:     map[string]*Session{},
	}
}

func (auth *Authenticator) GetSession(w http.ResponseWriter, req *http.Request) *Session {
	sessionCookie := getSessionCookie(req)
	s := auth.Sessions[sessionCookie]
	if s != nil {
		s.SetCookie(w)
		return s
	}
	s = &Session{
		Cookie: randomString(),
		State:  randomString(),
	}
	s.SetCookie(w)
	auth.Sessions[s.Cookie] = s
	return s
}

func (auth *Authenticator) BeginAuthURL(s *Session) string {
	u := url.URL{Scheme: "https", Host: "auth.monzo.com"}
	q := u.Query()
	q.Set("client_id", auth.ClientID)
	q.Set("redirect_uri", auth.CallbackURI)
	q.Set("response_type", "code")
	q.Set("state", s.State)
	u.RawQuery = q.Encode()
	return u.String()
}

func (auth *Authenticator) Login(w http.ResponseWriter, req *http.Request) {
	s := auth.GetSession(w, req)
	redirectURL := auth.BeginAuthURL(s)
	http.Redirect(w, req, redirectURL, http.StatusTemporaryRedirect)
}

func (auth *Authenticator) Logout(w http.ResponseWriter, req *http.Request) {
	delete(auth.Sessions, getSessionCookie(req))
}

func (auth *Authenticator) Callback(w http.ResponseWriter, req *http.Request) *Session {
	s := auth.GetSession(w, req)
	query := req.URL.Query()
	if query.Get("state") != s.State {
		// TODO 401
		fmt.Fprintf(w, "invalid state")
		return nil
	}

	authorizationCode := query.Get("code")
	if authorizationCode == "" {
		// TODO 400
		fmt.Fprintf(w, "missing code")
		return nil
	}

	cl, err := auth.Authorize(authorizationCode)
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err)
		// TODO 400 or maybe 500?
		return nil
	}
	s.Client = cl
	s.SetCookie(w)
	return s
}

func (auth *Authenticator) EnsureAuthenticated(w http.ResponseWriter, req *http.Request) *Client {
	s := auth.GetSession(w, req)
	if s.IsAuthenticated() {
		return s.Client
	}

	// Redirect to login page
	auth.Login(w, req)
	return nil
}

func randomString() string {
	return "foo"
}

func getSessionCookie(req *http.Request) string {
	for _, cookie := range req.Cookies() {
		if cookie.Name == "monzoSession" {
			return cookie.Value
		}
	}
	return ""
}

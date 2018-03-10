go-monzo
========

A Go client for the public Monzo API.


Usage
-----

Import the `monzo` package.

```go
import (
    "github.com/tjvr/go-monzo"
)
```

Make a client from an existing `access_token` (e.g. from the [Developer Playground](https://developers.monzo.com/api/playground)).

```go
cl := monzo.Client{
    BaseURL: "https://api.monzo.com",
    AccessToken: os.Getenv("access_token"),
}
```

Get an account.

```go
accounts, err := cl.Accounts("uk_retail")
if err != nil {
    return err
}
if len(accounts) == 0 {
    return fmt.Errorf("no retail account")
}
acc := accounts[0]
```

Get the balance of an account.

```go
rsp, err := cl.Balance(acc.ID)
if err != nil {
    return err
}
mainBalance := rsp.Balance
```

List transactions.

```go
transactions, err := cl.Transactions(acc.ID, false) // don't expandMerchant
if err != nil {
    return err
}
firstTx := transactions[0]
```

Get transaction details.

```go
tx, err := cl.Transaction(firstTx.ID)
if err != nil {
    return err
}
merchant := tx.Merchant
```

List pots.

```go
pots, err := cl.Pots()
if err != nil {
    return err
}
firstPot := pots[0]
```

Get a single pot.

```go
pot, err := cl.Pot(firstPot.ID)
if err != nil {
    return err
}
potBalance := pot.Balance
```

Withdraw (or deposit!) from a pot.

```go
updatedPot, err := cl.Deposit(&monzo.DepositRequest{
	PotID: pot.ID,
	AccountID: acc.ID,
	Amount: potBalance,
	IdempotencyKey: idempotencyKey,
}
```

Post to the user's feed.

```go
err := CreateFeedItem(&monzo.FeedItem{
	AccountID: acc.ID,
	Type: "basic",
	URL: "https://www.example.com/a_page_to_open_on_tap.html",
	Title: "My custom item",
	Body: "Some body text to display",
	ImageURL: "www.example.com/image.png",
})
if err != nil {
    return err
}
```

OAuth
-----

`go-monzo` tries to provide helpful HTTP handlers for you, so you don't have to think about OAuth very much. This is probably a mistake, but there we go!

Create an object to handle authentication and sessions.

```go
auth := monzo.NewAuthenticator(
    os.Getenv("CLIENT_ID"),
    os.Getenv("CLIENT_SECRET"),
    os.Getenv("CALLBACK_URI"),
)
```

Get the `go-monzo` Session object from an HTTP request.

```go
func myHandler(w http.ResponseWriter, req *http.Request) {
    session := auth.
```

Use the built-in `auth.Login` handler to redirect an HTTP request (e.g. all requests to `/login`) to Monzo's authentication page.

```go
http.HandleFunc("/login", auth.Login)
```

Add a handler for the OAuth callback. This needs to exactly match both the `REDIRECT_URI` you passed to `NewAuthenticator`, and the Redirect URI set for your OAuth Client on the Monzo developer site.

```go
http.HandleFunc("/callback", func(w http.ResponseWriter, req *http.Request) {
    session := auth.Callback(w, req)
    if session == nil {
        // TODO something went wrong
    }
    http.Redirect(w, req, "/", http.StatusTemporaryRedirect)
})
```

Add a logout handler.

```go
http.HandleFunc("/logout", func(w http.ResponseWriter, req *http.Request) {
    auth.Logout(w, req)
    http.Redirect(w, req, "/", http.StatusTemporaryRedirect)
})
```

Check whether a request is authenticated.

```go
session := auth.GetSession(w, req)
if !session.IsAuthenticated() {
    http.Error(w, "Not Authenticated", 403)
    return
}
cl := session.Client
```




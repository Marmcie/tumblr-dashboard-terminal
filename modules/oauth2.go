package modules

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/user"
	"time"

	"github.com/zalando/go-keyring"
	"golang.org/x/oauth2"
)

func getAuthConfig() *oauth2.Config {
	config := GetConfig()
	return &oauth2.Config{
		ClientID:     config.Consumer_key,
		ClientSecret: config.Secret_key,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://www.tumblr.com/oauth2/authorize",
			TokenURL: "https://api.tumblr.com/v2/oauth2/token",
		},
		RedirectURL: "http://localhost:" + config.Redirect_port + "/",
		Scopes: []string{
			"basic",
			"offline_access",
		},
	}

}

type OAuth2Token struct {
	Token      *oauth2.Token
	created_at int64
}

var Expires_at int64

func GetClient() *http.Client {
	conf := getAuthConfig()
	ctx := context.Background()

	usr, _ := user.Current()
	user := usr.Username
	service := "tumblr-terminal-token"

	// get password
	tokenStr, _ := keyring.Get(service, user)

	var token *OAuth2Token
	if len(tokenStr) == 0 {
		//INFO:Create new token
		token = &OAuth2Token{}
		token.Token = Auth(ctx)
		token.created_at = time.Now().Unix()
		Expires_at = token.created_at + token.Token.ExpiresIn
		bytes, _ := json.Marshal(token)
		keyring.Set(service, user, string(bytes))
	} else {
		//INFO:Get token from keyring
		token = &OAuth2Token{}
		json.Unmarshal([]byte(tokenStr), token)

		Expires_at = token.created_at + token.Token.ExpiresIn

		//INFO:Attempt token refresh
		if TokenExpired() {
			token.Token = Refresh(ctx, token.Token.RefreshToken)
			token.created_at = time.Now().Unix()
			Expires_at = token.created_at + token.Token.ExpiresIn
			bytes, _ := json.Marshal(token)
			keyring.Set(service, user, string(bytes))
		}
	}

	return conf.Client(ctx, token.Token)

}

func TokenExpired() bool {
	now := time.Now().Unix()
	//INFO:Attempt token refresh
	return now >= Expires_at
}

func RemoveToken() {
	usr, _ := user.Current()
	user := usr.Username
	service := "tumblr-terminal-token"
	keyring.Delete(service, user)
}

func Auth(ctx context.Context) *oauth2.Token {

	conf := getAuthConfig()
	config := GetConfig()

	verifier := oauth2.GenerateVerifier()

	state := rand.Text()

	requestUrl := conf.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.S256ChallengeOption(verifier))

	fmt.Printf("Visit the URL for the auth dialog: %v", requestUrl)
	srv := &http.Server{Addr: ":" + config.Redirect_port}
	var token *oauth2.Token
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		returnedState := r.URL.Query().Get("state")
		if state != returnedState {
			panic("Incorrect state was returned to redirect URL.")
		}
		code := r.URL.Query().Get("code")
		if len(code) > 0 {
			tok, err := conf.Exchange(ctx, code, oauth2.VerifierOption(verifier))
			if err != nil {
				print("OAuth2 key exchange failed.\n")
				log.Fatal(err)
			}
			token = tok
			srv.Shutdown(ctx)
			return
		}
	})

	srv.ListenAndServe()
	return token
}

func Refresh(ctx context.Context, refreshToken string) *oauth2.Token {
	conf := getAuthConfig()
	url := "http://api.tumblr.com/v2/oauth2/token"

	b := map[string]string{
		"grant_type":    "refresh_token",
		"client_id":     conf.ClientID,
		"client_secret": conf.ClientSecret,
		"refresh_token": refreshToken,
	}

	str, _ := json.Marshal(b)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(str))
	if err != nil {
		print("Request creation error")
		RemoveToken()
		panic(err)
	}
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		print("Request error")
		RemoveToken()
		panic(err)
	}

	response, _ := io.ReadAll(resp.Body)
	var tok *oauth2.Token

	json.Unmarshal(response, &tok)

	return tok

}

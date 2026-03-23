package modules

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/user"

	"github.com/zalando/go-keyring"
	"golang.org/x/oauth2"
)

func GetClient() *http.Client {
	config := GetConfig()
	conf := &oauth2.Config{
		ClientID:     config.Consumer_key,
		ClientSecret: config.Secret_key,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://www.tumblr.com/oauth2/authorize",
			TokenURL: "https://api.tumblr.com/v2/oauth2/token",
		},
		RedirectURL: "http://localhost:6969/",
	}

	usr, _ := user.Current()
	user := usr.Username
	service := "tumblr-terminal-token"

	// get password
	tokenStr, _ := keyring.Get(service, user)

	var token *oauth2.Token
	if len(tokenStr) == 0 {
		token = Auth()
		bytes, _ := json.Marshal(token)
		keyring.Set(service, user, string(bytes))
	} else {
		token = &oauth2.Token{}
		json.Unmarshal([]byte(tokenStr), token)
	}

	ctx := context.Background()

	return conf.Client(ctx, token)

}

func RemoveToken() {
	usr, _ := user.Current()
	user := usr.Username
	service := "tumblr-terminal-token"
	err := keyring.Delete(service, user)

	print("Invalid OAuth2 token\n")
	log.Fatal(err)
}

func Auth() *oauth2.Token {
	ctx := context.Background()
	config := GetConfig()

	conf := &oauth2.Config{
		ClientID:     config.Consumer_key,
		ClientSecret: config.Secret_key,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://www.tumblr.com/oauth2/authorize",
			TokenURL: "https://api.tumblr.com/v2/oauth2/token",
		},
		RedirectURL: "http://localhost:6969/",
	}
	verifier := oauth2.GenerateVerifier()

	requestUrl := conf.AuthCodeURL("state", oauth2.AccessTypeOffline, oauth2.S256ChallengeOption(verifier))

	fmt.Printf("Visit the URL for the auth dialog: %v", requestUrl)
	srv := &http.Server{Addr: ":6969"}
	var token *oauth2.Token
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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

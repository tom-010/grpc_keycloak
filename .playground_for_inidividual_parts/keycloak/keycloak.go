package main

import (
	"context"
	"crypto/tls"
	"log"

	gocloak "github.com/Nerzal/gocloak/v11"
)

const (
	realm        = "skytala"
	clientID     = "clientid-03"
	clientSecret = "IslMJ7tUpEEWAjN6osCD5RLqmIZVsMnM"
	username     = "t.deniffel"
	password     = "password"
	URL          = "http://localhost:8080"
)

func main() {
	client := gocloak.NewClient(URL, gocloak.SetAuthAdminRealms("admin/realms"), gocloak.SetAuthRealms("realms"))
	restyClient := client.RestyClient()
	restyClient.SetDebug(false)
	restyClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	ctx := context.Background()
	token, err := client.Login(ctx, clientID, clientSecret, realm, username, password)
	if err != nil {
		log.Fatalf("Login failed: %v", err)
	}

	log.Println(token.AccessToken)

	log.Println(client.GetUserInfo(ctx, token.AccessToken, realm))

	rptResult, err := client.RetrospectToken(ctx, token.AccessToken, clientID, clientSecret, realm)
	if err != nil {
		log.Fatalf("Inspection failed: %v", err)
	}

	if !*rptResult.Active {
		log.Fatal("Token is not active")
	}

	// permissions := rptResult.Permissions
	log.Println(rptResult)
}

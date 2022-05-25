package main

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Nerzal/gocloak/v11"
)

const (
	adminUser       = "admin"
	adminPassword   = "admin"
	realm           = "skytala"
	masterRealm     = "master"
	clientID        = "clientid-03"
	clientSecret    = "kpB0LPPyMeRlhAq8fk2bkWOqlxz3d7Ah"
	keycloakUrl     = "http://localhost:8080"
	initialPassword = "password"
)

func main() {
	username := "tom" + fmt.Sprint(time.Now().Unix())
	createNewUser("Thomas", "Deniffel", username, username+"@localhost")
	token := login(username, initialPassword)

	// output the result (by base64-decoding the second part)
	bytes, _ := base64.RawStdEncoding.DecodeString(strings.Split(token.AccessToken, ".")[1])
	log.Println(string(bytes))
}

func createNewUser(fistName string, lastName string, username string, email string) {

	client := gocloak.NewClient(keycloakUrl, gocloak.SetAuthAdminRealms("admin/realms"), gocloak.SetAuthRealms("realms"))
	restyClient := client.RestyClient()
	restyClient.SetDebug(false)
	restyClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

	ctx := context.Background()
	adminToken, err := client.LoginAdmin(ctx, adminUser, adminPassword, masterRealm)
	if err != nil {
		log.Fatalf("something wrong with the credentials or url: %v", err)
	}

	user := gocloak.User{
		FirstName: gocloak.StringP(fistName),
		LastName:  gocloak.StringP(lastName),
		Email:     gocloak.StringP(email),
		Enabled:   gocloak.BoolP(true),
		Username:  gocloak.StringP(username),
	}

	log.Println(adminToken.AccessToken)

	newUserId, err := client.CreateUser(ctx, adminToken.AccessToken, realm, user)
	if err != nil {
		log.Fatalf("failed to create user: %v", err)
	}

	err = client.SetPassword(ctx, adminToken.AccessToken, newUserId, realm, initialPassword, false)
	if err != nil {
		log.Printf("Could not set password: %v", err)
	}
}

func login(username string, password string) *gocloak.JWT {
	ctx := context.Background()
	client := gocloak.NewClient(keycloakUrl, gocloak.SetAuthAdminRealms("admin/realms"), gocloak.SetAuthRealms("realms"))
	restyClient := client.RestyClient()
	restyClient.SetDebug(false)
	restyClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

	jwt, err := client.Login(ctx, clientID, clientSecret, realm, username, password)
	if err != nil {
		log.Fatalf("Could not login: %v", err)
	}

	return jwt
}

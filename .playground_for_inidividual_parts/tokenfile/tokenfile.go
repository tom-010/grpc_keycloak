package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/Nerzal/gocloak"
	"github.com/dgrijalva/jwt-go"
)

const (
	realm        = "skytala"
	clientID     = "clientid-03"
	clientSecret = "bI6wznsqBH3dN1UWeRk8Yz4Xepp1li7D"
	username     = "t.deniffel"
	password     = "password"
	URL          = "http://localhost:8080"
)

// "email_verified": true,
//   "name": "Thomas Deniffel",
//   "preferred_username": "t.deniffel",
//   "given_name": "Thomas",
//   "family_name": "Deniffel",
//   "email": "t.deniffel@skytala-gmbh.com"

type user struct {
	Username      string
	Name          string
	GivenName     string
	FamilyName    string
	Email         string
	EmailVerified bool
}

func toUser(claims jwt.MapClaims) user {
	res := user{}
	m := map[string]string{}
	for key, val := range claims {
		fmt.Printf("Key: %v, value: %v\n", key, val)
	}
	res.Username = m["preferred_username"]
	return res
}

func main() {
	tokenString, err := ioutil.ReadFile("token.json")
	if err != nil {
		log.Fatalf("could not open token file: %v", err)
	}

	keycloakToken := new(gocloak.JWT)
	json.Unmarshal(tokenString, keycloakToken)

	claims := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(keycloakToken.AccessToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(clientID), nil
	})
	if err != nil {
		log.Fatalf("parse claims err: %v", err)
	}

	log.Println(toUser(claims))
}

// token, err := jwt.Parse(keycloakToken.AccessToken, func(token *jwt.Token) (interface{}, error) {
// 	// Don't forget to validate the alg is what you expect:
// 	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 		return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
// 	}

// 	// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
// 	hmacSampleSecret := []byte("my_secret_key")
// 	return hmacSampleSecret, nil
// })

// // if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
// // 	fmt.Println(claims["foo"], claims["nbf"])
// // } else {
// // 	fmt.Println(err)
// // }

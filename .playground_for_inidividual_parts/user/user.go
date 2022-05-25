package main

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"strings"
)

type User struct {
	Id        string
	Name      string
	Username  string
	FirstName string
	LastName  string
	Email     string
	Jwt       JwtBody
}

type JwtBody struct {
	Exp            int      `json:"exp,omitempty"`
	Iat            int      `json:"iat,omitempty"`
	Jti            string   `json:"jti,omitempty"`
	Iss            string   `json:"iss,omitempty"`
	Aud            string   `json:"aud,omitempty"`
	Sub            string   `json:"sub,omitempty"`
	Typ            string   `json:"typ,omitempty"`
	Azp            string   `json:"azp,omitempty"`
	Sessionstate   string   `json:"session_state,omitempty"`
	Acr            string   `json:"acr,omitempty"`
	AllowedOrigins []string `json:"allowed-origins,omitempty"`
	RealmAccess    struct {
		Roles []string `json:"roles,omitempty"`
	} `json:"realm_access,omitempty"`
	ResourceAccess struct {
		Account struct {
			Roles []string `json:"roles,omitempty"`
		} `json:"account,omitempty"`
	} `json:"resource_access,omitempty"`
	Scope             string `json:"scope,omitempty"`
	Sid               string `json:"sid,omitempty"`
	EmailVerified     bool   `json:"email_verified,omitempty"`
	Name              string `json:"name,omitempty"`
	PreferredUsername string `json:"preferred_username,omitempty"`
	GivenName         string `json:"given_name,omitempty"`
	FamilyName        string `json:"family_name,omitempty"`
	Email             string `json:"email,omitempty"`
}

func main() {
	tokenString := "eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJqMmwxTUJLM0VOazg1SDBLbzVsY1lKN09RMk9pd1NVNHRVOS1RelNvLVZrIn0.eyJleHAiOjE2NTMzMDMyMTcsImlhdCI6MTY1MzMwMjkxNywianRpIjoiMTAzNDIwNWEtZTg2OS00YjlhLWFkNmUtOTBhMWRmNjkyZTU1IiwiaXNzIjoiaHR0cDovL2xvY2FsaG9zdDo4MDgwL3JlYWxtcy9za3l0YWxhIiwiYXVkIjoiYWNjb3VudCIsInN1YiI6IjkzOTNhMTNmLWMzZTItNGI3MS05YTUxLWJhN2EyOWY3NDdlNSIsInR5cCI6IkJlYXJlciIsImF6cCI6ImNsaWVudGlkLTAzIiwic2Vzc2lvbl9zdGF0ZSI6IjY1NjIwYWY0LWEzMTktNDExYy1hNDAyLTc1MTUzMzE3ZThhYiIsImFjciI6IjEiLCJhbGxvd2VkLW9yaWdpbnMiOlsiaHR0cDovL2xvY2FsaG9zdDo4MDgwIl0sInJlYWxtX2FjY2VzcyI6eyJyb2xlcyI6WyJkZWZhdWx0LXJvbGVzLXNreXRhbGEiLCJvZmZsaW5lX2FjY2VzcyIsInVtYV9hdXRob3JpemF0aW9uIl19LCJyZXNvdXJjZV9hY2Nlc3MiOnsiYWNjb3VudCI6eyJyb2xlcyI6WyJtYW5hZ2UtYWNjb3VudCIsIm1hbmFnZS1hY2NvdW50LWxpbmtzIiwidmlldy1wcm9maWxlIl19fSwic2NvcGUiOiJvcGVuaWQgZW1haWwgcHJvZmlsZSIsInNpZCI6IjY1NjIwYWY0LWEzMTktNDExYy1hNDAyLTc1MTUzMzE3ZThhYiIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJuYW1lIjoiVGhvbWFzIERlbmlmZmVsIiwicHJlZmVycmVkX3VzZXJuYW1lIjoidC5kZW5pZmZlbCIsImdpdmVuX25hbWUiOiJUaG9tYXMiLCJmYW1pbHlfbmFtZSI6IkRlbmlmZmVsIiwiZW1haWwiOiJ0LmRlbmlmZmVsQHNreXRhbGEtZ21iaC5jb20ifQ.neqvbOBFsL81jxGD-ivC-KEeJUMOJELwsKi1j5W63quS4-z5ZvaOxQI4d-CllKmHpNJalnyo7C8QKUFsj7FSISihlWV-llzSZm1st09j0HndVTBk5FMYBgK9TGl0wSxPN9CHCNXGKMHcvezApnnAEb3l5EJkz3199ZPh8-lszUrOc6WIxs0sa_FBPeTQodrUWXwKLLjjpz5GCpgredX7EuWuqOtZN3pZZbK6Ob1KQbwEU0apet0TTZD-MQm3l88yTOznGo7bgWLq7eliwkaSlFb1CfPFEbv2rum1iTBNcQg51__6HEW9LbPS1LbVneLFFjf2UZ5UfyOuyMQ5qYI7xw"
	toUser(tokenString)
}

func toUser(tokenString string) {
	parts := strings.Split(tokenString, ".")
	userBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		log.Fatalf("failed to base64-decode: %v", err)
	}
	user := User{}
	err = json.Unmarshal(userBytes, &user.Jwt)
	if err != nil {
		log.Fatalf("Failed to parse user: %v", err)
	}

	user.Id = user.Jwt.Sub
	user.Username = user.Jwt.PreferredUsername
	user.Name = user.Jwt.Name
	user.FirstName = user.Jwt.GivenName
	user.LastName = user.Jwt.FamilyName
	user.Email = user.Jwt.Email

	log.Printf("%v", user.FirstName)
}

package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	certsUrl := "http://localhost:8080/realms/skytala/protocol/openid-connect/certs"
	tokenString := "eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJqMmwxTUJLM0VOazg1SDBLbzVsY1lKN09RMk9pd1NVNHRVOS1RelNvLVZrIn0.eyJleHAiOjE2NTMzMDMyMTcsImlhdCI6MTY1MzMwMjkxNywianRpIjoiMTAzNDIwNWEtZTg2OS00YjlhLWFkNmUtOTBhMWRmNjkyZTU1IiwiaXNzIjoiaHR0cDovL2xvY2FsaG9zdDo4MDgwL3JlYWxtcy9za3l0YWxhIiwiYXVkIjoiYWNjb3VudCIsInN1YiI6IjkzOTNhMTNmLWMzZTItNGI3MS05YTUxLWJhN2EyOWY3NDdlNSIsInR5cCI6IkJlYXJlciIsImF6cCI6ImNsaWVudGlkLTAzIiwic2Vzc2lvbl9zdGF0ZSI6IjY1NjIwYWY0LWEzMTktNDExYy1hNDAyLTc1MTUzMzE3ZThhYiIsImFjciI6IjEiLCJhbGxvd2VkLW9yaWdpbnMiOlsiaHR0cDovL2xvY2FsaG9zdDo4MDgwIl0sInJlYWxtX2FjY2VzcyI6eyJyb2xlcyI6WyJkZWZhdWx0LXJvbGVzLXNreXRhbGEiLCJvZmZsaW5lX2FjY2VzcyIsInVtYV9hdXRob3JpemF0aW9uIl19LCJyZXNvdXJjZV9hY2Nlc3MiOnsiYWNjb3VudCI6eyJyb2xlcyI6WyJtYW5hZ2UtYWNjb3VudCIsIm1hbmFnZS1hY2NvdW50LWxpbmtzIiwidmlldy1wcm9maWxlIl19fSwic2NvcGUiOiJvcGVuaWQgZW1haWwgcHJvZmlsZSIsInNpZCI6IjY1NjIwYWY0LWEzMTktNDExYy1hNDAyLTc1MTUzMzE3ZThhYiIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJuYW1lIjoiVGhvbWFzIERlbmlmZmVsIiwicHJlZmVycmVkX3VzZXJuYW1lIjoidC5kZW5pZmZlbCIsImdpdmVuX25hbWUiOiJUaG9tYXMiLCJmYW1pbHlfbmFtZSI6IkRlbmlmZmVsIiwiZW1haWwiOiJ0LmRlbmlmZmVsQHNreXRhbGEtZ21iaC5jb20ifQ.neqvbOBFsL81jxGD-ivC-KEeJUMOJELwsKi1j5W63quS4-z5ZvaOxQI4d-CllKmHpNJalnyo7C8QKUFsj7FSISihlWV-llzSZm1st09j0HndVTBk5FMYBgK9TGl0wSxPN9CHCNXGKMHcvezApnnAEb3l5EJkz3199ZPh8-lszUrOc6WIxs0sa_FBPeTQodrUWXwKLLjjpz5GCpgredX7EuWuqOtZN3pZZbK6Ob1KQbwEU0apet0TTZD-MQm3l88yTOznGo7bgWLq7eliwkaSlFb1CfPFEbv2rum1iTBNcQg51__6HEW9LbPS1LbVneLFFjf2UZ5UfyOuyMQ5qYI7xw"
	DownloadAndStorePublicKeyForTokenCached(certsUrl, tokenString, "public_key")
	DownloadAndStorePublicKeyForTokenCached(certsUrl, tokenString, "public_key")
	DownloadAndStorePublicKeyForTokenCached(certsUrl, tokenString, "public_key")
}

var cache = ""

func DownloadAndStorePublicKeyForTokenCached(certsUrl, tokenString, outputpath string) string {
	if cache == "" {
		cache = DownloadAndStorePublicKeyForToken(certsUrl, tokenString, outputpath)
	}
	return cache
}

func DownloadAndStorePublicKeyForToken(keycloakUrl, tokenString, outputpath string) string {
	fileContent, err := os.ReadFile(outputpath)
	if err == nil {
		return string(fileContent)
	}
	key := DownloadPublicKeyForToken(keycloakUrl, tokenString)
	os.WriteFile(outputpath, []byte(key), 0644)
	return key
}

func FileExists(name string) bool {
	_, err := os.Stat(name)
	if err == nil {
		return true
	}
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	return false
}

func DownloadPublicKeyForToken(certsUrl, tokenString string) string {
	kid := decodeJwtHeader(tokenString).Kid

	// Step 2: Load certs
	keys := loadKeys(certsUrl)

	// Step 3: Extract the correct cert for a kid
	key := findKeyWithKid(keys, kid)

	// Step 4: Write out the cert with pre- and postfix
	key = formatKey(key)

	return key
}

func decodeJwtHeader(tokenString string) JwtHeader {
	parts := strings.Split(tokenString, ".")
	headerBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		log.Fatalf("Could not decode header-part of token-base64: %v", err)
	}

	header := JwtHeader{}
	err = json.Unmarshal(headerBytes, &header)
	if err != nil {
		log.Fatalf("Failed to parse jwt-header: %v", err)
	}

	return header
}

func loadKeys(url string) keyResponse {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("could not query certs endpoint (GET-Http-request): %v", err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Could not decode server-key-response: %v", err)
	}

	keys := keyResponse{}
	json.Unmarshal(body, &keys)

	return keys
}

func findKeyWithKid(keys keyResponse, kid string) string {
	for _, key := range keys.Keys {
		if key.Kid == kid {
			return key.X5c[0]
		}
	}
	return ""
}

func formatKey(raw string) string {
	return "-----BEGIN CERTIFICATE-----\n" + strings.Join(chunks(raw, 64), "\n") + "\n-----END CERTIFICATE-----\n"
}

// this code comes from: https://stackoverflow.com/questions/25686109/split-string-by-length-in-golang
func chunks(s string, chunkSize int) []string {
	if len(s) == 0 {
		return nil
	}
	if chunkSize >= len(s) {
		return []string{s}
	}
	var chunks []string = make([]string, 0, (len(s)-1)/chunkSize+1)
	currentLen := 0
	currentStart := 0
	for i := range s {
		if currentLen == chunkSize {
			chunks = append(chunks, s[currentStart:i])
			currentLen = 0
			currentStart = i
		}
		currentLen++
	}
	chunks = append(chunks, s[currentStart:])
	return chunks
}

type keyResponse struct {
	Keys []struct {
		Kid string   `json:"kid"`
		X5c []string `json:"x5c"`
	} `json:"keys"`
}

type JwtHeader struct {
	Alg string `json:"alg,omitempty"`
	Typ string `json:"typ,omitempty"`
	Kid string `json:"kid,omitempty"`
}

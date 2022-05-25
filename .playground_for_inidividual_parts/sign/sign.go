package main

import (
	"io/ioutil"
	"log"
	"strings"

	"github.com/golang-jwt/jwt"
)

func main() {
	tokenString := "eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJWb2RoUUFoNGFkU2MwejJTcl9CamdyWHAzT2xmbXJ6YzI0Vjlnc0w4a3VJIn0.eyJleHAiOjE2NTM0MTQ4NTUsImlhdCI6MTY1MzQxNDU1NSwianRpIjoiZDU5Mjc4N2YtZmQ0ZS00MTdmLWEzZDQtYTliNTUzYzAxN2FjIiwiaXNzIjoiaHR0cDovL2xvY2FsaG9zdDo4MDgwL3JlYWxtcy9za3l0YWxhIiwiYXVkIjoiYWNjb3VudCIsInN1YiI6IjQxZDE1ODFiLTQwOTAtNDU1OS04MzMxLWMzMDdhZjdhNWI1YiIsInR5cCI6IkJlYXJlciIsImF6cCI6ImNsaWVudGlkLTAzIiwic2Vzc2lvbl9zdGF0ZSI6IjQzYWQ2OTRiLWE5ODQtNDFkNi1hYjI3LTE0NmQwNThhNjBiOCIsImFjciI6IjEiLCJhbGxvd2VkLW9yaWdpbnMiOlsiaHR0cDovL2xvY2FsaG9zdDo4MDgwIl0sInJlYWxtX2FjY2VzcyI6eyJyb2xlcyI6WyJkZWZhdWx0LXJvbGVzLXNreXRhbGEiLCJvZmZsaW5lX2FjY2VzcyIsInVtYV9hdXRob3JpemF0aW9uIl19LCJyZXNvdXJjZV9hY2Nlc3MiOnsiYWNjb3VudCI6eyJyb2xlcyI6WyJtYW5hZ2UtYWNjb3VudCIsIm1hbmFnZS1hY2NvdW50LWxpbmtzIiwidmlldy1wcm9maWxlIl19fSwic2NvcGUiOiJwcm9maWxlIGVtYWlsIiwic2lkIjoiNDNhZDY5NGItYTk4NC00MWQ2LWFiMjctMTQ2ZDA1OGE2MGI4IiwiZW1haWxfdmVyaWZpZWQiOnRydWUsIm5hbWUiOiJUaG9tYXMgRGVuaWZmZWwiLCJwcmVmZXJyZWRfdXNlcm5hbWUiOiJ0LmRlbmlmZmVsIiwiZ2l2ZW5fbmFtZSI6IlRob21hcyIsImZhbWlseV9uYW1lIjoiRGVuaWZmZWwiLCJlbWFpbCI6InQuZGVuaWZmZWxAc2t5dGFsYS1nbWJoLmNvbSJ9.Zma30kGohKBxageFVOTW64guXqYxnLrPH3bDNY6lSdzlkUqhjFJj4L0yok6x4p2b73ufWB7r249_bNy3ypl1KP2FlScfzZMlGmVGlwLOkgWt1m9JFtZ55ofJ0HdnDUktHrveczolECIii4r-mZb-bpaIaDHfX3wVggOv7D3ba_kgcOmFGChR-3_7YIj_YX4yZdHduZ3GYvHT-BcAlcVIsI9KzpvsGZb7hdCU6dB8ssQq7U9ogmSMcbWY-giLssUJnDSR4B8ibwMCpuQEwQjghIdzE6BqQUL2rRQHbWbr7PABgvVcg0URaleEKzF7Wvrkwg_UneiFGlDwd8JyImBr6w"
	keyData, err := ioutil.ReadFile("public_key")
	if err != nil {
		log.Fatalf("Could not open public key: %v", err)
	}

	valid, err := verifyToken(tokenString, keyData)
	if err != nil {
		log.Fatalf("token verification failed: %v", err)
	}
	if valid {
		log.Println("Token is valid")
	} else {
		log.Println("Token is not valid")
	}
}

func verifyToken(token string, keyData []byte) (bool, error) {
	key, err := jwt.ParseRSAPublicKeyFromPEM(keyData)
	if err != nil {
		return false, err
	}
	parts := strings.Split(token, ".")
	// TODO: check parts length
	err = jwt.SigningMethodRS256.Verify(strings.Join(parts[0:2], "."), parts[2], key)
	if err != nil {
		return false, nil
	}
	return true, nil
}

// Obtain certificate from
// http://localhost:8080/realms/skytala/protocol/openid-connect/certs

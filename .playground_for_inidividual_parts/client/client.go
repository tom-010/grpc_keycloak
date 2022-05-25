package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"path/filepath"
	"time"

	"github.com/Nerzal/gocloak"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "skytala.com/jwt-example/proto"
)

const (
	address  = "localhost:50051"
	hostname = "localhost"

	keycloakUrl  = "http://localhost:8080"
	realm        = "skytala"
	clientId     = "clientid-03"
	clientSecret = "IslMJ7tUpEEWAjN6osCD5RLqmIZVsMnM"
)

type KeycloakTokenAuth struct {
	basePath     string
	realm        string
	clientID     string
	clientSecret string
	client       gocloak.GoCloak
	token        *gocloak.JWT
}

// NewKeycloakTokenAuth creates a KeycloakTokenAuth.
func NewKeycloakTokenAuth(
	basePath, realm, clientID, clientSecret string,
) (*KeycloakTokenAuth, error) {
	t := &KeycloakTokenAuth{
		basePath:     basePath,
		realm:        realm,
		clientID:     clientID,
		clientSecret: clientSecret,
	}
	t.client = gocloak.NewClient(t.basePath)
	// clientID string, clientSecret string, realm string, username string, password string

	file, err := ioutil.ReadFile("token.json")
	if err != nil {
		log.Fatalf("could not open token file: %v", err)
	}
	token := new(gocloak.JWT)
	json.Unmarshal(file, token)

	if err != nil {
		log.Fatalf("LoginClient: %v", err)
	}
	t.token = token
	return t, nil
}

func (t KeycloakTokenAuth) GetRequestMetadata(
	ctx context.Context, in ...string,
) (map[string]string, error) {

	file, err := ioutil.ReadFile("token.json")
	if err != nil {
		log.Fatalf("could not open token file: %v", err)
	}
	token := new(gocloak.JWT)
	json.Unmarshal(file, token)

	return map[string]string{
		"authorization": "Bearer " + t.token.AccessToken,
	}, nil
}

func (KeycloakTokenAuth) RequireTransportSecurity() bool {
	return true
}

func main() {
	ta, err := NewKeycloakTokenAuth(keycloakUrl, realm, clientId, clientSecret)
	if err != nil {
		log.Fatalf("Keycloak error: %v", err)
	}

	crtFile := filepath.Join("server.crt")
	creds, err := credentials.NewClientTLSFromFile(crtFile, hostname)
	if err != nil {
		log.Fatalf("failed to load credentials: %v", err)
	}

	opts := []grpc.DialOption{
		grpc.WithPerRPCCredentials(ta),
		// transport credentials.
		grpc.WithTransportCredentials(creds),
	}

	// Set up a connection to the server.
	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()

	c := pb.NewUserManagementClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.CreateNewUser(ctx, &pb.NewUser{Name: "Tom", Age: 22})
	if err != nil {
		log.Fatalf("Could not create user: %v", err)
	}
	log.Printf("Created user with id=%d", r.Id)
}

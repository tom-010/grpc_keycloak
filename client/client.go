package main

import (
	"context"
	"log"
	"path/filepath"
	"time"

	pb "deniffel.com/grpc_keycloak/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	address      = "localhost:50051"
	loginAddress = "localhost:50052"
	hostname     = "localhost"
	crtFile      = "server.crt"
)

func login(username, password string) (string, error) {
	crtFileContent := filepath.Join(crtFile)
	creds, err := credentials.NewClientTLSFromFile(crtFileContent, hostname)
	if err != nil {
		log.Fatalf("failed to load credentials: %v", err)
	}

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
	}

	conn, err := grpc.Dial(loginAddress, opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewLoginServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.Login(ctx, &pb.LoginData{Username: username, Password: password})
	if err != nil {
		return "", err
	}

	if !r.Ok {
		return "", err
	}

	return r.Token, nil
}

func main() {
	token, err := login("t.deniffel", "password")
	log.Println(token)

	crtFileContent := filepath.Join(crtFile)
	creds, err := credentials.NewClientTLSFromFile(crtFileContent, hostname)
	if err != nil {
		log.Fatalf("failed to load credentials: %v", err)
	}

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
	}

	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewUserManagementClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var new_users = make(map[string]int32)
	new_users["Alice"] = 43
	new_users["Bob"] = 30

	for name, age := range new_users {
		r, err := c.CreateNewUser(ctx, &pb.NewUser{Name: name, Age: age})
		if err != nil {
			log.Fatalf("could not create user %v", err)
		}
		log.Printf(`User Details:
		NAME: %s,
		AGE: %d,
		ID: %d`, r.GetName(), r.GetAge(), r.GetId())
	}
}

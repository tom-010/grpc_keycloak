package main

import (
	"context"
	"crypto/tls"
	"log"
	"net"

	pb "deniffel.com/grpc_keycloak/proto"
	"github.com/Nerzal/gocloak/v11"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	port    = ":50052"
	crtFile = "server.crt"
	keyFile = "server.key"

	realm           = "skytala"
	masterRealm     = "master"
	clientID        = "clientid-03"
	clientSecret    = "La59FL56BdLR9vBmzrGRXptk0HYfLwxT"
	keycloakUrl     = "http://localhost:8080"
	initialPassword = "password"
)

type LoginServer struct {
	pb.UnimplementedLoginServiceServer
}

func (s *LoginServer) Login(ctx context.Context, in *pb.LoginData) (*pb.LoginResult, error) {
	client := gocloak.NewClient(keycloakUrl, gocloak.SetAuthAdminRealms("admin/realms"), gocloak.SetAuthRealms("realms"))
	restyClient := client.RestyClient()
	restyClient.SetDebug(false)
	restyClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

	jwt, err := client.Login(ctx, clientID, clientSecret, realm, in.Username, in.Password)
	if err != nil {
		log.Printf("could not login at keycloak: %v", err)
		return &pb.LoginResult{Ok: false}, nil
	}

	return &pb.LoginResult{Ok: true, Token: jwt.AccessToken}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	cert, err := tls.LoadX509KeyPair(crtFile, keyFile)
	if err != nil {
		log.Fatalf("failed to load key pair: %s", err)
	}

	opts := []grpc.ServerOption{
		grpc.Creds(credentials.NewServerTLSFromCert(&cert)),
	}

	s := grpc.NewServer(opts...)
	pb.RegisterLoginServiceServer(s, &LoginServer{})

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

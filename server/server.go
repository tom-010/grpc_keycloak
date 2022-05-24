package main

import (
	"context"
	"crypto/tls"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"strings"

	pb "deniffel.com/grpc_keycloak/proto"
	"github.com/golang-jwt/jwt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	port    = ":50051"
	crtFile = "server.crt"
	keyFile = "server.key"
)

var (
	errMissingMetadata = status.Errorf(codes.InvalidArgument, "missing metadata")
	errInvalidToken    = status.Errorf(codes.Unauthenticated, "invalid credentials")
)

type UserManagementServer struct {
	pb.UnimplementedUserManagementServer
}

func (s *UserManagementServer) CreateNewUser(ctx context.Context, in *pb.NewUser) (*pb.User, error) {
	log.Printf("Received: %v", in.GetName())
	var user_id int32 = int32(rand.Intn(100000))
	return &pb.User{Name: in.GetName(), Age: in.GetAge(), Id: user_id}, nil
}

type LoginServer struct {
	pb.UnimplementedLoginServiceServer
}

func (s *LoginServer) Login(ctx context.Context, in *pb.LoginData) (*pb.LoginResult, error) {
	log.Printf("Login: %v", in)
	return &pb.LoginResult{Ok: false}, nil
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
		grpc.UnaryInterceptor(ensureValidJwtCredentials),
	}

	s := grpc.NewServer(opts...)
	pb.RegisterUserManagementServer(s, &UserManagementServer{})

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func ensureValidJwtCredentials(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errMissingMetadata
	}

	if !valid(md["authorization"]) {
		return nil, errInvalidToken
	}

	res, err := handler(ctx, req)

	return res, err
}

func valid(authorization []string) bool {
	if len(authorization) < 1 {
		return false
	}

	token := strings.TrimPrefix(authorization[0], "Bearer ") // no not forget the whitespace
	log.Println(token)

	keyData, err := ioutil.ReadFile("public_key")
	if err != nil {
		log.Fatalf("Could not open public key: %v", err)
	}

	key, err := jwt.ParseRSAPublicKeyFromPEM(keyData)
	if err != nil {
		log.Printf("could not parse rsa key: %v", err)
		return false
	}
	parts := strings.Split(token, ".")
	// TODO: check parts length
	err = jwt.SigningMethodRS256.Verify(strings.Join(parts[0:2], "."), parts[2], key)
	if err != nil {
		log.Println("was not valid")
		return false
	}
	return true
}

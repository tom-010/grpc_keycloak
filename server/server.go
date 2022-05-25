package main

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
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
	creator, err := getUser(ctx)
	if err != nil {
		return nil, err
	}
	var user_id int32 = int32(rand.Intn(100000))
	return &pb.User{Name: in.GetName(), Age: in.GetAge(), Id: user_id, CreatedBy: creator.Name}, nil
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

func getUser(ctx context.Context) (User, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return User{}, errMissingMetadata
	}

	var authorization = md["authorization"]
	if len(authorization) < 1 {
		return User{}, errMissingMetadata
	}

	token := strings.TrimPrefix(authorization[0], "Bearer ") // no not forget the whitespace

	return toUser(token), nil
}

func toUser(tokenString string) User {
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

	return user
}

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

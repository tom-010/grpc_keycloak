grpc_keycloak
=============

An example how to use keycloak to secure gRPC communication.

## Setup Keycoak

```bash
./scripts/start_keycloak.sh
```

1. Go to [http://localhost:8080](http://localhost:8080) > Administration Console > "admin" "admion" > Login
2. Hover "Master" in the left upper corner > "Add Realm" > Name="skytala" > Create
3. Left Menu "Clients" > Right upper corner "Create" > Client ID="clientid-03", Root URL = "http://localhost:8080/" > Save
4. Access Type="confidential" > Save
5. Tab Menu "Credentials" > Copy the secret and use it in the app

## Generate Certificates

```bash
./scripts/create_cert.sh
```


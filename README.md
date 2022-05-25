grpc_keycloak
=============

An example how to use keycloak to secure gRPC communication.

## THIS REPO IS A PROTOTYPE

While it should be secure, it has a very bad code-style and does a lot of work twice, as it is a showcase
and optimized for that and not for maintainability. It also lacks any unit-testing or production-experience.

## Setup Keycoak

```bash
./scripts/start_keycloak.sh
```

1. Go to [http://localhost:8080](http://localhost:8080) > Administration Console > "admin" "admion" > Login
2. Hover "Master" in the left upper corner > "Add Realm" > Name="skytala" > Create
3. Left Menu "Clients" > Right upper corner "Create" > Client ID="clientid-03", Root URL = "http://localhost:8080/" > Save
4. Access Type="confidential" > Save
5. Tab Menu "Credentials" > Copy the secret and use it in the app

Now, create a user:

1. Make sure, you are in the realm `skytala` and the client `clientid-03`
2. Left menu "Users" > Right upper corner "Create User" > "Username"="m.mustermann", "First Name"="Max", "Last Name"="Mustermann", "User Enabled"="True", "E-Mail"="m.mustermann@email.com" > save
3. In tabs "Credentials" > "Password"="password", "Confirm Password"="password", "Temporary"="Off" > "Set Password"

## Generate Certificates

```bash
./scripts/create_cert.sh
```

## Add the keycloak public-key

Visit 

```bash
realm="skytala"
curl http://localhost:8080/realms/$realm/protocol/openid-connect/certs
```

To get the public-keys of the server. There should be multiple and everyone has its own `kid`.

Use
```bash
./scripts/fetch_token.sh
```
to fetch a token (do not forget to change the client-secret!).

Decrypt the token at [jwt.io](jwt.io) and look out for the `kid` in the header-part. Now you know
which certificate you need.

Create a new file `public_key` in the root of this project with the content:

```
-----BEGIN CERTIFICATE-----
<<paste your copied token here>>
-----END CERTIFICATE-----
```

As an alternative, take a look at `tools/public_key/public_key.go`.

## Start everything up

You already started keycloak. In three other terminals, run:
```bash
./scripts/run_login_server.sh
./scripts/run_server.sh
./scripts/run_client.sh
```

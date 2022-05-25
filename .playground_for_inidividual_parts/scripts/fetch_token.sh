
realm="skytala"
client="clientid-03"
secret="bI6wznsqBH3dN1UWeRk8Yz4Xepp1li7D"

username="t.deniffel"
password="password"

curl -L -X POST "http://localhost:8080/realms/$realm/protocol/openid-connect/token" \
-H "Content-Type: application/x-www-form-urlencoded" \
--data-urlencode "client_id=$client" \
--data-urlencode "grant_type=password" \
--data-urlencode "client_secret=$secret" \
--data-urlencode "scope=openid" \
--data-urlencode "username=$username" \
--data-urlencode "password=$password" \
> token.json

echo $?

cat token.json
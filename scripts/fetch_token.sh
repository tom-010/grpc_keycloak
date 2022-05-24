
realm="skytala"
client="clientid-03"
secret="La59FL56BdLR9vBmzrGRXptk0HYfLwxT"

echo "$secret"
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

code=$?

echo $code

cat token.json

exit $code
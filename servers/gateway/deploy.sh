sh build.sh

ssh ec2-user@api.jimhua32.me << EOF

docker rm -f gateway
docker rm -f redis

docker pull jimhua32/finalprojgateway

docker run \
    -d \
    --name redis \
    --network serv \
    redis

docker run \
    -d \
    -e ADDR=:443 \
    -e SESSIONKEY="key" \
    -e MYSQL_ROOT_PASSWORD="blah" \
    -e MYSQL_DATABASE="mydatabase" \
    -e REDISADDR=redis:6379 \
    -e DSN="root:blah@tcp(database:3306)/mydatabase?parseTime=true" \
    -e MEETINGADDR="meetings:5000" \
    -v /etc/letsencrypt:/etc/letsencrypt:ro \
    -e TLSCERT=/etc/letsencrypt/live/api.jimhua32.me/fullchain.pem \
    -e TLSKEY=/etc/letsencrypt/live/api.jimhua32.me/privkey.pem \
    -p 443:443 \
    --name gateway \
	--network serv \
    jimhua32/finalprojgateway
exit

EOF
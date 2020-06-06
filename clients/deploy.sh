sh build.sh

ssh ec2-user@jimhua32.me << EOF

docker rm -f client

docker pull jimhua32/finalprojclient

docker run \
    -d \
    -e ADDR=:443 \
    -v /etc/letsencrypt:/etc/letsencrypt:ro \
    -e TLSCERT=/etc/letsencrypt/live/jimhua32.me/fullchain.pem \
    -e TLSKEY=/etc/letsencrypt/live/jimhua32.me/privkey.pem \
    -p 443:443 \
    -p 80:80 \
    --name client \
    jimhua32/finalprojclient
exit

EOF
sh build.sh

ssh ec2-user@api.jimhua32.me << EOF

docker rm -f database

docker pull jimhua32/finalprojdb

docker run \
    -d \
    -e MYSQL_ROOT_PASSWORD="blah" \
    -e MYSQL_DATABASE="mydatabase" \
    --name database \
	--network serv \
    jimhua32/finalprojdb

EOF
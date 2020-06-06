sh build.sh

ssh ec2-user@api.jimhua32.me << EOF

docker rm -f meetings

docker pull jimhua32/meetings

docker run \
    -d \
    -e ADDR=:5000 \
    -e DSN="root:blah@tcp(database:3306)/mydatabase?parseTime=true" \
    --name meetings \
	--network serv \
    jimhua32/meetings
exit

EOF
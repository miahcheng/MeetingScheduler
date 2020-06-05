GOOS=linux go build
docker build -t jimhua32/meetings .
go clean

docker push jimhua32/meetings
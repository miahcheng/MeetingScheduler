GOOS=linux go build
docker build -t jimhua32/finalprojgateway .
go clean

docker push jimhua32/finalprojgateway
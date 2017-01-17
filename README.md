```
./cockroach version
Build Tag:   beta-20170112
Build Time:  2017/01/12 18:27:36
Platform:    linux amd64
Go Version:  go1.7.3
C Compiler:  gcc 4.9.2
./cockroach start --background
./cockroach start --store=cockroach-data2 --port=26258 --http-port=8081 --join=localhost:26257 --background
./cockroach start --store=cockroach-data3 --port=26259 --http-port=8082 --join=localhost:26257 --background
./cockroach sql -e "create database bench;"

go get github.com/lib/pq
go run ./fill.go
go test -test.run=^$ -test.bench=Bench .
```

@echo off
echo Building Airbnb CLI tools...

echo Building main CLI...
go build -o airbnb-cli.exe cmd/main.go

echo Building producer...
go build -o airbnb-producer.exe cmd/producer/main.go

echo Building consumer...
go build -o airbnb-consumer.exe cmd/consumer/main.go

echo Build completed. 
@echo off
cd api

echo Setting up Go modules...
go mod init recommendation-engine/api
go get github.com/lib/pq
go get github.com/redis/go-redis/v9
go mod tidy

echo Go modules setup complete!
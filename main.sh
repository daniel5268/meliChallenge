migrate -database $DATABASE_URL -path ./migrations up
go run ./src/main.go
$connection="postgresql://postgres:postgres@localhost:5432/final?sslmode=disable"
migrate -path migrations -database "$connection" up

tidy:
	go mode tidy
run:
	go run .
test:
	go test -v ./... -cover -count=1

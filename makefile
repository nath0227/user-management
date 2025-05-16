tidy:
	go mode tidy
run:
	go run .
test:
	go test -v ./... -cover -count=1

docker-gen-proto: del-output
	docker run --volume "./app/user/grpc:/workspace" --workdir /workspace bufbuild-go generate proto

del-output: 
	rm -rf ./app/user/grpc/doc/*
	rm -rf ./app/user/grpc/gen/*
docker-build: 
	docker build -t bufbuild-go .

#Important Add shared paths from Docker -> Preferences... -> Resources -> File Sharing.
#don't shared paths. Please copy command in 7 th line to run in terminal
docker-gen-proto: del-output
	docker run --volume "$(pwd)/grpc:/workspace" --workdir /workspace bufbuild-go generate proto

del-output: 
	rm -rf grpc/doc/*
	rm -rf grpc/gen/*
.PHONY: proto

proto: # compile .proto files (output: ./proto/*.pb.go)
	docker build -t protogen -f Dockerfile.protogen .
	docker run --name protogen protogen
	docker cp protogen:/proto/gen/echo.pb.go ./echo
	docker rm protogen

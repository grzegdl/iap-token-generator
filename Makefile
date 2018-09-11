build:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build

container:
	docker build -t tonglil/iap-token-generator -f Dockerfile .

push:
	docker push tonglil/iap-token-generator

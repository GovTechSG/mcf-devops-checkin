dep:
	GO111MODULE=on go mod vendor

build:
	GO111MODULE=on go build -mod vendor -o ./bin/pinger ./cmd/pinger

run:
	GO111MODULE=on go run ./cmd/pinger

test:
	GO111MODULE=on go test ./...

testenv:
	# do not change this
	docker-compose up -f ./deployments/docker-compose.yml

docker_image:
	# do not change this
	docker build -f ./build/Dockerfile -t devops/pinger:latest .

docker_testrun:
	# do not change this
	docker run -it -p 8000:8000 devops/pinger:latest

docker_tar: docker_image
	# do not change this
	docker save -o ./build/pinger.tar devops/pinger:latest 

docker_untar:
	# do not change this
	docker load -i ./build/pinger.tar

export GO111MODULE=on

.PHONY: build
build:
	@echo "-- building binary"
	go build -o ./bin/manyface cmd/server/main.go

.PHONY: dev
dev: 
	@echo "-- starting air wrapper"
	air -c ./air.toml

.PHONY: docker_dev
docker_dev: 
	@echo "-- building docker dev container "
	docker build -f build/Dockerfile_golang1.17.5 -t olaesean/manyface:golang1.17.5 .

.PHONY: docker_prod
docker_prod: 
	@echo "-- building docker prod container"
	docker build -f build/Dockerfile4_scratch_multistage -t olaesean/manyface .

.PHONY: docker_run
docker_run: 
	@echo "-- starting docker dev container"
	docker run -d -p 8080:8080 olaesean/manyface:golang1.17.5

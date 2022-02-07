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
	docker push olaesean/manyface:golang1.17.5

.PHONY: docker_prod
docker_prod: 
	@echo "-- building docker prod container"
	docker build -f build/Dockerfile4_scratch_multistage -t olaesean/manyface .
	docker push olaesean/manyface

.PHONY: docker_run
docker_run: 
	@echo "-- starting docker dev container"
	docker run -d -p 8080:8080 olaesean/manyface:golang1.17.5

.PHONY: mount
mount: 
	@echo "-- mounting local directory into minikube"
	minikube mount ./db:/data

.PHONY: copydb
copydb:
	@echo "-- coping db file into PersistentVolume"
	kubectl cp ./db/data.db test:/data

.PHONY: deploy
deploy:
	@echo "-- deployment k8s objects into kube"
	kubectl delete deploy manyface
	sleep 4
	# kubectl apply -f deployments/pvc.yaml
	kubectl apply -f deployments/configmap.yaml
	kubectl apply -f deployments/pod.yaml
	# sleep 4
	# kubectl cp ./db/data.db test:/data
	kubectl apply -f deployments/clusterip.yaml
	kubectl apply -f deployments/deployment.yaml
	kubectl get po

.PHONY: start_mtrx
start_mtrx:
	cd ../synapse
	source venv/bin/activate
	synctl start

.PHONY: start_conduit
start_conduit:
docker run -d -p 8009:6167 -v ~/Github/matrix/conduit/.conduit/conduit.toml:/srv/conduit/conduit.toml -v ~/Github/matrix/conduit/.conduit/db:/srv/conduit/.local/share/conduit matrixconduit/matrix-conduit:latest
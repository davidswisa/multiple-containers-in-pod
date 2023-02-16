.PHONY: run-compose run build clean apply

run-compose:
	docker-compose up -d --build

apply:
	kubectl apply -f ./kube

run: build apply
	

build:
	docker rmi $$(docker images | grep 'multiple-containers-in-pod') || true
	docker-compose build --no-cache

.PHONY: rebuild-%
rebuild-%:
	kubectl delete deployment "$*" || true
	# sleep 15
	./scripts/wait-for-pods.sh "$*"
	docker rmi "multiple-containers-in-pod_$*"
	docker-compose build --no-cache "$*"
	$(MAKE) apply


clean: kube-clean
	docker rmi $$(docker images | grep 'multiple-containers-in-pod') || true

kube-clean:
	kubectl delete deployments --all
	kubectl delete services --all
	kubectl delete pods --all
	kubectl delete daemonset --all
	kubectl delete NetworkPolicy --all
	kubectl delete replicaset --all


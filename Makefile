export APDIR=$(go list ./... | grep -v /vendor/)

build:
	CGO_ENABLED=0 go install -tags netgo ${APDIR}

run: build
	${GOPATH}/bin/broker

run-local: build
	PORT=8181 KUBERNETES_URL=127.0.0.1:8080 ${GOPATH}/bin/broker


build-docker: build
	cp ${GOPATH}/bin/broker ./broker.elf
	docker build -t tap-broker .

push-latest: build-docker
	docker tag tap-broker  5edf9636-8df5-4609-b93b-4af5e0b1d81f.tmp.us.enableiot.com:5000/tap-broker:latest
	docker push 5edf9636-8df5-4609-b93b-4af5e0b1d81f.tmp.us.enableiot.com:5000/tap-broker:latest

push-stable: build-docker
	@echo "is it stable?"
	exit 1
	docker tag tap-broker  5edf9636-8df5-4609-b93b-4af5e0b1d81f.tmp.us.enableiot.com:5000/tap-broker:0.8.1
	docker push 5edf9636-8df5-4609-b93b-4af5e0b1d81f.tmp.us.enableiot.com:5000/tap-broker:0.8.1

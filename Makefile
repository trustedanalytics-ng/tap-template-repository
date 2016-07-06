GOBIN=$(GOPATH)/bin
APP_DIR_LIST=$(shell go list ./... | grep -v /vendor/)

build: verify_gopath
	CGO_ENABLED=0 go install -tags netgo $(APP_DIR_LIST)
	go fmt $(APP_DIR_LIST)

verify_gopath:
	@if [ -z "$(GOPATH)" ] || [ "$(GOPATH)" = "" ]; then\
		echo "GOPATH not set. You need to set GOPATH before run this command";\
		exit 1 ;\
	fi

docker_build: build
	rm -Rf application && mkdir application
	cp -Rf $(GOBIN)/tapng-template-repository application/
	docker build -t tapng-template-repository .

push_docker: docker_build
	docker tag tapng-template-repository $(REPOSITORY_URL)/tapng-template-repository:latest
	docker push $(REPOSITORY_URL)/tapng-template-repository:latest

kubernetes_deploy:
	kubectl create -f configmap.yaml
	kubectl create -f service.yaml
	kubectl create -f deployment.yaml

deps_fetch_newest:
	$(GOBIN)/govendor remove +all
	@echo "Update deps used in project to their newest versions"
	$(GOBIN)/govendor fetch -v +external, +missing

deps_update: verify_gopath
	$(GOBIN)/govendor remove +all
	$(GOBIN)/govendor add +external
	@echo "Done"

bin/govendor: verify_gopath
	go get -v -u github.com/kardianos/govendor

tests: verify_gopath
	go test --cover $(APP_DIR_LIST)

prepare_dirs:
	mkdir -p ./temp/src/github.com/trustedanalytics/tapng-template-repository
	$(eval REPOFILES=$(shell pwd)/*)
	ln -sf $(REPOFILES) temp/src/github.com/trustedanalytics/tapng-template-repository

build_anywhere: prepare_dirs
	$(eval GOPATH=$(shell cd ./temp; pwd))
	$(eval APP_DIR_LIST=$(shell GOPATH=$(GOPATH) go list ./temp/src/github.com/trustedanalytics/tapng-template-repository/... | grep -v /vendor/))
	GOPATH=$(GOPATH) CGO_ENABLED=0 go install -tags netgo $(APP_DIR_LIST)
	rm -Rf application && mkdir application
	cp $(GOPATH)/bin/tapng-template-repository ./application/tapng-template-repository

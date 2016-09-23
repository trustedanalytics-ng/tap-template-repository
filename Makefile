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

run: build_anywhere
	./application/tap-template-repository

run-local: build
	PORT=8083 TEMPLATE_REPOSITORY_USER=admin TEMPLATE_REPOSITORY_PASS=password ${GOPATH}/bin/tap-template-repository

docker_build: build_anywhere
	docker build -t tap-template-repository .

push_docker: docker_build
	docker tag -f tap-template-repository $(REPOSITORY_URL)/tap-template-repository:latest
	docker push $(REPOSITORY_URL)/tap-template-repository:latest

kubernetes_deploy: docker_build
	kubectl create -f configmap.yaml
	kubectl create -f service.yaml
	kubectl create -f deployment.yaml

kubernetes_update: docker_build
	kubectl delete -f deployment.yaml
	kubectl create -f deployment.yaml

deps_fetch_specific: bin/govendor
	@if [ "$(DEP_URL)" = "" ]; then\
		echo "DEP_URL not set. Run this comand as follow:";\
		echo " make deps_fetch_specific DEP_URL=github.com/nu7hatch/gouuid";\
	exit 1 ;\
	fi
	@echo "Fetching specific dependency in newest versions"
	$(GOBIN)/govendor fetch -v $(DEP_URL)

deps_update_tap: verify_gopath
	$(GOBIN)/govendor update github.com/trustedanalytics/...
	rm -Rf vendor/github.com/trustedanalytics/tap-template-repository
	@echo "Done"

prepare_dirs:
	mkdir -p ./temp/src/github.com/trustedanalytics/tap-template-repository
	$(eval REPOFILES=$(shell pwd)/*)
	ln -sf $(REPOFILES) temp/src/github.com/trustedanalytics/tap-template-repository

build_anywhere: prepare_dirs
	$(eval GOPATH=$(shell cd ./temp; pwd))
	$(eval APP_DIR_LIST=$(shell GOPATH=$(GOPATH) go list ./temp/src/github.com/trustedanalytics/tap-template-repository/... | grep -v /vendor/))
	GOPATH=$(GOPATH) CGO_ENABLED=0 go build -tags netgo $(APP_DIR_LIST)
	rm -Rf application && mkdir application
	cp -RL ./tap-template-repository ./application/tap-template-repository
	rm -Rf ./temp

install_mockgen:
	scripts/install_mockgen.sh

mock_update: install_mockgen
	$(GOBIN)/mockgen -source=catalog/template.go -package=catalog -destination=catalog/template_mock.go

test: verify_gopath  mock_update
	go test --cover $(APP_DIR_LIST)

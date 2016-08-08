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

docker_build: build_anywhere
	docker build -t tapng-template-repository .

push_docker: docker_build
	docker tag tapng-template-repository $(REPOSITORY_URL)/tapng-template-repository:latest
	docker push $(REPOSITORY_URL)/tapng-template-repository:latest

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

deps_update_tapng: verify_gopath
	$(GOBIN)/govendor update github.com/trustedanalytics/...
	rm -Rf vendor/github.com/trustedanalytics/tapng-template-repository
	@echo "Done"

bin/gomock: verify_gopath
	go get -v -u github.com/golang/mock/mockgen

mock_update: bin/gomock
	$(GOBIN)/mockgen -source=catalog/template.go -package=catalog -destination=catalog/template_mock.go

tests: verify_gopath mock_update
	go test --cover $(APP_DIR_LIST)

prepare_dirs:
	mkdir -p ./temp/src/github.com/trustedanalytics/tapng-template-repository
	$(eval REPOFILES=$(shell pwd)/*)
	ln -sf $(REPOFILES) temp/src/github.com/trustedanalytics/tapng-template-repository

build_anywhere: prepare_dirs
	$(eval GOPATH=$(shell cd ./temp; pwd))
	$(eval APP_DIR_LIST=$(shell GOPATH=$(GOPATH) go list ./temp/src/github.com/trustedanalytics/tapng-template-repository/... | grep -v /vendor/))
	GOPATH=$(GOPATH) CGO_ENABLED=0 go build -tags netgo $(APP_DIR_LIST)
	rm -Rf application && mkdir application
	cp -RL ./tapng-template-repository ./application/tapng-template-repository
	rm -Rf ./temp

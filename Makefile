GOBIN=$(GOPATH)/bin
APP_DIR_LIST=$(shell go list ./... | grep -v /vendor/)

bin/app: verify_gopath
	CGO_ENABLED=0 go install -tags netgo $(APP_DIR_LIST)
	go fmt $(APP_DIR_LIST)

verify_gopath:
	@if [ -z "$(GOPATH)" ] || [ "$(GOPATH)" = "" ]; then\
		echo "GOPATH not set. You need to set GOPATH before run this command";\
		exit 1 ;\
	fi

build: bin/app
	rm -Rf application && mkdir application
	cp -Rf $(GOBIN)/tap-template-repository application/

push_docker: build
	docker build -t tap-template-repository .
	docker tag tap-template-repository $(REPOSITORY_URL)/tap-template-repository:latest
	docker push $(REPOSITORY_URL)/tap-template-repository:latest

deps_fetch_newest:
	$(GOBIN)/govendor remove +all
	@echo "Update deps used in project to their newest versions"
	$(GOBIN)/govendor fetch -v +external, +missing

deps_update: verify_gopath
	$(GOBIN)/govendor update +external
	@echo "Done"

bin/govendor: verify_gopath
	go get -v -u github.com/kardianos/govendor

tests: verify_gopath
	go test --cover $(APP_DIR_LIST)
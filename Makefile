

.PHONY:	scan
scan: 
	trivy fs . 

.PHONY: build
build: 
	goreleaser build --clean

.PHONY: build-snapshot
build-snapshot: 
	goreleaser build --clean --snapshot --single-target


.PHONY: release-skip-publish
release-skip-publish: 
	goreleaser release --clean --skip-publish  --skip-sign

.PHONY: release-snapshot
release-snapshot: 
	goreleaser release --clean --skip-publish --snapshot --skip-sign
 

.PHONY: lint
lint: 
	golangci-lint run ./... 


.PHONY: changelog
changelog: 
	git-chglog -o CHANGELOG.md 


.PHONY: test
test:
	echo TODO
	


.PHONY: gen-doc
gen-doc: 
	mkdir -p ./doc
	dist/json-patcher_linux_amd64_v1/json-patcher documentation  --dir ./doc



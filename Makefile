gofmt:
	find . -path ./vendor -prune -type f -o -name '*.go' -exec gofmt -d {} + | tee /dev/stderr
	find . -path ./vendor -prune -type f -o -name '*.go' -exec gofmt -w {} + | tee /dev/stderr
test:
	test -z '$(shell find . -path ./vendor -prune -type f -o -name '*.go' -exec gofmt -d {} + | tee /dev/stderr)'
	go test $(shell glide novendor)

build: test
	go build -a -tags netgo -installsuffix netgo -o bin/concourse-aws ./

publish: build
	ghr -u mumoshu -r concourse-aws -c master --prerelease v0.0.2 bin/

publish-latest: build
	ghr -u mumoshu -r concourse-aws -c master --replace --prerelease latest bin/

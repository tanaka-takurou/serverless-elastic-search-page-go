root	:=		$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

.PHONY: clean build deploy

clean:
	rm -rfv bin
	$(MAKE) -C "${root}/api/load" clean
	$(MAKE) -C "${root}/api/save" clean

build:
	mkdir -p bin
	scripts/create_template.sh
	GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o bin/bootstrap
	$(MAKE) -C "${root}/api/load" build
	$(MAKE) -C "${root}/api/save" build

deploy:
	sam package --output-template-file "${root}"/packaged.yml --s3-bucket "${bucket}"
	sam deploy --stack-name "${stack}" --capabilities CAPABILITY_IAM --template-file "${root}/packaged.yml"

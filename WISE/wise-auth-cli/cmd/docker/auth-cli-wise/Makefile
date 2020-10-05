all: build docker-build push

test:
	WISE_INTERNAL_AUTH_REDIRECT="https://aws.us-west-2.internal.wise.us" \
	CONTAINER_LISTEN_PORT=4433 \
	go run .
	
test-envvars:
	@[ "${BUILD_TAG}" ] || ( echo "BUILD_TAG var is not set"; exit 1 )
	@[ "${ECR_IMAGE}" ] || ( echo "ECR_IMAGE var is not set"; exit 1 )
	@[ "${ECS_ENV}" ] || ( echo "ECS_ENV var is not set"; exit 1 )

clean:
	rm -rf static templates main

build: clean
	cp -r ../../../static static
	cp -r ../../../templates templates
	GOOS=linux GOARCH=amd64 go build -mod=vendor -o main .

docker-build: test-envvars
	docker build \
		-t $(ECR_IMAGE):latest \
		-t $(ECR_IMAGE):${BUILD_TAG} \
		.

push: test-envvars
	$$(aws ecr get-login --no-include-email)
	docker push ${ECR_IMAGE}:latest
	docker push ${ECR_IMAGE}:${BUILD_TAG}

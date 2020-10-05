all: build docker-build push deploy

test-envvars:
	@[ "${BUILD_TAG}" ] || ( echo "BUILD_TAG var is not set"; exit 1 )
	@[ "${ECR_IMAGE}" ] || ( echo "ECR_IMAGE var is not set"; exit 1 )
	@[ "${ECS_ENV}" ] || ( echo "ECS_ENV var is not set"; exit 1 )  # dev, qa, stg, prd

clean:
	rm -rf static main

build:
	cp -r ../../../static static
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

deploy: test-envvars
	aws ecs update-service --cluster ${ECS_ENV}-auth-cli --service ${ECS_ENV}-auth-cli-aws --force-new-deployment

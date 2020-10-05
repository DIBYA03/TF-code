export
AWS_DEFAULT_PROFILE=master-us-west-2-saml-roles-deployment
AWS_DEFAULT_REGION=us-west-2
ECR_URL=379379777492.dkr.ecr.us-west-2.amazonaws.com
ECR_REPO=master-core-platform-ecr-shopify-order

all: build docker-build push

test-envvars:
	@[ "${BUILD_NUMBER}" ] || ( echo "BUILD_NUMBER var is not set"; exit 1 )

build:
	GOOS=linux GOARCH=amd64 go build -mod=vendor -o main .

docker-build: test-envvars
	docker build \
		-t $(ECR_URL)/$(ECR_REPO):latest \
		-t $(ECR_URL)/$(ECR_REPO):build${BUILD_NUMBER} \
		.

push: test-envvars
	$$(aws ecr get-login --no-include-email --region us-west-2)
	docker push $(ECR_URL)/$(ECR_REPO):latest
	docker push $(ECR_URL)/$(ECR_REPO):build${BUILD_NUMBER}

export
AWS_DEFAULT_PROFILE=master-us-west-2-saml-roles-deployment
AWS_DEFAULT_REGION=us-west-2
ECR_URL=379379777492.dkr.ecr.us-west-2.amazonaws.com
ECR_REPO=master-core-platform-ecr-app-payments
WEB_PAYMENTS_BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
ALLOWED_WEB_PAYMENTS_BRANCHES=(develop qa staging master)

all: build web-payments docker-build push clean-web-payments

test-envvars:
	@[ "${BUILD_NUMBER}" ] || ( echo "BUILD_NUMBER var is not set"; exit 1 )

clean-web-payments:
	rm -rf web-payments

web-payments-branch:
ifeq (,$(findstring ${WEB_PAYMENTS_BRANCH},$(ALLOWED_WEB_PAYMENTS_BRANCHES)))
  WEB_PAYMENTS_BRANCH=develop
endif

web-payments: clean-web-payments web-payments-branch
	git clone git@github.com:wiseco/web-payments.git
	cd web-payments && git checkout ${WEB_PAYMENTS_BRANCH}

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

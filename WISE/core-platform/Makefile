# Makefile

export
COGNITO_LAMBDA_DIRS=$(shell find cmd/lambda/cognitoauth -type d)
CORE_LAMBDA_DIRS=$(shell find cmd/lambda/clientapi -type d)
UTIL_LAMBDA_DIRS=$(shell find cmd/lambda/util -type d)
CSP_LAMBDA_DIRS=$(shell find cmd/lambda/csp -type d)
DOCKER_DIRS=${shell find cmd/docker/* -name Dockerfile -exec dirname {} \;}
OPENAPI_S3_BUCKET=global-account-resources-mobile-pkg
SPECS_API_DIR=specs/api
GO111MODULE=on
GOPRIVATE=github.com/wiseco

.PHONY: go-mod-vendor

go-mod-vendor:
	go mod vendor

dev-install:
	go get -u golang.org/x/lint/golint

clean:
	cd ../../cmd/lambda/cognitoauth && find ./ -name lambda.zip -exec rm {} \;
	cd ../../cmd/lambda/cognitoauth && find ./ -name main -exec rm {} \;
	cd ../../cmd/lambda/clientapi && find ./ -name lambda.zip -exec rm {} \;
	cd ../../cmd/lambda/clientapi && find ./ -name main -exec rm {} \;
	cd ../../cmd/lambda/util && find ./ -name lambda.zip -exec rm {} \;
	cd ../../cmd/lambda/util && find ./ -name main -exec rm {} \;

vet:
	go vet ./...

lint: vet
	golint $(go list ./... | grep -v /vendor/)

test:
	go test -cover -race ./...

test-stage:
	test -e "$(SPECS_API_DIR)/clientapi/${STAGE}/client-api.yaml"

test-rendered-openapi:
	test -e "$(SPECS_API_DIR)/clientapi/${STAGE}/rendered-client-api.yaml"

clean-openapi:
	rm -f $(SPECS_API_DIR)/clientapi/${STAGE}/client-api.yaml
	rm -f $(SPECS_API_DIR)/clientapi/${STAGE}/rendered-client-api.yaml

test-tf-vars:
	@[ "${STAGE}" ] || ( echo "STAGE var is not set"; exit 1 )
	@[ "${API_ENV}" ] || ( echo "API_ENV var is not set"; exit 1 )
	@[ "${API_REGION}" ] || ( echo "API_REGION var is not set"; exit 1 )

build-go-lambdas: go-mod-vendor
	# cognito lambdas
	$(foreach \
		LAMBDA,$(COGNITO_LAMBDA_DIRS), \
		echo "building $(LAMBDA)"; \
		env GOOS=linux GOARCH=amd64 go build -mod=vendor -o $(LAMBDA)/main $(LAMBDA)/*.go; \
		zip -j $(LAMBDA)/lambda.zip $(LAMBDA)/main; \
	)
	# core lambdas
	$(foreach \
		LAMBDA,$(CORE_LAMBDA_DIRS), \
		echo "building $(LAMBDA)"; \
		env GOOS=linux GOARCH=amd64 go build -mod=vendor -o $(LAMBDA)/main $(LAMBDA)/*.go; \
		zip -j $(LAMBDA)/lambda.zip $(LAMBDA)/main; \
	)
	# util lambdas
	$(foreach \
		LAMBDA,$(UTIL_LAMBDA_DIRS), \
		echo "building $(LAMBDA)"; \
		env GOOS=linux GOARCH=amd64 go build -mod=vendor -o $(LAMBDA)/main $(LAMBDA)/*.go; \
		zip -j $(LAMBDA)/lambda.zip $(LAMBDA)/main; \
	)

build-go-csp-lambdas: go-mod-vendor
	# csp lambdas
	$(foreach \
		LAMBDA,$(CSP_LAMBDA_DIRS), \
		echo "building $(LAMBDA)"; \
		env GOOS=linux GOARCH=amd64 go build -mod=vendor -o $(LAMBDA)/main $(LAMBDA)/*.go; \
		zip -j $(LAMBDA)/lambda.zip $(LAMBDA)/main; \
	)

docker-build: go-mod-vendor
	@[ "${BUILD_NUMBER}" ] || ( echo "BUILD_NUMBER var is not set"; exit 1 )
	$(foreach \
		DOCKERF,$(DOCKER_DIRS), \
		echo "building ${DOCKERF}"; \
		cd ./$(DOCKERF); \
		BUILD_NUMBER=${BUILD_NUMBER} make all; \
		cd -; \
	)

openapi-build: clean-openapi
	cd $(SPECS_API_DIR)/clientapi/${STAGE}/ && yq m client-api-*.yaml > client-api.yaml

terraform-partner-resources: test-tf-vars test-stage
	cd terraform/partner_resources && make ${API_REGION}-${API_ENV} apply-auto-approve

terraform-fargate: test-tf-vars test-stage
	cd terraform/fargate && make ${API_REGION}-${API_ENV} apply-auto-approve

terraform-csp: test-tf-vars test-stage
	cd terraform/csp && make ${API_REGION}-${API_ENV} apply-auto-approve

terraform-client-api-gw: test-tf-vars test-stage
	cd terraform/client_api/api_gateway && make ${API_REGION}-${API_ENV} apply-auto-approve

terraform-client-lambdas: test-tf-vars test-stage
	cd terraform/client_api/lambda_functions && make ${API_REGION}-${API_ENV} apply-auto-approve

# API_ENV=dev API_REGION=us-west-2 STAGE="2019-10-01" make all
all: openapi-build build-go-lambdas terraform-client-api-gw terraform-client-lambdas terraform-fargate

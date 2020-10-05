# Makefile

export
CSP_LAMBDA_DIRS=$(shell find ../../cmd/lambda/csp -type d)
UTIL_LAMBDA_DIRS=$(shell find cmd/lambda/util -type d)
OPENAPI_S3_BUCKET=global-account-resources-mobile-pkg
SPECS_API_DIR="../../specs/api/csp"

dev-install:
	go get -u golang.org/x/lint/golint

clean: clean-openapi
	cd ../../cmd/lambda/csp && find ./ -name lambda.zip -exec rm {} \;
	cd ../../cmd/lambda/csp && find ./ -name main -exec rm {} \;
	cd ../../cmd/lambda/util && find ./ -name lambda.zip -exec rm {} \;
	cd ../../cmd/lambda/util && find ./ -name main -exec rm {} \;

vet:
	go vet ./...

lint: vet
	golint $(go list ./... | grep -v /vendor/)

test:
	go test -cover -race ./...

test-stage:
	test -e "$(SPECS_API_DIR)/${STAGE}/csp-api.yaml"

test-rendered-openapi:
	test -e "$(SPECS_API_DIR)/${STAGE}/rendered-csp-api.yaml"

clean-openapi:
	rm -f $(SPECS_API_DIR)/${STAGE}/csp-api.yaml
	rm -f $(SPECS_API_DIR)/${STAGE}/rendered-csp-api.yaml
	rm -rf $(SPECS_API_DIR)/sdk
	mkdir -p $(SPECS_API_DIR)/sdk

test-tf-vars:
	@[ "${STAGE}" ] || ( echo "STAGE var is not set"; exit 1 )
	@[ "${API_ENV}" ] || ( echo "API_ENV var is not set"; exit 1 )
	@[ "${API_REGION}" ] || ( echo "API_REGION var is not set"; exit 1 )

build:
	# CSP lambdas
	$(foreach \
		LAMBDA,$(CSP_LAMBDA_DIRS), \
		echo "building $(LAMBDA)"; \
		env GOOS=linux GOARCH=amd64 go build -o $(LAMBDA)/main $(LAMBDA)/*.go; \
		zip -j $(LAMBDA)/lambda.zip $(LAMBDA)/main; \
	)
	# util lambdas
	$(foreach \
		LAMBDA,$(UTIL_LAMBDA_DIRS), \
		echo "building $(LAMBDA)"; \
		env GOOS=linux GOARCH=amd64 go build -o $(LAMBDA)/main $(LAMBDA)/*.go; \
		zip -j $(LAMBDA)/lambda.zip $(LAMBDA)/main; \
	)

openapi-build: clean-openapi
	cd $(SPECS_API_DIR)/${STAGE}/ && yq m csp-*.yaml > csp-api.yaml

# API_ENV="dev" STAGE="2019-05-14" make terraform
terraform-apply: test-tf-vars test-stage build
	cd api_gateway && make ${API_REGION}-${API_ENV} apply-auto-approve
	cd lambda_functions && make ${API_REGION}-${API_ENV} apply-auto-approve
	cd fargate && make ${API_REGION}-${API_ENV} apply-auto-approve

# API_ENV="dev" STAGE="2019-06-01" make terraform-test
terraform-plan: test-tf-vars openapi-build test-stage build
	cd api_gateway && make ${API_REGION}-${API_ENV} plan
	cd lambda_functions && make ${API_REGION}-${API_ENV} plan
	cd fargate && && make ${API_REGION}-${API_ENV} plan

# API_ENV=dev API_REGION=us-west-2 STAGE="2019-10-01" make all
all: openapi-build build terraform-apply

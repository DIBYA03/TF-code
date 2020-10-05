# core-platform

### Requirements

```bash
# instal golint and dep
make dev-install
```

### monorepo layout

https://github.com/golang-standards/project-layout

- api: handlers for auth, client, and partner facing apis
- cmd: all of our lambda executables go here
- docs: shared documentation
- examples: example code of how to use various packages
- partner: partner functionality - e.g. banks or vendors
- services: core functionality for platform (banking/payroll/etc.)
- shared: shared or common utility code
- specs: Specifications such as for API's or Database
- terraform: Build and deployment
- test: tests
- tools: tools
- vendor: third party vendor packages

---

Here's a really good blog post from Shippable about monorepos: http://blog.shippable.com/our-journey-to-microservices-and-a-mono-repository

1. Better developer testing: Developers can easily run the entire platform on their machine and this helps them understand all services and how they work together. This has led our developers to find more bugs locally before even sending a pull request.

2. Reduced code complexity: Senior engineers can easily enforce standardization across all services since it is easy to keep track of pull requests and changes happening across the repository.

3. Effective code reviews: Most developers now understand the end to end platform leading to more bugs being identified and fixed at the code review stage.

4. Sharing of common components: Developers have a view of what is happening across all services and can effectively carve out common components. Over a few weeks, we actually found that the code for each microservice became smaller, as a lot of common functionality was identified and shared across services.

5. Easy refactoring: Any time we want to rename something, refactoring is as simple as running a grep command. Restructuring is also easier as everything is neatly in one place and easier to understand.

### Makefile

One Makefile will handle deployment for all of the services in this repo.

```bash
make test -PROJECT_PATH=cmd/lambda-api-gateway/business
make build -PROJECT_PATH=cmd/lambda-api-gateway/business
```

### Vendor

If you’re writing a library, especially an open-source library, it’s not generally a good idea to commit your `vendor` directory. If you’re writing an applicaiton that emits binaries, there are some arguments for committing your `vendor/` directory as an application that builds a binary and some arguments against it.

- You have all the source and binaries needed to build your application in your repository. This can speed CI builds by avoiding a lot of downloading and dependency resolution when building. It also gives you a repeatable set of source to build from.
- You’re protected from upstream changes breaking your builds—if a dependency unpublishes a previous release, you’ll still have a copy.

More Info: [golang-dep](https://lightstep.com/blog/golang-dep)

# Naming and environment variables
To help keep this project in uniform, Wise Platform/ CloudOps teams have a document that is used for naming resources and labeling what environment variables are for:

https://docs.google.com/spreadsheets/d/1w7hBqaDayQA32_xsBY4f1VVOTj99-Kr9QxGZyfOA7hs

# Teraform

## Deployment Steps
1. bbva
    >Don't run every deployment. There is a subscription that is happening and we don't need this unless a change is happening in this repo
2. csp
   > This is seperate from client API, so run as needed, but not needed on every client API deployment
3. client_api
4. fargate
5. cloudfront
    >Don't run every deployment. This is setup to automatically get a WAF and the deployment will be nonsense and say that a WAF should be removed

## Sample deployment
```
# from root of repo
API_ENV=dev BASE_PATH="2019.03.01" STAGE="2019-03-01" make all
```

The BASE_PATH is the extra basepath that will be added. Stage is the name of the new stage.

## Fargate resource information

## Task definition mem/cpu
* 512 (0.5 GB), 1024 (1 GB), 2048 (2 GB) - Available cpu values: 256 (.25 vCPU)
* 1024 (1 GB), 2048 (2 GB), 3072 (3 GB), 4096 (4 GB) - Available cpu values: 512 (.5 vCPU)
* 2048 (2 GB), 3072 (3 GB), 4096 (4 GB), 5120 (5 GB), 6144 (6 GB), 7168 (7 GB), 8192 (8 GB) - Available cpu values: 1024 (1 vCPU)
* Between 4096 (4 GB) and 16384 (16 GB) in increments of 1024 (1 GB) - Available cpu values: 2048 (2 vCPU)
* Between 8192 (8 GB) and 30720 (30 GB) in increments of 1024 (1 GB) - Available cpu values: 4096 (4 vCPU)
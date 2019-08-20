functions := $(shell find lambda-fns -name \*main.go | awk -F'/' '{print $$2}')

build: ## Build golang binaries
	@for function in $(functions) ; do \
		env GOOS=linux go build -ldflags="-s -w" -o bin/$$function lambda-fns/$$function/*.go ; \
	done

clean:
	rm -rf bin/*

watch:
	gobin -m -run github.com/cortesi/modd/cmd/modd --file lambda-fns/modd.conf

local-api:
	sam local start-api --skip-pull-image --env-vars .env.json
local-api-debug: ## to debug a lambda invocation with breakpoints
	sam local start-api --skip-pull-image --env-vars .env.json -d 5986 --debugger-path .debugger --debug-args "-delveAPI=2"
local-cors-proxy: ## because stupid SAM CLI does't support CORS in local
	npx lcp --proxyUrl http://localhost:3000 --proxyPartial ''

deploy-staging: build
	serverless deploy --stage staging --verbose

test:
	go test -v ./...

## targets to create tables on particular environments
create-tables-dev:
	env APP_ENV=development go run cmd/create-tables/main.go
create-tables-staging:
	env APP_ENV=staging go run cmd/create-tables/main.go
create-tables-prod:
	env APP_ENV=production go run cmd/create-tables/main.go

## targets to fetch cities on particular environments
fetch-cities-dev:
	env APP_ENV=development go run cmd/fetch-cities/main.go
fetch-cities-staging:
	env APP_ENV=staging go run cmd/fetch-cities/main.go
fetch-cities-prod:
	env APP_ENV=production go run cmd/fetch-cities/main.go

compile-debugger: ## because the binary runs inside linux container we need to compile the debugger with linux as the targer
	GOARCH=amd64 GOOS=linux go build -o .debugger/dlv github.com/go-delve/delve/cmd/dlv
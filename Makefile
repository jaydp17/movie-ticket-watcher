functions := $(shell find lambda-fns -name \*main.go | awk -F'/' '{print $$2}')

build: ## Build golang binaries
	@for function in $(functions) ; do \
		env GOOS=linux go build -ldflags="-s -w" -o bin/$$function lambda-fns/$$function/main.go ; \
	done

clean:
	rm -rf bin/*

watch:
	gobin -m -run github.com/cortesi/modd/cmd/modd --file lambda-fns/modd.conf

local-api:
	sam local start-api --skip-pull-image --env-vars .env.json

deploy-staging: build
	serverless deploy --stage staging --verbose

test:
	go test -v ./pkg/cinemas

## targets to create tables on particular environments
create-tables-dev:
	env APP_ENV=development go run cmd/create-tables/main.go
create-tables-staging:
	env APP_ENV=staging go run cmd/create-tables/main.go
create-tables-prod:
	env APP_ENV=production go run cmd/create-tables/main.go
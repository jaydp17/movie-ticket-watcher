lambda-watch:
	gobin -m -run github.com/cortesi/modd/cmd/modd --file lambda-fns/modd.conf

lambda-serve:
	sam local start-api --skip-pull-image

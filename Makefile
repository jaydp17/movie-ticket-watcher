lambda-watch:
	gobin -m -run github.com/cortesi/modd/cmd/modd --file lambda-fns/modd.conf

lambda-serve:
	# dummy aws credentials are required else it go to Disk for finding it on every request
	env AWS_ACCESS_KEY_ID=hello AWS_SECRET_ACCESS_KEY=world sam local start-api --skip-pull-image

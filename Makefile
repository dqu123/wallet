run:
	source secrets.env && go run main.go

psql:
	psql wallet

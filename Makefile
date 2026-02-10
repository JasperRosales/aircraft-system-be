dev:
	air

build:
	@go build -o bin/app cmd/app/main.go

push:
	@git add .
	@git commit -m "$(m)"
	@git push origin main

branch:
	@git checkout -b $(b)
	@git add .
	@git commit -m "$(m)"
	@git push --set-upstream origin $(b)

ping:
	@curl http://localhost:8080/ping

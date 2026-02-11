dev:
	air

build:
	@go build -o bin/app cmd/app/main.go

push:
	@git add .
	@git commit -m "$(m)"
	@git push

branch:
	@git checkout -b $(b)
	@git add .
	@git commit -m "$(m)"
	@git push --set-upstream origin $(b)

merge:
	@git checkout main
	@git pull origin main
	@git merge $(b)
	@git push origin main

ping:
	@curl http://localhost:8080/ping

set-github:
	@git config --global user.name "$(name)"
	@git config --global user.email "$(email)"


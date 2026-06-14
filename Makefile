# FileENIAC Build

.PHONY: build test clean

build:
	cd backend && go build -o ../bin/fileeniac .

test:
	cd backend && go test ./... -v

test-short:
	cd backend && go test ./...

clean:
	rm -rf bin/
	rm -rf apps/desktop/node_modules/
	rm -rf apps/desktop/src-tauri/target/

backend:
	cd backend && go mod tidy && go build -o ../bin/fileeniac .

desktop:
	cd apps/desktop && npm install && npm run tauri dev

lint:
	cd backend && golangci-lint run ./...

cross-compile:
	cd backend && GOOS=linux GOARCH=amd64 go build -o ../bin/fileeniac-linux .
	cd backend && GOOS=darwin GOARCH=amd64 go build -o ../bin/fileeniac-darwin .
	cd backend && GOOS=windows GOARCH=amd64 go build -o ../bin/fileeniac.exe .

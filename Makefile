all: assets
	@mkdir -p build
	@go build -o build/joe cmd/joe/*.go
	@ls -la build/joe

assets:
	@rm -rf doc/templates/compiled.go
	@go-bindata -o doc/templates/compiled.go -pkg templates doc/templates/

install:
	@cp build/joe /usr/local/bin/

clean:
	@rm -rf build
all:
	@mkdir -p build
	@go build -o build/joe cmd/joe/*.go
	@ls -la build/joe

install_example_rule:
	@mkdir -p /etc/joe/queries
	@cp example.yml /etc/joe/queries/

install:
	@cp build/joe /usr/local/bin/

clean:
	@rm -rf build
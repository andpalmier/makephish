all: makephish

makephish: build
	@go build -o build/makephish cmd/makephish/*.go

install: makephish
	@cp build/makephish /usr/bin/
	@chmod a+x /usr/bin/makephish

build:
	@mkdir -p build

clean:
	@rm -rf build

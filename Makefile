.PHONY: build
build:
	go build

.PHONY: run
run:
	go run .


.PHONY: clean
clean:
	rm redisdel
.PHONY: package
package: build
	tar czvf redisdel.tar.gz redisdel
	shasum -a 256 redisdel.tar.gz
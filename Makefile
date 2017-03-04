.PHONY: build
build:
	go build cmd/ohrad.go


.PHONY: clean
clean:
	rm -f ohrad

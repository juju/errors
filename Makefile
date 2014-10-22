default: check

check:
	go test && go test -compiler gccgo

docs:
	godoc2md github.com/juju/errors > README.md


# Main Makefile for fa-export

GOBIN=	${GOPATH}/bin

SRCS=	fa-tail.go

all:	fa-tail

install:
	go install ${SRCS}

clean:
	go clean -v

fa-tail: fa-tail.go
	echo "This is replaced by github.com/keltia/flightaware-go."
	exit 1

push:
	git push --all
	git push --all origin
	git push --all backup
	git push --tags
	git push --tags origin
	git push --tags backup

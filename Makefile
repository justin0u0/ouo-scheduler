
BIN_DIR=_output/bin

# If tag not explicitly set in users default to the git sha.
TAG ?= ${shell (git describe --tags --abbrev=14 | sed "s/-g\([0-9a-f]\{14\}\)$/+\1/") 2>/dev/null || git rev-parse --verify --short HEAD}

all: local

init:
	mkdir -p ${BIN_DIR}

local: init
	go build -o=${BIN_DIR}/ouo-scheduler ./cmd/scheduler

image: init
	docker build --no-cache . -t ouo-scheduler:${TAG}

clean:
	rm -rf _output/
	rm -f *.log

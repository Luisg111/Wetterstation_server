sudo docker run --rm -v "$PWD":/app -w /app -e CGO_ENABLED=1 golang:1.23 go build -v -buildvcs=false -tags osusergo,netgo,sqlite_omit_load_extension -ldflags="-extldflags=-static"

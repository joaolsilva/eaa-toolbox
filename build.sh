cd "`dirname \"$0\"`"
rm cmd/eaa-toolbox/eaa-toolbox
go fmt ./...
source ~/go/src/gocv.io/x/gocv/env.sh
cd cmd/eaa-toolbox
CGO_CPPFLAGS_ALLOW=.* CGO_CFLAGS_ALLOW=.* CGO_LDFLAGS_ALLOW=.* go build
cd ../..

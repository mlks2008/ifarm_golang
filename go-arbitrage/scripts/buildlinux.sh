ScriptDir=$(cd $(dirname $0);pwd)
cd $ScriptDir/../

export GO111MODULE=on
export GOARCH=amd64
export GOOS=linux
export CGO_ENABLED=0
export GOROOT=/usr/local/go
go build -o ./scripts/goarbitrage_oneplat


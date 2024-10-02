# 同步执行程序
ScriptDir=$(cd $(dirname $0);pwd)

export GOARCH=amd64
export GOOS=linux
export CGO_ENABLED=0
export GOROOT=/usr/local/go
go mod tidy
go build -ldflags "-X 'main.gitHash=$(git show -s --format=%H)'" -o fil

COPYFILE_DISABLE=1 tar --no-xattrs -zcvf fil.bin fil-mainapi-restart.sh fil-mainapi-stop.sh fil-oneplat-restart.sh fil-oneplat-stop.sh fil
rm -rf fil

echo "上传执行程序..."
sshpass -p A8b4df5i scp -r -P 22 fil.bin root@47.236.130.49:/home/admin/go-tests-up
rm -rf fil.bin
echo "done"

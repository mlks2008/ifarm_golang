# 同步执行程序
ScriptDir=$(cd $(dirname $0);pwd)

export GOARCH=amd64
export GOOS=linux
export CGO_ENABLED=0
export GOROOT=/usr/local/go
go mod tidy
go build -o doge

COPYFILE_DISABLE=1 tar --no-xattrs -zcvf doge.bin restart.sh stop.sh doge
rm -rf doge

echo "上传执行程序..."
sshpass -p A8b4df5i scp -r -P 22 doge.bin root@47.236.130.49:/home/admin/go-tests
rm -rf doge.bin
echo "done"

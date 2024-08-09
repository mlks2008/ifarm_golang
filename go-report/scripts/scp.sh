ScriptDir=$(cd $(dirname $0);pwd)
cd $ScriptDir/../

export GOARCH=amd64
export GOOS=linux
export CGO_ENABLED=0
export GOROOT=/usr/local/go
go mod tidy
go build -o ./scripts/go-report

cd $ScriptDir/
COPYFILE_DISABLE=1 tar --no-xattrs -zcvf go-report.tar.gz ./go-report ./scp.sh ./restart.sh

echo "上传执行程序..."
sshpass -p A8b4df5i scp -r -P 22 go-report.tar.gz root@47.236.130.49:/home/admin/go-report/
rm -rf go-report.tar.gz
rm -rf go-report
echo "done"

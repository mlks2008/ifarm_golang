# 同步执行程序
ScriptDir=$(cd $(dirname $0);pwd)
cd $ScriptDir/../

export GOARCH=amd64
export GOOS=linux
export CGO_ENABLED=0
export GOROOT=/usr/local/go
go mod tidy
go build -o ./dca

COPYFILE_DISABLE=1 tar --no-xattrs -zcvf bin/dca.bin ./configs ./dca
rm -rf ./dca

echo "上传执行程序..."
sshpass -p A8b4df5i scp -r -P 22 bin/dca.bin root@47.236.130.49:/home/admin
echo "拉取运行数据..."
sshpass -p A8b4df5i ssh root@47.236.130.49 -t 'cd /home/admin;tar -zcvf data.tar.gz files log-*';
sshpass -p A8b4df5i scp -r -P 22 root@47.236.130.49:/home/admin/data.tar.gz /Users/iworkspace/ifarm_golang/go-arbitrage
tar -zxvf data.tar.gz
rm -rf data.tar.gz
echo "done"

#mac安装sshpass
#-----------------------------------------------------------------------------------------------------------------------------------
#安装包：https://sourceforge.net/projects/sshpass/
#开始安装
#tar -zxvf sshpass-1.06.tar.gz
#cd sshpass-版本号
#sh configure
#make && make install
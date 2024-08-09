# 同步执行程序
ScriptDir=$(cd $(dirname $0);pwd)

sshpass -p A8b4df5i scp -r -P 22 root@47.236.130.49:/home/admin/go-tests/l.log /Users/iworkspace/ifarm_golang/go-arbitrage/tests/doge2

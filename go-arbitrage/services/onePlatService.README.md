### 开启套利机器人
启动oneplat套利机器人
```
cd /Users/iworkspace/gowallets/goarbitrage/services
source ~/.bash_profile
proxy
./Test_OnePlat_Start.test >> ../../logs/goarbitrage/test_oneplat.log
```

套利日志
```
tail -f /Users/iworkspace/gowallets/logs/goarbitrage/test_oneplat.log
```

启动/暂停oneplat套利机器人
```
curl -X POST 'http://127.0.0.1:6868' --header 'Content-Type: application/json' -d '{"method":"StopOnePlat","params":true}'
```

设置oneplat套利参数
```
curl -X POST 'http://127.0.0.1:6868' --header 'Content-Type: application/json' -d '{"method":"SetOnePlat","params":{"OnePlatUsdtAmount1":11, "OnePlatUsdtAmount2":15, "OnePlatUsdtAmount3":20}}'
```

### tmux工具
#### 开启新会话
```tmux new -s oneplat```
#### 执入会话
```tmux attach -t oneplat```
#### 杀死会话
```tmux kill-session -t oneplat```
#### 查看会话
```tmux ls```
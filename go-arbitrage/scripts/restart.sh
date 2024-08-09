# 启动机器人

ps -ef | grep "./dca" | awk '{print $2}'|xargs kill

tar -zxvf dca.bin

rm -rf dca.bin

nohup ./dca -cmd=arb > log-dca-arb.log &

sleep 2

nohup ./dca -cmd=bnb > log-dca-bnb.log &

sleep 1
ps -ef | grep "./dca"
ps -ef | grep "./doge" | awk '{print $2}'|xargs kill

tar -zxvf doge.bin

rm -rf doge.bin

nohup ./doge > log-doge.log &

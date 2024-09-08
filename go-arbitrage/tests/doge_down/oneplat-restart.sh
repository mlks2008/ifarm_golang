ps -ef | grep "./doge -robot=oneplat" | awk '{print $2}'|xargs kill

tar -zxvf doge.bin

rm -rf doge.bin

nohup ./doge -robot=oneplat -maxSellOrders=6 > log.oneplat &

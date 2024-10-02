ps -ef | grep "./doge -robot=oneplat" | awk '{print $2}'|xargs kill

tar -zxvf doge.bin

rm -rf doge.bin

nohup ./doge -robot=oneplat -minAllowPrice=0.096 -maxSellOrders=8 > doge.oneplat.log &

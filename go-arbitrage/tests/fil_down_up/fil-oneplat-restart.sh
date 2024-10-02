ps -ef | grep "./fil -robot=oneplat" | awk '{print $2}'|xargs kill

tar -zxvf fil.bin

rm -rf fil.bin

nohup ./fil -robot=oneplat -minAllowPrice=3 -maxSellOrders=8 > fil.oneplat.log &

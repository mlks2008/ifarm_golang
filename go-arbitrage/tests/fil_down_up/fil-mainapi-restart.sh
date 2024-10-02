ps -ef | grep "./fil -robot=mainapi" | awk '{print $2}'|xargs kill

tar -zxvf fil.bin

rm -rf fil.bin

nohup ./fil -robot=mainapi -minAllowPrice=3 -maxSellOrders=8 > fil.mainapi.log &

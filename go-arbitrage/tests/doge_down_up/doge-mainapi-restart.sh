ps -ef | grep "./doge -robot=mainapi" | awk '{print $2}'|xargs kill

tar -zxvf doge.bin

rm -rf doge.bin

nohup ./doge -robot=mainapi -minAllowPrice=0.09 -maxSellOrders=8 > doge.mainapi.log &

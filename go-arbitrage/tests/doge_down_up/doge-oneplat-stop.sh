# 启动机器人

ps -ef | grep "./doge -robot=oneplat" | awk '{print $2}'|xargs kill
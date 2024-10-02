# 启动机器人

ps -ef | grep "./doge -robot=mainapi" | awk '{print $2}'|xargs kill
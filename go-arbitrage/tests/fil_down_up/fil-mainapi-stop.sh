# 启动机器人

ps -ef | grep "./fil -robot=mainapi" | awk '{print $2}'|xargs kill
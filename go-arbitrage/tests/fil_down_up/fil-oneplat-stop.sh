# 启动机器人

ps -ef | grep "./fil -robot=oneplat" | awk '{print $2}'|xargs kill
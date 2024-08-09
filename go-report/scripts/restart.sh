ps -ef | grep "./go-report" | awk '{print $2}'|xargs kill

tar -zxvf go-report.tar.gz

rm -rf go-report.tar.gz

nohup ./go-report > log-go-report.log &

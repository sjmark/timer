# 下载
go get github.com/sjmark/timer
# 使用
go Start()

1. 一次性任务
AddOnce(key string, sec time.Duration, fn func())
2. 永久任务
AddForever(key string, sec time.Duration, fn func())
3. 停止某个计时器
StopCron(key string)



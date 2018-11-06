# 下载
go get github.com/sjmark/timer
# 使用
1. 开启计时器

go Start()

2. 一次性任务

AddOnce(key string, sec time.Duration, fn func())

3. 永久任务

AddForever(key string, sec time.Duration, fn func())

4. 停止某个计时器

StopCron(key string)



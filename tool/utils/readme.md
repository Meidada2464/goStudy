## utils
一些公用的组件、工具之类的
```
.
├── command.go
├── command_test.go
# float相关的函数
├── float.go
├── float_test.go
├── gzip.go
├── gzip_test.go
# 获取本机IP（利用net包）
├── ips.go
├── ips_test.go
# 耗时选择器，按照访问耗时选择节点，围绕Latency（根据耗时）和LatencyShuffle（随机）的一系列方法
├── latency.go
├── latency_test.go
# 计算MD5，支持string, []byte, filepath
├── md5.go
├── md5_test.go
├── md5_test.txt
├── readme.md
# 时间相关函数：任务循环执行，检测时间是否在区间内，检测函数执行耗时
├── timer.go
├── timer_test.go
├── unique.go
└── unique_test.go
```
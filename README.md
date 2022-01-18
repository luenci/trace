# trace

本库基于极客时间课程[《Tony Bai · Go语言第一课》]([https://time.geekbang.org/column/intro/100093501?code=cQ4ugiP4uzDdDVD1T-HXXlTv9Fdl-SpdsPnSfxf0%2FuU%3D]) 教程开发

使用方法：
```shell
// 执行编译
make build

// 根据编译后的命令 对目标文件进行追踪插入代码

// 不对目标文件(examples/demo.go)进行插入,将输出打印到控制台
./bin/core examples/demo.go

// 对目标文件(examples/demo.go)进行追踪插入代码
./bin/core -w=true examples/demo.go

```


## TODO List
- [ ] 使用 traceId 替代 goroutineID
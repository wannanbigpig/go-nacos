# go-nacos
[![Go](https://github.com/wannanbigpig/go-nacos/actions/workflows/go.yml/badge.svg)](https://github.com/wannanbigpig/go-nacos/actions/workflows/go.yml)
### 构建 liunx 执行文件
- mac 环境下构建liunx可执行文件

`CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/go-nacos .`
- 构建本机可执行文件

`go build -o ./bin/go-nacos .`

#### 使用说明：
> 可执行文件要赋予执行权限
```shell
    ./go-nacos --help  // 查看参数
    
    ./go-nacos -dataId=test -group=dev -path=/home/www/xxxx/ -filename=.env 
    
    // OR 后台运行
    
    nohup /home/go-nacos/go-nacos-liunx -dataId=test -group=dev -path=/home/www/xxxx/ -filename=.env  > runoob.log 2>&1 &
```

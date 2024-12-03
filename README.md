# TxNet

TxNet 是一个用于与天下网络系统进行交互的 Golang SDK。它提供了一些 API 接口，便于用户管理、转账操作、查询记录等功能。

## 安装

1. 安装依赖包

   ```bash
   go get github.com/wfunc/txnet

## 配置
1. 配置文件 example.properties 用于设置以下选项：
```
txsrv/verbose: 设置日志详细级别，1 为开启详细日志。
txsrv/proxy_addr: 设置代理地址（可选）。
txsrv/timeout: 设置请求超时（单位：秒）。
txsrv/api_host: 设置 API 主机地址。
txsrv/website: 设置网站标识。
txsrv/uppername: 设置上层账号。
txapi: 配置各个 API 的详细信息。
```

2. 配置示例
```
txsrv/verbose=1
txsrv/proxy_addr=http://proxy.example.com
txsrv/timeout=5
txsrv/api_host=http://api.example.com
txsrv/website=mywebsite.com
txsrv/uppername=admin
txapi/CreateMember={"keyA":10, "keyB":"some_key", "keyC":20}
```

## 使用
1. 初始化
首先，调用`Bootstrap`函数进行初始化，传入配置文件路径作为参数。
```
package main

import (
    "github.com/wfunc/txnet"
)

func main() {
    txnet.Bootstrap("example.properties")
}
```

## API调用示例
* 创建用户
```
resp, err := txnet.CreateMember("new_user")
if err != nil {
    log.Fatal(err)
}
fmt.Println(resp)
```
* 登陆
```
redirectURL := txnet.Login("existing_user")
fmt.Println("Redirect URL:", redirectURL)
```



## test-data-service-go
分布式数据获取服务:为自动化测试脚本，压测脚本等测试脚本提供测试数据的服务

[v1.0.0 下载地址](https://github.com/dperfly/test-data-service-go/releases/tag/1.0.0) 

 ## 接口列表

|  接口  | 说明 |
|------ |----- |
|/{fileName} | 获取数据接口|

***

## 接口详情
* 接口地址：/fileName

* 请求方式：Get

* 请求示例：http://{ip}:{Port}/demo.txt?num=2&type=random

* 接口备注：通过fileName的后坠来自动区别接口返回内容。

* 请求参数说明：

| 名称 | 类型 | 必填 |说明|
|----- |------| ---- |----|
|fileName |string|true|文件名|
|<font color=red>num | string |false|csv和txt 生效 返回数据的数量,csv从第二行开始返回,默认返回1行|
|<font color=red>type | string |false|可选项目type=random,默认根据文本内容顺序返回数据，随机数据未去重|

* 返回参数说明：

| 名称 | 类型 |说明|
|----- |------|----|
| msg | string|错误信息
|result | object|具体数据|

* JSON返回示例：
```
    json data :  http://ip:port/json_demo.json
         {
            "msg": "",
            "result": {
                "age": 18,
                "name": "zhang san"
            }
        }

    csv data:  http://ip:port/csv_demo.csv
        {
            "msg": "",
            "result": [
                [
                    "user3",
                    "123456"
                ]
            ]
        }


    txt data:  http://ip:port/phone_demo.txt
        {
            "msg": "",
            "result": [
                "18666660001"
            ]
        }

```
---

## 注意事项
* 测试数据默认读取启动文件同级的file文件,将测试数据放入到此文件中;也可以通过启动服务时指定file文件读取时的绝对路径
 ```
 ./jmeter-data-service-go --file=/tmp/file
 ```
* 修改测试数据后需要手动重启服务;
* IP默认localhost,端口号默认8080,可以通过启动服务指定host 和 port;
 ```
 ./jmeter-data-service-go --host=192.168.31.111 --port=9999
 ```
* 尽量避免在测试中对此服务进行压测，此服务的目的是数据分发，也就是压测前的必要数据获取;
* 关于数据提取可以参考example目录下的demo文件;

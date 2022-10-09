# Net
- 自动登录华科校园网
## 原理
访问`123.123.123.123`会接收到一串下面的信息
```
<script>top.self.location.href='http://192.168.50.3:8080/eportal/index.jsp?wlanuserip=xxx&wlanacname=xxx&ssid=
&nasip=xxx&snmpagentip=&mac=xxx&t=wireless-v2&url=xxx&apmac=&nasid=xxx&vid=xxx&port=xxx&nasportid=xxx'</script>
```
其中的参数就是网络及设备信息，利用这些信息，加上个人账户信息（从`config.json`读取），即可登录`172.18.18.60:8080`
## 编译方法
直接下载可执行文件和`login.json`文件，并放在同一目录下  
[Windows](./main.exe)  
[Linux](./main)  
[配置文件](./login.json)
```shell
# Windows //-H windowsgui 隐藏窗口// -s -w 缩减大小
go build -ldflags="-H windowsgui -s -w" main.go

# Linux
go build -ldflags="-s -w" main.go
```
## 提供账户密码
编辑`login.json`文件

## 使用
- 直接运行（确保目录下有`login.json`文件）  
- 也可设为开机启动

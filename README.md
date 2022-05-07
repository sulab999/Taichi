# TAIJISUITE
## 簡介
基於go語言實現的高交互滲透測試框架，已實現如下功能：  
1、"RDP","JAVADEBUG","REDIS", "FTP", "SNMP", "POSTGRESQL", "SSH", "MONGO", "SMB", "MSSQL", "MYSQL", "ELASTICSEARCH"服務的弱口令掃描；  
2、敏感路徑掃描（基於字典）  
3、子域名掃描（基於字典）    
4、增加poc模塊（內置兼容了kunpeng的50+poc）  
![TAIJI](https://github.com/sulab999/Taichi/raw/main/demo.png "demo")
## 編譯運行
1、安裝第三方庫（命令：go get xxx）  
2、go run main.go  
3、編譯
go build  

## 使用
1.加載模塊  
load <模塊> <協議>  
e.g:  
load portscan  
load burt ftp  
2.設置參數  
set ip/file  xxx  
3.展示參數  
show  
4.運行  
go  

端口掃描
1.load portscan
2.set ip 127.0.0.1
或者 set file xxx.txt
c段 set ip 192.168.1.1-255
3.go

敏感路径扫描（需要本地有urldic.txt）  
1.load urlscan  
2.set ip/domain xxx  
3.go

子域名扫描（需要本地有subdic.txt）  
1.load subscan  
2.set domain http://xxx.cn  
3.go  

poc功能  
本地创建plugins文件夹，用于存放yml文件  
1.poc  
2.init（首次使用或新增poc時）  
3.show  
4.use xxx  
5.set ip/url xxx  
6.go  

## 更新：  
v0.1  
1、已實現端口掃描和爆破模塊
## 来元世界交流一下啊
![TAIJI](https://github.com/sulab999/Taichi/blob/main/nworld.jpg)
![TAIJI](https://github.com/sulab999/Taichi/blob/main/webchat.png)

下載地址：https://github.com/sulab999/Taichi/releases/tag/v0.2.2

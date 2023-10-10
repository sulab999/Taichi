# TAICHISUITE
## 簡介
基於go語言實現的高交互滲透測試框架，已實現如下功能：  
1、"RDP","JAVADEBUG","REDIS", "FTP", "SNMP", "POSTGRESQL", "SSH", "MONGO", "SMB", "MSSQL", "MYSQL", "ELASTICSEARCH"服務的弱口令掃描；  
2、敏感路徑掃描（基於字典）  
3、子域名掃描（基於字典）    
4、增加poc模塊（已更新700+poc）  
5、url存活檢測  
6、端口掃描&服務識別  
![TAIJI](https://github.com/sulab999/Taichi/raw/main/demo.png "demo")
## 編譯運行
1、安裝第三方庫（命令：go get xxx）  
2、go run main.go  
3、編譯
go build  

## 基本使用
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

## 端口掃描  
1.load portscan  
2.set ip 127.0.0.1  
可設置文件 set file xxx.txt  
c段 set ip 192.168.1.1-255  
3.go  

## 敏感路徑掃描（需要本地有urldic.txt）  
1.load urlscan  
2.set ip/domain xxx  
3.go

## 子域名掃描（需要本地有subdic.txt）  
1.load subscan  
2.set domain http://xxx.cn  
3.go  

## url存活掃描  
![TAICHI](https://github.com/sulab999/Taichi/blob/main/test/livescan.png)  
1.load urlscan live   
2.set file url.txt  
3.go 

## poc功能  
本地創建taichi-pocs文件夾，用於存放yml文件  
1.poc  
2.init（首次使用或新增poc時）  
3.show  
4.use xxx 或set xxx（poc）  
5.set ip/url xxx  
6.go  
掃描結束後，生成的報告在reports文件夾中  

## 更新：
後期更新見realse  
v0.1  
1、已實現端口掃描和爆破模塊
## 來元世界交流一下啊
![TAICHI](https://github.com/sulab999/Taichi/blob/main/nworld.jpg)
![TAICHI](https://github.com/sulab999/Taichi/blob/main/webchat.png)

# 免責聲明
該程序及其相關技術僅用於安全自查檢測。

由於傳播、利用此程序所提供的信息而造成的任何直接或者間接的後果及損失，均由使用者本人負責，作者不為此承擔任何責任。

本人擁有對此程序的修改和解釋權。未經網絡安全部門及相關部門允許，不得善自使用本程序進行任何攻擊活動，不得以任何方式將其用於商業目的。

下載地址：https://github.com/sulab999/Taichi/releases/tag/v0.5

# TAIJISUITE
## 簡介
基於go語言實現的仿MSF滲透測試框架，已實現"RDP","JAVADEBUG","REDIS", "FTP", "SNMP", "POSTGRESQL", "SSH", "MONGO", "SMB", "MSSQL", "MYSQL", "ELASTICSEARCH"服務的弱口令掃描  
![TAIJI](https://github.com/sulab999/Taichi/blob/main/demo.png "demo")
## 編譯運行
1、安裝第三方庫（命令：go get xxx）  
2、go run main.go

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
## 更新：  
v0.1  
1、已實現端口掃描和爆破模塊

下載地址：https://github.com/sulab999/Taichi/releases/tag/v0.1

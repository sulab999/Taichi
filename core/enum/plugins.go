package enum

type ScanFunc func(ip string, port string, username string, password string) (result bool, err error)

var (
	ScanFuncMap map[string]ScanFunc
)

func init() {
	ScanFuncMap = make(map[string]ScanFunc)
	ScanFuncMap["FTP"] = ScanFtp
	ScanFuncMap["SSH"] = ScanSsh
	ScanFuncMap["SMB"] = ScanSmb
	ScanFuncMap["MSSQL"] = ScanMssql
	ScanFuncMap["MYSQL"] = ScanMysql
	ScanFuncMap["POSTGRESQL"] = ScanPostgres
	ScanFuncMap["REDIS"] = ScanRedis

	ScanFuncMap["MONGO"] = ScanMongodb
	ScanFuncMap["JAVADEBUG"] = JavaDebug

	ScanFuncMap["RDP"] = ScanRdp

	ScanFuncMap["SNMP"] = ScanSnmp

}

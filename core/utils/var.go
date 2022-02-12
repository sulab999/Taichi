package utils

import (
	"strings"
	"sulab/core/model"
	"sync"
)

var (
	Mutex sync.Mutex

	PortNames = map[int]string{
		21:    "FTP",
		22:    "SSH",
		161:   "SNMP",
		445:   "SMB",
		1433:  "MSSQL",
		3306:  "MYSQL",
		5432:  "POSTGRESQL",
		6379:  "REDIS",
		9200:  "ELASTICSEARCH",
		27017: "MONGODB",
	}

	SupportProtocols map[string]bool

	BruteResult map[string]model.Service
)

func init() {

	BruteResult = make(map[string]model.Service)

	SupportProtocols = make(map[string]bool)
	for _, proto := range PortNames {
		SupportProtocols[strings.ToUpper(proto)] = true
	}

}

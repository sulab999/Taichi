package cli

import (
	"strings"

	"github.com/c-bata/go-prompt"
)

func excludeOptions(args []string) []string {
	ret := make([]string, 0, len(args))
	for i := range args {
		if !strings.HasPrefix(args[i], "-") {
			ret = append(ret, args[i])
		}
	}
	return ret
}

func Completer2(d prompt.Document) []prompt.Suggest {
	if d.TextBeforeCursor() == "" {
		return []prompt.Suggest{}
	}
	args := strings.Split(d.TextBeforeCursor(), " ")

	for i := range args {
		if args[i] == "|" {
			return []prompt.Suggest{}
		}
	}
	return argumentsCompleter2(d, excludeOptions(args))

}

func Completer3(d prompt.Document) []prompt.Suggest {
	if d.TextBeforeCursor() == "" {
		return []prompt.Suggest{}
	}
	args := strings.Split(d.TextBeforeCursor(), " ")

	for i := range args {
		if args[i] == "|" {
			return []prompt.Suggest{}
		}
	}
	return argumentsCompleter3(d, excludeOptions(args))

}

func Completer(d prompt.Document) []prompt.Suggest {
	if d.TextBeforeCursor() == "" {
		return []prompt.Suggest{}
	}
	args := strings.Split(d.TextBeforeCursor(), " ")

	for i := range args {
		if args[i] == "|" {
			return []prompt.Suggest{}
		}
	}
	return argumentsCompleter2(d, excludeOptions(args))
}

var commands = []prompt.Suggest{

	{Text: "load", Description: "加載模塊."},
	{Text: "portscan", Description: "端口掃描."},
	{Text: "show", Description: "顯示設置信息"},
	{Text: "set", Description: "設置參數"},
	{Text: "help", Description: "幫助"},
	{Text: "exit", Description: "退出"},
}

func argumentsCompleter3(d prompt.Document, args []string) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "load", Description: "加載模塊"},
		{Text: "show", Description: "顯示設置信息"},
		{Text: "set", Description: "設置參數"},
		{Text: "go", Description: "執行"},
		{Text: "help", Description: "幫助"},
	}
	if len(args) <= 1 {
		return prompt.FilterHasPrefix(s, args[0], true)
	}

	first := args[0]
	switch first {

	case "set":
		if len(args) == 2 {
			subcommands := setParse2()
			return prompt.FilterHasPrefix(subcommands, args[1], true)
		}
	case "load":
		if len(args) == 2 {
			subcommands := loadParse()
			return prompt.FilterHasPrefix(subcommands, args[1], true)
		}

		if len(args) == 3 {
			switch args[1] {
			case "burst":
				subcommands := loadParse2()
				return prompt.FilterHasPrefix(subcommands, args[2], true)
			}
		}
	default:
		return []prompt.Suggest{}
	}

	return []prompt.Suggest{}
}

func argumentsCompleter2(d prompt.Document, args []string) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "load", Description: "加載模塊"},
		{Text: "show", Description: "顯示設置信息"},
		{Text: "set", Description: "設置參數"},
		{Text: "go", Description: "執行"},
		{Text: "help", Description: "幫助"},
	}
	if len(args) <= 1 {
		return prompt.FilterHasPrefix(s, args[0], true)
	}

	first := args[0]
	switch first {

	case "set":
		if len(args) == 2 {
			subcommands := setParse()
			return prompt.FilterHasPrefix(subcommands, args[1], true)
		}
	case "load":
		if len(args) == 2 {
			subcommands := loadParse()
			return prompt.FilterHasPrefix(subcommands, args[1], true)
		}

		if len(args) == 3 {
			switch args[1] {
			case "burst":
				subcommands := loadParse2()
				return prompt.FilterHasPrefix(subcommands, args[2], true)
			}
		}
	default:
		return []prompt.Suggest{}
	}

	return []prompt.Suggest{}
}

func loadParse() []prompt.Suggest {
	subcommands := []prompt.Suggest{
		{Text: "portscan", Description: "端口掃描"},
		{Text: "urlscan", Description: "路径扫描"},
		{Text: "subscan", Description: "子域名扫描"},
		{Text: "burst", Description: "爆破"},
		{Text: "new", Description: "..."},
	}
	return subcommands
}

func loadParse2() []prompt.Suggest {
	subcommands := []prompt.Suggest{
		{Text: "ftp", Description: "ftp爆破"},
		{Text: "ssh", Description: "ssh爆破"},
		{Text: "mysql", Description: "mysql爆破"},
		{Text: "3389", Description: "3389爆破"},
		{Text: "new", Description: "..."},
	}
	return subcommands
}

func setParse() []prompt.Suggest {
	subcommands := []prompt.Suggest{
		{Text: "ip", Description: "設置ip"},
		{Text: "file", Description: "設置文件"},
	}
	return subcommands
}

func setParse2() []prompt.Suggest {
	subcommands := []prompt.Suggest{
		{Text: "ip", Description: "設置ip"},
		{Text: "file", Description: "設置文件"},
		{Text: "port", Description: "設置端口（不設置就使用默認端口）"},
	}
	return subcommands
}

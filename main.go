package main

import (
	"fmt"
	"strings"
	"sulab/core/cli"
	"sulab/core/utils"

	"github.com/c-bata/go-prompt"
	"github.com/fatih/color"
)

var (
	author  string = "sulab"
	version string = "0.1"
)

func showBanner() {
	name := fmt.Sprintf("Taiji (v.%s)", version)
	banner := `
	_________   ____   ____________  ________________
	/_  __/ _ | /  _/_ / /  _/ __/ / / /  _/_  __/ __/
	 / / / __ |_/ // // // /_\ \/ /_/ // /  / / / _/  
	/_/ /_/ |_/___/\___/___/___/\____/___/ /_/ /___/  												  																																
   `

	all_lines := strings.Split(banner, "\n")
	w := len(all_lines[1])

	fmt.Println(banner)
	color.Green(fmt.Sprintf("%[1]*s", -w, fmt.Sprintf("%[1]*s", (w+len(name))/2, name)))
	color.Blue(fmt.Sprintf("%[1]*s", -w, fmt.Sprintf("%[1]*s", (w+len(author))/2, author)))
	fmt.Println()
}

func initCore() {

	showBanner()

	utils.InitConfig()
}

func main() {

	initCore()

	p := prompt.New(
		cli.Executor,
		cli.Completer,
		prompt.OptionTitle("Taiji: Penetration Test Framework"),
		prompt.OptionPrefix("[Taiji] > "),
		prompt.OptionLivePrefix(cli.ChangeLivePrefix),
		prompt.OptionInputTextColor(prompt.White),
	)
	p.Run()
}

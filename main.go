package main

import (
	"fmt"
	"os"

	"github.com/splitscriptjs/cli/config"
	"github.com/splitscriptjs/cli/utils"
)

func main() {
	addCmd("dev", make(args), func(args []string) {
		fmt.Println(utils.Block.Render("Splitscript"))

		conf, err := config.Read()
		if os.IsNotExist(err) {
			fmt.Println(utils.Error.Render(err.Error()))
			fmt.Println(utils.Warning.Render("Try running `splitscript init`"))
			os.Exit(1)
		} else if err != nil {
			fmt.Println(utils.Error.Render("Couldn't read config"))
			fmt.Println(err.Error())
			os.Exit(1)
		}
		watchDir(conf, "./")
	})
	parseCmd()
}

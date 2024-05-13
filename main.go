package main

import (
	"fmt"
	"os"
	"splitscript/config"
	"splitscript/utils"
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

package main

import (
	"fmt"
	"os"
	"splitscript/config"
)

func main() {
	addCmd("dev", make(args), func(args []string) {
		fmt.Println(block.Render("Splitscript"))

		conf, err := config.Read()
		if os.IsNotExist(err) {
			fmt.Println(errMessage.Render(err.Error()))
			fmt.Println(warning.Render("Try running `splitscript init`"))
			os.Exit(1)
		} else if err != nil {
			fmt.Println(errMessage.Render("Couldn't read config"))
			fmt.Println(err.Error())
			os.Exit(1)
		}
		watchDir(conf, "./")
	})
	parseCmd()
}

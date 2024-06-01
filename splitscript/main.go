package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	"github.com/splitscriptjs/cli/cmd"
	"github.com/splitscriptjs/cli/config"
	"github.com/splitscriptjs/cli/utils"
)

func main() {
	cmd.Add("dev", "Run and watch an app", make(cmd.Args), func(args []string) {
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
	cmd.Add("help", "View this message", make(cmd.Args), func(_ []string) {
		columns := []table.Column{
			{Title: "Name", Width: 8},
			{Title: "Description", Width: 32},
		}
		rows := cmd.List()

		t := table.New(table.WithColumns(columns), table.WithRows(rows), table.WithHeight(len(rows)-1))
		s := table.DefaultStyles()
		s.Header = s.Header.
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240")).
			BorderBottom(true).
			Bold(false)

		t.SetStyles(s)
		var box = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240"))
		fmt.Println(box.Render(t.View()))
	})
	cmd.Parse()
}

package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/splitscriptjs/cli/utils"
)

func tsBoilerplate(packageName, eventName string) string {
	return fmt.Sprintf("import { Events } from '%s';\nexport default async function (event: Events.%s) {\n\n}", packageName, eventName)
}
func esmBoilerplate(packageName, eventName string) string {
	return fmt.Sprintf("/** @typedef {import('%s').Events.%s} Event */\n/** @param {Event} event */\n\nexport default async function (event) {\n\n}", packageName, eventName)
}
func cjsBoilerplate(packageName, eventName string) string {
	return fmt.Sprintf("/** @typedef {import('%s').Events.%s} Event */\n/** @param {Event} event */\n\nmodule.exports = async function (event) {\n\n}", packageName, eventName)
}

const (
	ts = iota
	esm
	cjs
)

func Boilerplate(t int, packageName, eventName string) string {
	switch t {
	case ts:
		return tsBoilerplate(packageName, eventName)
	case esm:
		return esmBoilerplate(packageName, eventName)
	case cjs:
		return cjsBoilerplate(packageName, eventName)
	}
	return ""
}

type Package struct {
	Type string `json:"type"`
}

func GetProjectType() int {
	bytes, err := os.ReadFile("package.json")
	if os.IsNotExist(err) && err != nil {
		fmt.Println(utils.Warning.Render("Failed to read package.json, defaulting type to CommonJS"))
		return cjs
	}
	var pkg Package
	err = json.Unmarshal(bytes, &pkg)
	if err != nil {
		fmt.Println(utils.Warning.Render("Failed to read package.json, defaulting type to CommonJS"))
		return cjs
	}
	if pkg.Type == "module" {
		return esm
	}
	return cjs
}

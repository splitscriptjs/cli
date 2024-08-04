package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/splitscriptjs/cli/config"
	"github.com/splitscriptjs/cli/debounce"
	"github.com/splitscriptjs/cli/run"
	"github.com/splitscriptjs/cli/utils"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/farmergreg/rfsnotify"
	"gopkg.in/fsnotify.v1"
)

func watchDir(conf config.Config, dir string) error {
	t := GetProjectType()
	err := buildAll(conf)
	if err != nil {
		return err
	}
	run.Run(conf)
	watcher, err := rfsnotify.NewWatcher()
	if err != nil {
		return err
	}
	debounceBuild := debounce.New(250 * time.Millisecond)
	debounceRun := debounce.New(250 * time.Millisecond)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok || (filepath.Ext(event.Name) != ".ts" && filepath.Ext(event.Name) != ".js") {
					break
				}
				isTs := filepath.Ext(event.Name) == ".ts"
				if event.Op == fsnotify.Create {
					var bp []byte

					pkg := getEventPackage(event.Name)
					if pkg == "" {
						break
					}
					eventName := getEventName(event.Name, pkg)
					typeName := formatToTypeName(eventName)
					validEvent, pkgName := isValidEvent(pkg, eventName)
					fmt.Println(validEvent, pkgName)
					if !validEvent {
						break
					}
					if isTs {
						bp = []byte(Boilerplate(ts, pkgName, typeName))
					} else {
						bp = []byte(Boilerplate(t, pkgName, typeName))
					}
					err := os.WriteFile(event.Name, bp, 0666)
					if err != nil {
						fmt.Println(utils.Error.Render("Failed to boilerplate " + event.Name))
						fmt.Println(err.Error())
						break
					}
					debounceBuild(func() { build(conf, event.Name) })
					debounceRun(func() { run.Run(conf) })
				} else if event.Op == fsnotify.Write {
					debounceBuild(func() {
						build(conf, event.Name)
					})
					debounceRun(func() { run.Run(conf) })
				} else if event.Op == fsnotify.Remove || event.Op == fsnotify.Rename {
					fileToRemove, err := utils.GenerateDevFileName(conf, event.Name)
					if err != nil {
						fmt.Println(err.Error())
						break
					}
					err = os.Remove(fileToRemove)
					if err != nil {
						fmt.Println(err.Error())
						break
					}
					debounceRun(func() { run.Run(conf) })
					fmt.Println(utils.Info.Render("Updated"))
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					break
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.AddRecursive(dir)
	if err != nil {
		return err
	}
	if len(conf.Ignore) > 0 {
		ignore = conf.Ignore
	}
	ignore = append(ignore, conf.Dev, conf.Out)
	for _, folder := range ignore {
		watcher.RemoveRecursive(path.Join(dir, folder))
	}

	<-make(chan struct{})

	return nil
}

var ignore = []string{"node_modules", ".git"}

func build(conf config.Config, path string) {
	outFile, err := utils.GenerateDevFileName(conf, path)
	if err != nil {
		fmt.Println(utils.Error.Render("Failed to build " + path))
		fmt.Println(err.Error())
		return
	}
	_ = api.Build(api.BuildOptions{
		EntryPoints: []string{path},
		Outfile:     outFile,
		Bundle:      false,
		Write:       true,
		LogLevel:    api.LogLevelInfo,
	})
}
func buildAll(conf config.Config) error {
	fmt.Println(utils.Info.Render("Rebuilding"))

	clearDevDir(conf)
	files := walk()
	includesMain := false
	for i := range files {
		if conf.Main == files[i] {
			includesMain = true
			break
		}
	}
	if !includesMain {
		fmt.Println(utils.Error.Render("Main file `" + conf.Main + "` not found"))
		os.Exit(1)
	}
	for i := range files {
		build(conf, files[i])
	}
	return nil
}

var result = []string{}

func walk() []string {
	result = []string{}
	err := filepath.WalkDir("./", visit)
	if err != nil {
		fmt.Println(err.Error())
	}
	return result
}
func visit(path string, di fs.DirEntry, err error) error {
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	if di.IsDir() {
		folders := strings.Split(path, string(filepath.Separator))
		for i := range folders {
			for j := range ignore {
				if folders[i] == ignore[j] {
					return filepath.SkipDir
				}
			}
		}
	} else {
		ext := filepath.Ext(path)
		if ext == ".js" || ext == ".ts" {
			result = append(result, path)
		}
	}
	return nil
}

func clearDevDir(conf config.Config) {
	err := os.RemoveAll(conf.Dev)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func getEventPackage(file string) string {
	segments := strings.Split(file, string(filepath.Separator))
	if len(segments) < 2 {
		return ""
	}
	if segments[0] != "functions" {
		return ""
	}
	return segments[1]
}
func getEventName(file, eventPackage string) string {
	rel, err := filepath.Rel(path.Join("functions", eventPackage), file)
	if err != nil {
		return ""
	}
	segments := filepath.SplitList(path.Dir(rel))
	return strings.Join(segments, "/")
}
func formatToTypeName(eventName string) string {
	segments := strings.Split(eventName, "/")
	formatted := make([]string, len(segments))
	for i := range segments {
		if len(segments[i]) == 0 {
			continue
		}
		formatted[i] = strings.ToUpper(string(segments[i][0])) + segments[i][1:]
	}
	return strings.Join(formatted, "")
}

type SSJson struct {
	Packages map[string]Package `json:"packages"`
}
type Package struct {
	ValidEvents []string `json:"validEvents"`
	PackageName string   `json:"packageName"`
}

func isValidEvent(eventPackage, eventName string) (bool, string) {
	ss, err := os.ReadFile("./ss.json")
	if os.IsNotExist(err) {
		return false, ""
	}
	var ssJson SSJson
	err = json.Unmarshal(ss, &ssJson)
	if err != nil {
		utils.Warning.Render("ss.json is incorrectly formatted")
	}
	if ssJson.Packages[eventPackage].PackageName == "" || len(ssJson.Packages[eventPackage].ValidEvents) == 0 {
		return false, ""
	}
	return utils.Includes(ssJson.Packages[eventPackage].ValidEvents, eventName), ssJson.Packages[eventPackage].PackageName
}

package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"splitscript/config"
	"splitscript/debounce"

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
	watcher, err := rfsnotify.NewWatcher()
	if err != nil {
		return err
	}
	debounce := debounce.New(250 * time.Millisecond)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok || (filepath.Ext(event.Name) != ".ts" && filepath.Ext(event.Name) != ".js") {
					return
				}
				isTs := filepath.Ext(event.Name) == ".ts"
				if event.Op == fsnotify.Create {
					var bp []byte
					if isTs {
						bp = []byte(Boilerplate(ts, "@splitscript.js/discord", "MessageCreate"))
					} else {
						bp = []byte(Boilerplate(t, "@splitscript.js/discord", "MessageCreate"))
					}
					err := os.WriteFile(event.Name, bp, 0666)
					if err != nil {
						fmt.Println(errMessage.Render("Failed to boilerplate " + event.Name))
						fmt.Println(err.Error())
						return
					}
					debounce(func() { build(conf, event.Name) })
				} else if event.Op == fsnotify.Write {
					debounce(func() { build(conf, event.Name) })
				} else if event.Op == fsnotify.Remove || event.Op == fsnotify.Rename {
					fileToRemove, err := generateDevFileName(conf, event.Name)
					if err != nil {
						fmt.Println(err.Error())
						return
					}
					err = os.Remove(fileToRemove)
					if err != nil {
						fmt.Println(err.Error())
						return
					}
					fmt.Println(info.Render("Updated"))
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
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
	outFile, err := generateDevFileName(conf, path)
	if err != nil {
		fmt.Println(errMessage.Render("Failed to build " + path))
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
	fmt.Println(info.Render("Rebuilding"))

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
		fmt.Println(errMessage.Render("Main file `" + conf.Main + "` not found"))
		os.Exit(1)
	}
	for i := range files {
		build(conf, files[i])
	}
	return nil
}

func generateDevFileName(conf config.Config, path string) (string, error) {
	rel, err := filepath.Rel("./", path)
	if err != nil {
		return "", err
	}
	ext := filepath.Ext(rel)
	if ext == ".ts" {
		rel = strings.TrimSuffix(rel, ext) + ".js"
	}
	return filepath.Join(conf.Dev, rel), nil
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

package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/alecthomas/kong"
)

const (
	ConfigName = "tp.json"
)

type Config struct {
	License struct {
		List       bool     `name:"list" short:"l" help:"list all available license templates."`
		Files      []string `arg:"" name:"files" help:"license template names." optional:""`
		OutputPath string   `name:"output" short:"o"`
	} `cmd:"" help:"Fetch software licenses."`
	Gitignore struct {
		List       bool     `name:"list" short:"l" help:"list all available gitignore templates."`
		Files      []string `arg:"" name:"files" help:"gitignore template names." optional:""`
		OutputPath string   `name:"output" short:"o"`
	} `cmd:"" help:"fetch gitignore templates."`
}

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// This should rarely happen
		errExit(1, err.Error())
	}

	configDir, err := os.UserConfigDir()
	if err != nil {
		errExit(1, err.Error())
	}
	_, _ = homeDir, configDir

	var config Config
	ctx := kong.Parse(&config)
	switch ctx.Command() {
	case "license":
		fmt.Println("list files")
	case "gitignore":
		fmt.Println("list files")
	case "license <files>":
		var f *os.File
		if config.License.OutputPath != "" {
			f, err = os.Create(config.License.OutputPath)
			if err != nil {
				errExit(1, err.Error())
			}
		}
		defer f.Close()

		name := config.License.Files[0]

		path, err := url.JoinPath("https://raw.githubusercontent.com/spdx/license-list-data/main/text", name+".txt")
		if err != nil {
			errExit(1, err.Error())
		}

		resp, err := http.Get(path)
		if err != nil {
			errExit(1, err.Error())
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			errExit(1, err.Error())
		}

		if f != nil {
			if _, err := f.Write(body); err != nil {
				errExit(2, err.Error())
			}
		} else {
			fmt.Println(string(body))
		}
	case "gitignore <files>":
		fmt.Println("hello")
	default:
		panic(ctx.Command())
	}
}

func errExit(code int, format string, args ...any) {
	fmt.Fprintf(os.Stderr, "tp error: "+format+"\n", args...)
	os.Exit(code)
}

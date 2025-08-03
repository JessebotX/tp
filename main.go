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

type Context struct {
	Offline bool
}

type LicenseCommand struct {
	List       bool     `name:"list" short:"l" help:"list all available license templates."`
	Stdout     bool     `name:"stdout" help:"Print contents to stdout instead of writing to a file path (i.e. output to terminal)"`
	Names      []string `arg:"" name:"names" help:"license template identifiers/names." optional:""`
	OutputPath string   `name:"output" short:"o" default:"LICENSE"`
}

func (l *LicenseCommand) Run(ctx *Context) error {
	if l.List {
		fmt.Println("TODO: list all available licenses.")

		return nil
	}

	if len(l.Names) == 0 {
		return fmt.Errorf("missing arguments")
	}

	var f *os.File
	var err error
	if !l.Stdout {
		f, err = os.Create(l.OutputPath)
		if err != nil {
			return err
		}
	}
	defer f.Close()

	for _, name := range l.Names {
		path, err := url.JoinPath("https://raw.githubusercontent.com/spdx/license-list-data/main/text", name+".txt")
		if err != nil {
			return err
		}

		resp, err := http.Get(path)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		if !l.Stdout {
			if _, err := f.Write(body); err != nil {
				return err
			}
		} else {
			fmt.Println(string(body))
		}
	}

	return nil
}

type GitignoreCommand struct {
	List       bool     `name:"list" short:"l" help:"list all available gitignore templates."`
	Stdout     bool     `name:"stdout" help:"Print contents to stdout instead of writing to a file path (i.e. output to terminal)"`
	Names      []string `arg:"" name:"names" help:"gitignore template identifier/names." optional:""`
	OutputPath string   `name:"output" short:"o"`
}

func (g *GitignoreCommand) Run(ctx *Context) error {
	fmt.Println("gitignore")
	return nil
}

type Config struct {
	License   LicenseCommand   `cmd:"" help:"Fetch software licenses."`
	Gitignore GitignoreCommand `cmd:"" help:"fetch gitignore templates."`
}

func main() {
	var config Config
	ctx := kong.Parse(&config)

	if err := ctx.Run(&Context{}); err != nil {
		ctx.FatalIfErrorf(err)
	}
}

// func errExit(code int, format string, args ...any) {
// 	fmt.Fprintf(os.Stderr, "tp error: "+format+"\n", args...)
// 	os.Exit(code)
// }

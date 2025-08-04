package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/alecthomas/kong"
)

const (
	LicenseListURL       = "https://api.github.com/repos/spdx/license-list-data/contents/text"
	LicenseDownloadURL   = "https://raw.githubusercontent.com/spdx/license-list-data/main/text"
	GitignoreListURL     = "https://api.github.com/repos/github/gitignore/contents"
	GitignoreDownloadURL = "https://raw.githubusercontent.com/github/gitignore/main"
)

type Context struct {
	Offline bool
}

type FetchItem struct {
	Name string
	Type string
}

type LicenseCommand struct {
	List       bool     `name:"list" short:"l" help:"List all available license templates."`
	Stdout     bool     `name:"stdout" help:"Print contents to stdout instead of writing to a file path (i.e. output to terminal)"`
	Names      []string `arg:"" name:"names" help:"License template identifiers/names." optional:""`
	OutputPath string   `name:"output" short:"o" default:"LICENSE"`
}

func (l *LicenseCommand) Run(ctx *Context) error {
	if l.List {
		body, err := fetchBytes(LicenseListURL)
		if err != nil {
			return err
		}

		var respItems []FetchItem
		if err := json.Unmarshal(body, &respItems); err != nil {
			return err
		}

		for _, v := range respItems {
			fmt.Println(strings.TrimSuffix(v.Name, ".txt"))
		}

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
		name = strings.TrimSuffix(name, ".txt")

		path, err := url.JoinPath(LicenseDownloadURL, name+".txt")
		if err != nil {
			return err
		}
		body, err := fetchBytes(path)
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
	List       bool     `name:"list" short:"l" help:"List all available gitignore templates."`
	Stdout     bool     `name:"stdout" help:"Print contents to stdout instead of writing to a file path (i.e. output to terminal)"`
	Names      []string `arg:"" name:"names" help:"Gitignore template identifier/names." optional:""`
	OutputPath string   `name:"output" short:"o" default:".gitignore"`
}

func (g *GitignoreCommand) Run(ctx *Context) error {
	if g.List {
		body, err := fetchBytes(GitignoreListURL)
		if err != nil {
			return err
		}

		var respItems []FetchItem
		if err := json.Unmarshal(body, &respItems); err != nil {
			return err
		}

		for _, v := range respItems {
			if v.Type == "dir" && v.Name[0] != '.' {
				fmt.Println(v.Name + "/")
				continue
			}

			if v.Type != "dir" && (!strings.HasSuffix(v.Name, ".gitignore") || v.Name == ".gitignore") {
				continue
			}

			if v.Name[0] != '.' {
				fmt.Println(strings.TrimSuffix(v.Name, ".gitignore"))
			}
		}

		return nil
	}

	if len(g.Names) == 0 {
		return fmt.Errorf("missing arguments")
	}

	var f *os.File
	var err error
	if !g.Stdout {
		f, err = os.Create(g.OutputPath)
		if err != nil {
			return err
		}
	}
	defer f.Close()

	for _, name := range g.Names {
		name = strings.TrimSuffix(name, ".gitignore")

		path, err := url.JoinPath(GitignoreDownloadURL, name+".gitignore")
		if err != nil {
			return err
		}

		body, err := fetchBytes(path)
		if err != nil {
			return err
		}

		if !g.Stdout {
			if _, err := f.Write(body); err != nil {
				return err
			}
		} else {
			fmt.Println(string(body))
		}
	}

	return nil
}

func fetchBytes(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

type Config struct {
	License   LicenseCommand   `cmd:"" help:"Fetch software licenses."`
	Gitignore GitignoreCommand `cmd:"" help:"Fetch gitignore templates."`
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

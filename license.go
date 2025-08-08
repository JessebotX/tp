package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"
)

const (
	LicenseListURL     = "https://api.github.com/repos/spdx/license-list-data/contents/text"
	LicenseDownloadURL = "https://raw.githubusercontent.com/spdx/license-list-data/main/text"
)

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

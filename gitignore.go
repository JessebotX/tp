package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

const (
	GitignoreListURL     = "https://api.github.com/repos/github/gitignore/contents"
	GitignoreDownloadURL = "https://raw.githubusercontent.com/github/gitignore/main"
)

type GitignoreCommand struct {
	List       bool     `name:"list" short:"l" help:"List all available gitignore templates."`
	Stdout     bool     `name:"stdout" help:"Print contents to stdout instead of writing to a file path (i.e. output to terminal)"`
	Names      []string `arg:"" name:"names" help:"Gitignore template identifier/names." optional:""`
	OutputPath string   `name:"output" short:"o" default:".gitignore"`
}

func (g *GitignoreCommand) Run(ctx *Context) error {
	if g.List {
		if err := printGitignoreList(GitignoreListURL); err != nil {
			return err
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

		var body []byte

		for body == nil || strings.HasSuffix(strings.TrimSpace(string(body)), ".gitignore") {
			linkedName := strings.TrimSpace(string(body))
			parent := filepath.Dir(name)

			if strings.HasSuffix(linkedName, ".gitignore") {
				path, err = url.JoinPath(GitignoreDownloadURL, parent, linkedName)
				if err != nil {
					return err
				}

				body, err = fetchBytes(path)
				if err != nil {
					return err
				}
			} else {
				body, err = fetchBytes(path)
				if err != nil {
					return err
				}
			}
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

func printGitignoreList(path string) error {
	body, err := fetchBytes(path)
	if err != nil {
		return err
	}

	var respItems []FetchItem
	if err := json.Unmarshal(body, &respItems); err != nil {
		return err
	}

	for _, v := range respItems {
		if v.Type != "dir" && (!strings.HasSuffix(v.Name, ".gitignore") || v.Name == ".gitignore") {
			continue
		}

		if v.Type == "dir" && v.Name[0] != '.' {
			newURL, err := url.JoinPath(path, v.Name)
			if err != nil {
				return err
			}

			if err := printGitignoreList(newURL); err != nil {
				return err
			}
		} else {
			fmt.Println(strings.TrimSuffix(v.Path, ".gitignore"))
		}
	}

	return nil
}

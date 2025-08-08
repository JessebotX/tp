package main

import (
	"fmt"

	"github.com/alecthomas/kong"
)

const (
	Version = "1.0.0"
)

type VersionCommand struct{}

func (v *VersionCommand) Run(ctx *Context) error {
	fmt.Printf("tp version %s\n", Version)
	return nil
}

// TODO...
type Context struct{}

type Config struct {
	License   LicenseCommand   `cmd:"" help:"Fetch software licenses."`
	Gitignore GitignoreCommand `cmd:"" help:"Fetch gitignore templates."`
	Version   VersionCommand   `cmd:"version" help:"Print program version."`
}

func main() {
	var config Config
	ctx := kong.Parse(&config)

	if err := ctx.Run(&Context{}); err != nil {
		ctx.FatalIfErrorf(err)
	}
}

package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/silvabyte/audeticlinkinstaller/link"
)

var cli = link.LinkInstaller{}

func main() {
	ctx := kong.Parse(&cli,
		kong.Name("audeticlink"),
		kong.Description("Audetic Link device management tool"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
			Summary: true,
		}),
	)
	err := ctx.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

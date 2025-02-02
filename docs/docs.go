//go:build docs

package main

import (
	"fmt"
	"os"
	"path"
	"strings"

	log "github.com/rs/zerolog"
	"github.com/spf13/cobra/doc"

	"github.com/maxgio92/yap/cmd"
	"github.com/maxgio92/yap/internal/commands/options"
)

const (
	cmdline      = "yap"
	docsDir      = "docs"
	fileTemplate = `---
title: %s
---	

`
)

var (
	filePrepender = func(filename string) string {
		title := strings.TrimPrefix(
			strings.TrimSuffix(strings.ReplaceAll(filename, "_", " "), ".md"),
			fmt.Sprintf("%s/", docsDir),
		)
		return fmt.Sprintf(fileTemplate, title)
	}
	linkHandler = func(filename string) string {
		if filename == cmdline+".md" {
			return "_index.md"
		}
		return filename
	}
)

func main() {
	if err := doc.GenMarkdownTreeCustom(
		cmd.NewRootCmd(options.NewCommonOptions(options.WithLogger(log.New(os.Stderr).Level(log.InfoLevel)))),
		docsDir,
		filePrepender,
		linkHandler,
	); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err := os.Rename(path.Join(docsDir, cmdline+".md"), path.Join(docsDir, "_index.md"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

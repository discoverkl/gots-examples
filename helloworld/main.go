package main

import (
	"bufio"
	"embed"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/discoverkl/gots/ui"
)

//go:embed index.html
var www embed.FS

func add(a, b int) int {
	return a + b
}

func main() {
	app := ui.New(
		ui.Mode(promptRunMod()),
		ui.Root(www),
		ui.OnlineAddr(":8000"),
	)
	app.BindFunc("add", add)
	app.Run()
}

func promptRunMod() string {
	for {
		fmt.Print(promptText)
		ch, _, err := bufio.NewReader(os.Stdin).ReadRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Println()
				os.Exit(0)
			}
			log.Fatal(err)
		}
		switch ch {
		case '1', 10:
			return "app"
		case '2':
			return "page"
		case '3':
			return "online"
		case 'q':
		default:
			os.Exit(0)
		}
	}
}

const promptText = `
*** Commands ***

1: LocalApp - start a local web server, open its' serving url within a native app
2: LocalPage - start a local web server, open its' serving url with your default web browser
3: Online   - run a online web server

Please enter (1-3)? `

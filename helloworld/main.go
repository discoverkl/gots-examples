package main

import (
	"bufio"
	"embed"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/discoverkl/gots/ui"
)

//go:embed fe/dist
var root embed.FS

func add(a, b int) int {
	return a + b
}

func main() {
	www, _ := fs.Sub(root, "fe/dist")
	app := ui.New(
		ui.Mode(promptRunMod()),
		ui.Root(http.FS(www)),
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
		case '1':
			return "page"
		case '2':
			return "app"
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

1: LocalPage - start a local web server, open its' serving url with your default web browser
2: LocalApp - start a local web server, open its' serving url within a native app (which is a chrome process)
3: Online   - run a online web server

Please enter (1-3)? `

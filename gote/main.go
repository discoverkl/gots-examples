package main

import (
	"embed"
	"flag"
	"io/fs"
	"log"
	"os"

	"github.com/discoverkl/gots/ui"
)

//go:embed fe/dist
var root embed.FS

//go:embed notes
var goNotes embed.FS

var port int

func main() {
	// parse arguments: 'path', 'port'
	flag.IntVar(&port, "p", -1, "binding port")
	flag.Parse()

	path := ""
	if flag.NArg() >= 1 {
		path = flag.Arg(0)
		log.Printf("notes path: %s", path)
	}

	// prepare data
	var notes fs.FS
	if path == "" {
		notes, _ = fs.Sub(goNotes, "notes")
	} else {
		notes = os.DirFS(path)
	}

	// options
	www, _ := fs.Sub(root, "fe/dist")
	ops := []ui.Option{ui.Root(www)}
	if port > 0 {
		ops = append(ops, ui.Mode("online"), ui.OnlinePort(port))
	} else if port == 0 {
		ops = append(ops, ui.Mode("page"))
	} else {
		ops = append(ops, ui.Mode("app"))
	}

	// create and run
	app := ui.New(ops...)
	app.BindObject(NewAPI(notes))
	app.Run()
}

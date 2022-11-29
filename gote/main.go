package main

import (
	"embed"
	"flag"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/discoverkl/gots/ui"
)

//go:embed fe/dist
var root embed.FS

//go:embed notes
var goNotes embed.FS

var port int
var keyPath string

func main() {
	// parse arguments: 'path', 'port'
	flag.IntVar(&port, "p", -1, "binding port")
	flag.StringVar(&keyPath, "key", "", "login key")
	flag.Parse()

	path := ""
	if flag.NArg() >= 1 {
		path = flag.Arg(0)
		log.Printf("notes path: %s", path)
	}

	// prepare data
	var notes fs.FS
	var abspath string
	if path == "" {
		notes, _ = fs.Sub(goNotes, "notes")
	} else {
		abspath, _ = filepath.Abs(path)
		notes = os.DirFS(abspath)
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

	api := NewAPI(notes, abspath)
	// login key
	if keyPath != "" {
		data, err := os.ReadFile(keyPath)
		if err != nil {
			log.Fatalf("read key failed: %v", err)
		}
		key := strings.TrimSpace(string(data))
		api.SetLoginKey(key)
	}

	// create and run
	app := ui.New(ops...)
	app.BindObject(api)
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

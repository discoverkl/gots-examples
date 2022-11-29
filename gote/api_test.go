package main

import (
	"embed"
	"os"
	"testing"
)

var api API

//go:embed notes
var www embed.FS

func init() {
	fsys := os.DirFS("notes")
	// fsys, _ := fs.Sub(www, "notes")
	api = *NewAPI(fsys, "notes")
}

func BenchmarkListNotes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		api.ListNotes()
	}
}

func BenchmarkListNotes2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		api.listNotes2()
	}
}

func BenchmarkListNotes3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		api.listNotes3()
	}
}

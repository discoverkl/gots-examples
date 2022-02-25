package main

import (
	"io/fs"
	"log"
	"path/filepath"
	"strings"
)

type NoteItem struct {
	Title   string   `json:"title"`
	Path    string   `json:"path"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

const MaxHeaderSize = 100

type API struct {
	fs fs.FS
}

func NewAPI(fs fs.FS) *API {
	return &API{fs: fs}
}

func (a *API) GetNote(path string) (*NoteItem, error) {
	data, err := fs.ReadFile(a.fs, path)
	if err != nil {
		return nil, err
	}
	note := parseNote(data)
	note.Path = path
	return note, nil
}

func (a *API) ListNotes() ([]*NoteItem, error) {
	names, err := fs.Glob(a.fs, "*.md")
	if err != nil {
		return nil, err
	}

	items := make([]*NoteItem, len(names))
	for i := 0; i < len(names); i++ {
		name := names[i]
		f, err := a.fs.Open(name)
		if err != nil {
			log.Printf("open file %s: %v", name, err)
			return nil, err
		}
		data := make([]byte, MaxHeaderSize)
		_, err = f.Read(data)
		if err != nil {
			f.Close()
			log.Printf("parse file %s: %v", name, err)
			return nil, err
		}
		f.Close()
		note, _ := parseNoteHeader(data)
		note.Path = name

		// default title
		if note.Title == "" {
			note.Title = name
			ext := filepath.Ext(name)
			if ext != "" {
				index := strings.LastIndex(name, ext)
				note.Title = name[0:index]
			}
		}
		items[i] = note
	}
	return items, nil
}

func (a *API) listNotes2() ([]*NoteItem, error) {
	names, err := fs.Glob(a.fs, "*.md")
	if err != nil {
		return nil, err
	}

	items := make([]*NoteItem, len(names))
	for i := 0; i < len(names); i++ {
		name := names[i]
		data, err := fs.ReadFile(a.fs, name)
		if err != nil {
			return nil, err
		}
		note := parseNote(data)
		note.Path = name

		// default title
		if note.Title == "" {
			note.Title = name
			ext := filepath.Ext(name)
			if ext != "" {
				index := strings.LastIndex(name, ext)
				note.Title = name[0:index]
			}
		}
		items[i] = note
	}
	return items, nil
}

func (a *API) listNotes3() ([]*NoteItem, error) {
	names, err := fs.Glob(a.fs, "*.md")
	if err != nil {
		return nil, err
	}

	items := make([]*NoteItem, len(names))
	for i := 0; i < len(names); i++ {
		items[i] = &NoteItem{Path: names[i], Title: ""}
	}
	return items, nil
}

func parseNoteHeader(data []byte) (*NoteItem, int) {
	header, offset := "", 0
	for i := 0; i < len(data); i++ {
		if data[i] == '\n' {
			if i-4 >= 0 &&
				data[i-1] == '\n' &&
				data[i-2] == '-' &&
				data[i-3] == '-' &&
				data[i-4] == '-' {
				header, offset = string(data[:i+1]), i+1
			}
		}
	}

	title := ""
	var tags []string
	if header != "" {
		sp := strings.Split(header, "\n")
		for _, line := range sp {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "-") {
				continue
			}
			kv := strings.SplitN(line, ":", 2)
			if len(kv) != 2 {
				continue
			}
			key, value := strings.TrimSpace(kv[0]), strings.TrimSpace(kv[1])
			switch key {
			case "title":
				title = value
			case "tags":
				tagList := strings.TrimSpace(value[1 : len(value)-1])
				if tagList == "" {
					tags = nil
				} else {
					tags = strings.Split(tagList, ", ")
				}
			}
		}
	}
	return &NoteItem{Title: title, Tags: tags}, offset
}

func parseNote(data []byte) *NoteItem {
	note, offset := parseNoteHeader(data)
	note.Content = string(data[offset:])
	return note
}

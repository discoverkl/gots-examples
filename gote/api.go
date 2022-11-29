package main

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type NoteItem struct {
	Title   string   `json:"title"`
	Path    string   `json:"path"` // file name
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

const MaxHeaderSize = 100

var ErrNeedLogin = errors.New("need login")

type API struct {
	fs       fs.FS
	notesDir string
	logined  bool
	key      string
}

func NewAPI(fs fs.FS, notesDir string) *API {
	ret := &API{fs: fs, notesDir: notesDir, logined: true}
	return ret
}

func (a *API) SetLoginKey(key string) {
	if key == "" {
		panic("invalid login key")
	}
	a.logined = false
	a.key = key
}

func (a *API) readonly() bool {
	return a.notesDir == ""
}

// TODO: read from config file
func (a *API) Login(key string) (bool, error) {
	if a.key == "" {
		return true, nil
	}
	ok := key == a.key
	a.logined = ok
	return ok, nil
}

func (a *API) GetNote(path string) (*NoteItem, error) {
	if !a.logined {
		return nil, ErrNeedLogin
	}
	data, err := fs.ReadFile(a.fs, path)
	if err != nil {
		return nil, err
	}
	note := parseNote(data)
	note.Path = path
	return note, nil
}

func (a *API) AddNote() (*NoteItem, error) {
	if !a.logined {
		return nil, ErrNeedLogin
	}
	names, err := fs.Glob(a.fs, "*.md")
	if err != nil {
		return nil, err
	}

	var path string
	for i := 1; true; i++ {
		path = noteName(i)
		hit := false
		for _, name := range names {
			if name == path {
				hit = true
				break
			}
		}
		if !hit {
			break
		}
	}

	ext := filepath.Ext(path)
	title := path[0 : len(path)-len(ext)]
	note := &NoteItem{
		Title:   title,
		Path:    path,
		Content: fmt.Sprintf("# %s\n", title),
	}
	note, err = a.SaveNote(note)
	if err != nil {
		return nil, err
	}
	return note, nil
}

func noteName(index int) string {
	if index <= 0 {
		panic("invalid argument")
	}
	name := "Untitled"
	if index == 1 {
		return fmt.Sprintf("%s.md", name)
	}
	return fmt.Sprintf("%s (%d).md", name, index)
}

func (a *API) SaveNote(item *NoteItem) (*NoteItem, error) {
	if !a.logined {
		return nil, ErrNeedLogin
	}
	log.Printf("save path %s, content length: %d", item.Path, len(item.Content))
	if a.readonly() {
		return nil, fmt.Errorf("readonly file system")
	}
	var oldHeader string
	data, err := fs.ReadFile(a.fs, item.Path)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return nil, err
		}
	}
	if err == nil {
		_, offset := parseNoteHeader(data)
		oldHeader = string(data[0:offset])
	}
	newFileContent := oldHeader + item.Content
	path := filepath.Join(a.notesDir, item.Path)
	err = os.WriteFile(path, []byte(newFileContent), 0644)
	if err != nil {
		return nil, err
	}

	// refresh title
	err = a.refreshTitle(item)
	if err != nil {
		log.Printf("refresh title failed: %s", err)
	}
	return item, nil
}

// refresh title by content
func (a *API) refreshTitle(item *NoteItem) error {
	// title
	lines := strings.Split(item.Content, "\n")
	var firstLine string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			firstLine = line
			break
		}
	}
	var i int
	for i = 0; i < len(firstLine); i++ {
		if firstLine[i] != '#' {
			break
		}
	}

	title := firstLine
	if i < len(firstLine) {
		title = strings.TrimSpace(firstLine[i:])
	}
	if title == "" {
		return fmt.Errorf("got empty title")
	}

	// update title
	if title != item.Title {
		log.Printf("refresh title: %s -> %s", item.Title, title)
		item.Title = title
	}

	// update path
	filename := fmt.Sprintf("%s.md", title)
	if filename != item.Path {
		// move file if we can
		sfs, ok := a.fs.(fs.StatFS)
		if !ok {
			return fmt.Errorf("filesystem do not support Stat operation")
		}
		_, err := sfs.Stat(filename)
		if err == nil {
			return fmt.Errorf("file already exists: %s", filename)
		}
		if err != nil && errors.Is(err, fs.ErrNotExist) {
			log.Printf("rename note: %s -> %s", item.Path, filename)
			oldpath := filepath.Join(a.notesDir, item.Path)
			newpath := filepath.Join(a.notesDir, filename)
			err = os.Rename(oldpath, newpath)
			if err != nil {
				return fmt.Errorf("rename note failed: %w", err)
			}
			item.Path = filename

			// overwrite header
			err = os.WriteFile(newpath, []byte(item.Content), 0644)
			if err != nil {
				return fmt.Errorf("override headr failed: %w", err)
			}
		} else {
			return err
		}
	}

	return nil
}

func (a *API) ListNotes() ([]*NoteItem, error) {
	if !a.logined {
		return nil, ErrNeedLogin
	}
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

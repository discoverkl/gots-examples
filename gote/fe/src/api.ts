import { getapi, Base } from "ts2go";
import { NoteItem } from "./types";

const dev = false || true;

interface API extends Base {
  listNotes(): NoteItem[];
  getNote(path: string): NoteItem;
  saveNote(item: NoteItem): NoteItem;
  login(key: string): boolean;
  addNote(): NoteItem;
}

export let api = mockAPI();

try {
  api = getapi() as API;
} catch (ex) {
  const msg = "API is not ready!";
  if (dev) console.log(msg);
  else throw msg;
}

function mockAPI(): API {
  return {
    listNotes() {
      return [
        {
          path: "clipboard.md",
          title: "clipboard",
          tags: [] as string[],
        },
        {
          path: "Go Keywords.md",
          title: "Go Keywords",
          tags: ["go", "hi", "3", "2"],
        },
        {
          path: "midway.md",
          title: "z-midway",
          tags: ["tool"],
        },
      ];
    },
    getNote(path) {
      switch (path) {
        case "Go Keywords.md":
          return {
            content:
              `# Go Keywords

A minimal runnable go code that contain all key words:

` +
              "```go" +
              `
package main

import _ "os"

func main() {
	goto a
a:
	const t = true
	defer main()
	go main()
	type X struct{}
	var (
		c chan X
		_ map[X]interface{}
	)
	for range c {
		break
		continue
	}
	if t {
	} else {
	}
	switch {
	case t:
		fallthrough
	default:
	}
	select {}
	return
}
` +
              "```" +
              `

Ref to: https://go.dev/ref/spec#Keywords

` +
              "```" +
              `
Keywords
The following keywords are reserved and may not be used as identifiers.

break        default      func         interface    select
case         defer        go           map          struct
chan         else         goto         package      switch
const        fallthrough  if           range        type
continue     for          import       return       var
` +
              "```" +
              `
`,
          };
        case "clipboard.md":
          return {
            content: "this is a clipboard",
          };
        case "Untitled.md":
          return {
            content: "# Untitled",
          };
        default:
          return {
            content: "",
          };
      }
    },
    saveNote(item: NoteItem) {
      console.log("save note:", item.title);
      let title = item.content.split("\n")[0].trim();
      while (title.length > 0 && title[0] == "#") {
        title = title.substring(1);
      }
      title = title.trim();

      return {
        path: title + ".md",
        title: title,
        tags: item.tags,
        content: item.content,
      };
    },
    login(key: string) {
      return true;
    },
    addNote() {
      return {
        path: "Untitled.md",
        title: "Untitled",
        tags: [] as string[],
      };
    },
  } as API;
}

if (dev) (window as any).api = api;

import { getapi, Base } from "ts2go"
import { NoteItem } from "./types"

const dev = false || true

interface API extends Base {
  listNotes(): NoteItem[]
  getNote(path: string): NoteItem
}

export let api = mockAPI()

try {
 api = getapi() as API 
} catch(ex) {
  const msg = "API is not ready!"
  if (dev) console.log(msg)
  else throw msg
}

function mockAPI(): API {
  return {
    listNotes() {
      return [
        {
          path: 'clipboard.md',
          title: 'clipboard',
          tags: [] as string[],
        },
        {
          path: 'Go Keywords.md',
          title: 'Go Keywords',
          tags: ['go', 'hi', '3', '2'],
        },
        {
          path: 'midway.md',
          title: 'midway',
          tags: ['tool'],
        },
        {
          path: 'clipboard.md',
          title: 'clipboard',
          tags: [] as string[],
        },
        {
          path: 'Go Keywords.md',
          title: 'Go Keywords',
          tags: ['go1'],
        },
        {
          path: 'midway.md',
          title: 'midway',
          tags: ['tool1'],
        },
        {
          path: 'clipboard.md',
          title: 'clipboard',
          tags: [] as string[],
        },
        {
          path: 'Go Keywords.md',
          title: 'Go Keywords',
          tags: ['go2'],
        },
        {
          path: 'midway.md',
          title: 'midway',
          tags: ['tool2'],
        },
      ]
    },
    getNote(path) {
      switch (path) {
      case "Go Keywords.md":
        return {
          content: `# Go Keywords

A minimal runnable go code that contain all key words:

`+ '```go' + `
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
`+ '```' + `

Ref to: https://go.dev/ref/spec#Keywords

`+ '```' + `
Keywords
The following keywords are reserved and may not be used as identifiers.

break        default      func         interface    select
case         defer        go           map          struct
chan         else         goto         package      switch
const        fallthrough  if           range        type
continue     for          import       return       var
`+ '```' + `
`}
      case 'clipboard.md':
        return {
          content: 'this is a clipboard'
        }
      default:
        return {
          content: ''
        }
      }
    },
  } as API
}

if (dev) (window as any).api = api
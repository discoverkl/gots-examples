<script setup lang="ts">
import { marked } from "marked"
import cookies from "browser-cookies"
import { computed, onMounted, onUnmounted, ref, watch } from 'vue';
import NoteList from './components/NoteList.vue'
import { NoteItem, NoteMode, State, TagInfo, Tag } from './types'
import { api } from './api'
import { mountEditor, setText } from './editor'

const inputElement = ref(null)
const contentElement = ref(null as HTMLElement | null)
const searchElement = ref(null as HTMLElement | null)

const menuOpen = ref(false)
const fullScreen = ref(false)
const noteItems = ref([] as NoteItem[])
const content = ref('')
const noteMode = ref(NoteMode.ReadOnly)
const searchText = ref('')
const selectedTag = ref(Tag.All as string | Tag)

let editor: any;

const htmlContent = computed(() => {
  if (noteMode.value != NoteMode.Html) return ''
  return marked.parse(content.value)
})

const filterNoteItems = computed(() => {
  // filter by selected tag group
  let tagNoteItems: NoteItem[] = []
  const tag = selectedTag.value
  switch (tag) {
    case Tag.All:
      tagNoteItems = noteItems.value
      break
    case Tag.Tagged:
      for (let item of noteItems.value) {
        if (item.tags) {
          tagNoteItems.push(item)
        }
      }
      break
    case Tag.Untagged:
      for (let item of noteItems.value) {
        if (!item.tags) {
          tagNoteItems.push(item)
        }
      }
      break
    default:
      if (typeof(tag) !== 'string') {
        console.log("invalid tag value: ", tag)
        tagNoteItems = noteItems.value
        break
      }
      for (let item of noteItems.value) {
        for (let one of item.tags??[]) {
          if (tag === one) {
            tagNoteItems.push(item)
            break
          }
        }
      }
      break
  }

  // filter by search text
  if (searchText.value.length <= 1) {
    return tagNoteItems
  }
  const exactItems: NoteItem[] = []
  const items: NoteItem[] = []
  const key = searchText.value
  for (let item of tagNoteItems) {
    for (let tag of item.tags??[]) {
      if (tag === key) {
        exactItems.push(item)
        break
      }
      if (tag.indexOf(key) !== -1) {
        items.push(item)
        break
      }
    }
  }
  return exactItems.concat(items)
})

const tagItems = computed((): TagInfo => {
  let total = 0, tagged = 0, untagged = 0, tags = [] as [string, number][]

  const tag2count = new Map<string, number>()
  for (let item of noteItems.value) {
    total++
    if (item.tags == null || item.tags.length == 0) {
      untagged++
    } else  {
      tagged++
      for (let tag of item.tags) {
        tag2count.set(tag, (tag2count.get(tag) || 0) + 1)
      }
    }
  }

  // sort
  const items = [] as [string, number][]
  tag2count.forEach((value, key) => {
    items.push([key, value])
  })
  items.sort((a, b) => {
    return a[0] == b[0] ? 0 : (a[0] < b[0] ? -1 : 1)
  })

  return {
    total: total,
    tagged: tagged,
    untagged: untagged,
    tags: items,
  }
})

const closeMenuOnSelect = () => {
  return window.innerWidth <= 600
}

const hljs = (window as any).hljs
marked.setOptions({
  highlight: function(code, lang) {
    if (lang === '') return code
    const ret = hljs.highlight(code, {language: lang}).value
    return ret
  }
})

//
// init
//
onMounted(async () => {
  const stateMenuOpen = cookies.get(State.MenuOpen) || JSON.stringify(true)
  menuOpen.value = JSON.parse(stateMenuOpen)
  fullScreen.value = true
  noteItems.value = await api.listNotes()
  noteMode.value = NoteMode.Html
  // searchText.value = cookies.get(State.SearchText) || ''

  editor = await mountEditor(inputElement.value)
  handlePageClick(true)
  document.addEventListener("keyup", onKeyUp)
})

onUnmounted(() => {
  handlePageClick(false)
  document.removeEventListener("keyup", onKeyUp)
})

// TODO: use editor component
watch([noteMode, content], () => {
  if (editor == null) {
    return
  }
  if (noteMode.value == NoteMode.ReadWrite) {
    setText(editor, content.value, ".md")
  } else {
    setText(editor, '', ".md")
  }
})

watch([menuOpen, noteMode], ()=>{
  const s = document.body.style
  if (menuOpen.value || noteMode.value == NoteMode.ReadWrite) {
    // body css default to open and none-readwrite
    s.removeProperty("height")
    s.removeProperty("overflow")
  } else {
    s.height = "initial"
    s.overflow = "auto"
  }
})

watch(menuOpen, ()=> {
  cookies.set(State.MenuOpen, JSON.stringify(menuOpen.value), {expires: 365})
})

async function selectNote(note: NoteItem) {
  // cookies.set(State.SearchText, searchText.value, {expires: 365})
  if (closeMenuOnSelect()) menuOpen.value = false
  window.document.title = note.title
  const noteContent = (await api.getNote(note.path)).content
  if (content.value != noteContent) {
    if (contentElement.value) {
      contentElement.value.scrollTo(0, 0)
    }
    content.value = noteContent
  }
}

function selectTag(tag: string | Tag) {
  selectedTag.value = tag
  if (tag != Tag.All && tag != Tag.Tagged) {
    searchText.value = ''
  }
}

function toggleEdit() {
  if (noteMode.value == NoteMode.ReadWrite) {
    noteMode.value = NoteMode.Html
  } else {
    noteMode.value = NoteMode.ReadWrite
  }
}

function noteContentClick(event: MouseEvent) {
  const width = window.innerWidth
  const pos = event.clientX
  let page = 0
  if (pos <= width / 3) {
    page = -1
  } else if (pos >= width * 2 / 3) {
    page = 1
  }
  if (page !== 0) {
    event.preventDefault()
    scrollPage(page)
  }
}

function scrollPage(offset: number) {
  const pageSize = window.innerHeight - 40 - 14
  window.scrollTo(window.scrollX, window.scrollY + offset * pageSize)
}

function usePageClick(onclick: (event: MouseEvent) => void): (on: boolean) => void {
  let drag = 0
  const mousedown = () => drag = 0
  const mousemove = () => drag++
  const mouseup = (event) => {
    const click = (drag < 10)
    if (click && !menuOpen.value && noteMode.value != NoteMode.ReadWrite) {
      onclick(event)
    }
  }

  const handlePageClick = (on: boolean) => {
    if (on) {
      document.addEventListener('mousedown', mousedown)
      document.addEventListener('mousemove', mousemove)
      document.addEventListener('mouseup', mouseup)
    } else {
      document.removeEventListener('mousedown', mousedown)
      document.removeEventListener('mousemove', mousemove)
      document.removeEventListener('mouseup', mouseup)
    }
  }
  return handlePageClick
}
const handlePageClick = usePageClick(noteContentClick)

function onKeyUp(event: KeyboardEvent) {
  if (event.key == "Escape") {
    searchText.value = ''
    if (searchElement.value != null) searchElement.value.focus()
    event.preventDefault()
  }
}
</script>

<template>
<div class="root" :class="{close: !menuOpen, fullscreen: fullScreen}">
  <span class="material-icons menu-icon" @click="menuOpen = !menuOpen">
    {{ menuOpen ? "menu_open" : "menu" }}</span>
  <div class="nav border">
    <div class="nav-head"></div>
    <div class="content">
      <div class="item-1" :class="{active: selectedTag == Tag.All}" @click="selectTag(Tag.All)">
        <span class="material-icons icon">description</span>
        <span class="text">All Notes</span>
        <span class="note-count">{{ tagItems.total }}</span>
      </div>
      <div class="item-1" :class="{active: selectedTag == Tag.Tagged}" @click="selectTag(Tag.Tagged)">
        <span class="material-icons icon">local_offer</span>
        <span class="text">Tags</span>
        <span class="note-count">{{ tagItems.tagged }}</span>
      </div>
      <ul class="demo">
        <li class="item-2"><span class="text">typescript</span><span class="note-count">2</span></li>
        <li class="item-2"><span class="text">vue</span><span class="note-count">1</span></li>
        <li class="item-2"><span class="text">go</span><span class="note-count">1</span></li>
      </ul>
      <ul class="">
        <li class="item-2" v-for="(kv, index) in tagItems.tags" :key="index"
          :class="{active: selectedTag === kv[0]}" @click="selectTag(kv[0])"
          ><span class="text">{{ kv[0] }}</span>
          <span class="note-count">{{ kv[1] }}</span></li>
      </ul>
      <div class="item-1" :class="{active: selectedTag == Tag.Untagged}" @click="selectTag(Tag.Untagged)"><span class="material-icons icon">edit_off</span>
        <span class="text">Untagged</span>
        <span class="note-count">{{ tagItems.untagged }}</span>
      </div>
    </div>
  </div>
  <div class="note-list right-border">
    <div class="head">
      <input class="search" type="text" ref="searchElement" placeholder="Search by tag..." v-model="searchText" spellcheck="false" />
      <span class="material-icons icon search-icon">search</span>
    </div>
    <div class="content">
      <NoteList :items="filterNoteItems" @select="selectNote" />
      <ul class="demo">
        <li>clipboard</li>
        <li class="active">Go Keywords</li>
        <li>midway</li>
        <li>gote</li>
        <li>Web Tech</li>
      </ul>
    </div>
  </div>
  <div class="note">
    <div class="head">
      <span class="material-icons icon"
        :class="{active: noteMode == NoteMode.ReadWrite}"
        @click="toggleEdit"
        >edit</span>
    </div>
    <div class="content" ref="contentElement">
      <pre class="readonly" v-if="noteMode == NoteMode.ReadOnly">{{ content }}</pre>
      <div class="readwrite" v-show="noteMode == NoteMode.ReadWrite" ref="inputElement"></div>
      <div class="html" v-if="noteMode == NoteMode.Html" v-html="htmlContent"></div>
      <pre class="demo">
Ref to: https://go.dev/ref/spec#Keywords

```
Keywords
The following keywords are reserved and may not be used as identifiers.

break        default      func         interface    select
case         defer        go           map          struct
chan         else         goto         package      switch
const        fallthrough  if           range        type
continue     for          import       return       var
```</pre>
    </div>
  </div>
</div>
</template>

<style>

/*
** Layout CSS
*/

* {
  box-sizing: border-box;
}

::-webkit-scrollbar {
  background-color: transparent;
  width: 8px;
}

::-webkit-scrollbar-thumb {
  background-color: #dedede50

}

:root {
  --radius-size: 4px;
  --border-style: solid;
  --demo-display: none;
  --content-display: block;
}

.right-border {
  border-style: var(--border-style);
  border-color: #e5e5e5;
  border-width: 0 1px 0 0;
}

html, body {
  height: 100%;
}

body {
  margin: 0;
  overflow: hidden;
}

body {
  display: flex;
  /* align-items: center; */
  justify-content: center;
  font-size: 14px;
  font-family: Menlo, Monaco, "Courier New", monospace;
  background: #017972;
}

#app {
  width: 100%;
  height: 100%;
}

.root {
  display: flex;
  width: 100%;
  height: 100%;
  /* width: 1224px; */
  max-width: 1400px;
  max-height: 865px;
  margin: 0 auto;
  background: white;
  position: relative;

  border-style: var(--border-style);
  border-color: #e5e5e5;
  border-width: 1px;
  border-radius: var(--radius-size);
}

.root.fullscreen {
  max-width: none;
  max-height: none;
  border-radius: 0;
}

/* left */
.nav {
  display: flex;
  flex-direction: column;
  background: #20272b;
  min-width: 180px;
  flex: 15 15 238px;
  color: white;
  margin: -1px 0 -1px -1px;
  border-radius: var(--radius-size) 0 0 var(--radius-size);
}

.nav-head {
  min-height: calc(38px + 1px);
}

.fullscreen .nav {
  border-radius: 0;
}

/* center */
.note-list {
  display: flex;
  flex-direction: column;
  min-width: 160px;
  flex: 5 80 336px;
}

/* right */
.note {
  display: flex;
  overflow: hidden;
  flex-direction: column;
  flex: 80 5 650px;
}

.head {
  min-height: 38px;
  max-height: 38px;
  background: #f5f5f5;
  border-style: var(--border-style);
  border-color: #e5e5e5;
  border-width: 0 0 1px 0;
}

.note .head {
  border-radius: 0 var(--radius-size) 0 0;
}

.fullscreen .note .head {
  border-radius: 0;
}

.content {
  overflow: auto;
}

.content > div, ul, pre {
  display: var(--content-display)
}

.content > .demo {
  display: var(--demo-display)
}

.menu-icon {
  position: absolute;
  color: white;
  font-size: 24px;
  left: 12px;
  top: 7px;
  cursor: pointer;
  display: none;
}

.head .icon {
  font-size: 18px;
  margin-top: 9px;
  cursor: pointer;
}

.head .icon.active {
  color: #ef6d02
}

.demo {
  display: var(--demo-display);
}

/* @media (max-width: 600px) { */
  .menu-icon {
    display: block;
  }

  .close .nav, .close .note-list {
    display: none;
  }

  .close .menu-icon {
    color: black;
  }

  /* ----------------- */
  /* for sticky scroll */
  .close .menu-icon {
    z-index: 1;
    position: fixed;
  }
  .close .note .head {
    position: fixed;
    width: 100%;
    border-style: var(--border-style);
    border-width: 1px;
    margin: -1px;
  }
  .close .note .content {
    margin-top: calc(38px + 1px)
  }
  /* ----------------- */
/* } */

/*
** .nav CSS
*/
.nav .item-1 {
  padding: 8px 12px;
  cursor: pointer;
  display: flex;
  height: 30px;
  align-items: center;
}
.nav .item-1 .icon {
  font-size: 18px;
  display: inline-block;
  margin-right: 5px;
}

.nav .item-1 .text {
  flex: 1;
}

.nav .note-count {
  font-size: 12px;
  /* font-family: monospace; */
}

.nav .item-1:hover,
.nav .item-2:hover {
  background: #292f35;
}

.nav .item-1.active,
.nav .item-2.active {
  background: #313831;
}
.nav ul {
  margin: 0;
  padding: 0;
}
.nav .item-2 {
  height: 30px;
  padding: 4px 12px 4px 35px;
  cursor: pointer;
  display: flex;
  align-items: center;
}
.nav .item-2 .text {
  flex: 1;
  font-size: 13px;;
}

/*
** .note-list CSS
*/
.note-list>.head {
  display: flex;
  position: relative;
}
.note-list>.head .search {
  height: 22px;
  border-radius: 4px;
  border-width: 1px;
  border-color: #e5e5e5;
  border-style: solid;
  width: 100%;
  margin: 6px 12px;
  padding-left: 6px;
}
.note-list>.head .search:focus {
  outline: none;
}

.note-list>.head .search-icon {
  position: absolute;
  right: 16px;
  top: 12px;
  margin: 0;
}

.note-list ul {
  margin: 0px;
  padding: 0px;
}
.note-list li {
  display: block;
  padding: 6px 12px;
  cursor: pointer;
}

.note-list li:hover {
  background: #f5f5f5;
}

.note-list li.active {
  background: #ebebeb;
}

/*
** .note CSS
*/
.note .head {
  padding-left: 12px;
}
.close .note .head {

  padding-left: calc(24px + 12px + 12px);
}

.note .content {
  height: 100%;
}

.note .content > pre {
  margin: 0;
  padding: 8px 12px;
  /* font-size: 18px; */
}

.note .content > .html {
  padding: 0 12px;
}

.note .content > .readwrite {
  padding-top: 8px;
  height: 100%;
}

/*
** print mode
*/
@media print {
  /* visual */
  :root {
    --border-style: none;
    --radius-size: 0;
  }
  /* content */
  body {
    height: initial;
    overflow: hidden;
    background: none;
  }
  .note pre {
    white-space: pre-wrap;
  }
  /* interaction */
  .nav, .note-list, .menu-icon, .note .head {
    display: none;
  }
  .content {
    overflow: hidden;
  }
  .close .note .content {
    margin-top: 0   /* the sticky header in close mode is gone, so reset margin-top */
  }
}
</style>

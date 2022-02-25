export interface NoteItem {
  title: string
  path: string
  content: string
  tags?: string[]
}

export enum NoteMode {
  ReadOnly = 1,
  ReadWrite,
  Html,
}

export enum Tag {
  All = 0,
  Tagged,
  Untagged,
}

export enum State {
  MenuOpen = "state.menuOpen",
  SelectedPath = "state.selectedPath",
  // SearchText = "state.searchText",
}

export interface TagInfo {
  total: number
  tagged: number
  untagged: number
  tags: [string, number][]
}
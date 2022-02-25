<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import cookies from "browser-cookies"
import { NoteItem, State } from '../types';

const props = defineProps<{
  items: NoteItem[]
}>()
const emit = defineEmits<{
  (e: 'select', note: NoteItem): void
}>()

const selectedPath = ref('')
let restoreStateOnce = true

onMounted(() => {
  // props.items may be empty here
  restoreState()
})

watch(props, ()=> {
  restoreState()
})

function restoreState() {
  if (restoreStateOnce && props.items.length > 0) {
    restoreStateOnce = false
    const path = cookies.get(State.SelectedPath)
    if (path !== '') {
      for (let note of props.items) {
        if (note.path === path) {
          onItemClick(note)
          break
        }
      }
    }
  }
}

function onItemClick(note: NoteItem) {
  selectedPath.value = note.path
  emit('select', note)
  cookies.set(State.SelectedPath, selectedPath.value, {expires: 365})
}
</script>

<template>
  <ul class="">
    <li v-for="(note, index) in items" :key="index"
        @click="onItemClick(note)"
        :class="{active: selectedPath == note.path}"
      >{{ note.title }}</li>
  </ul>
</template>
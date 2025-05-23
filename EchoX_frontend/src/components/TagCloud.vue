<template>
  <div class="tag-cloud-container">
    <vue-word-cloud
      :words="formattedTags"
      :color="([text, weight, tag]) => tag.color"
      font-family="Roboto, sans-serif"
      :font-size-ratio="5"
      :rotation="0"
      :spacing="5"
      @wordclick="handleWordClick"
    />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import VueWordCloud from 'vuewordcloud'
import type { Tag } from '../types/record'

interface TagWithCount extends Tag {
  count: number
}

const props = defineProps<{
  tags: TagWithCount[]
}>()

const emit = defineEmits<{
  (e: 'select', tag: Tag): void
}>()

const formattedTags = computed(() => {
  return props.tags.map(tag => [tag.name, tag.count, tag])
})

const handleWordClick = (word: string, weight: number, tag: Tag) => {
  emit('select', tag)
}
</script>

<style scoped>
.tag-cloud-container {
  width: 100%;
  height: 300px;
}
</style> 
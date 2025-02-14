<script setup lang="ts">
defineProps({
  text: String,
  disabled: Boolean,
  leftInnerIcon: Object,
  onClick: Function,
})
</script>

<template>
  <button
    type="button"
    @click="onClick !== undefined && onClick()"
    :disabled="disabled"
    class="focus-visible:outline-offset inline-flex items-center gap-x-2 rounded-md bg-red px-3 py-1.5 text-white shadow-md drop-shadow-md transition duration-200 ease-in-out hover:bg-opacity-80 focus:-outline-offset-2 focus-visible:outline focus-visible:outline-2 focus-visible:outline-dark disabled:bg-dark disabled:bg-opacity-50"
  >
    <component :is="leftInnerIcon" class="size-5" />
    <span class="block pb-1">{{ text }}</span>
  </button>
</template>

<style scoped></style>

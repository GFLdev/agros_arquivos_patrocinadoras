<script setup lang="ts">
import { PhPlus } from '@phosphor-icons/vue'
import type { PropType } from 'vue'

defineProps({
  text: {
    type: String,
    required: true,
  },
  disabled: {
    type: Boolean,
    default: false,
  },
  onClick: {
    type: Function as PropType<() => unknown>,
    default: (): void => {},
  },
})
</script>

<template>
  <button
    type="button"
    class="focus-visible:outline-offset inline-flex w-full items-center justify-center gap-x-2 rounded-md border-2 border-primary bg-primary px-3 py-1.5 text-base font-light text-white shadow-lg drop-shadow-lg transition duration-200 ease-in-out hover:font-normal focus:-outline-offset-2 focus-visible:bg-secondary focus-visible:outline focus-visible:outline-2 focus-visible:outline-primary enabled:hover:bg-white enabled:hover:text-primary disabled:bg-white disabled:text-primary"
    :disabled="disabled"
    @click.prevent="() => onClick !== undefined && onClick()"
  >
    <PhPlus class="size-5" />
    <span class="mb-1">{{ text }}</span>
  </button>
</template>

<style scoped></style>

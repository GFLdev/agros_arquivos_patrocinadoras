<script setup lang="ts">
import { PhCircleNotch, PhPaperPlaneTilt } from '@phosphor-icons/vue'

defineProps({
  text: {
    type: String,
    required: true,
  },
  loading: {
    type: Boolean,
    default: false,
  },
  loadingText: {
    type: String,
    default: 'Carregando...',
  },
  disabled: {
    type: Boolean,
    default: false,
  },
  leftInnerIcon: {
    type: Object,
    default: PhPaperPlaneTilt,
  },
})
</script>

<template>
  <button
    type="submit"
    :disabled="disabled"
    class="focus-visible:outline-offset inline-flex items-center gap-x-2 rounded-md bg-primary px-3 py-1.5 text-white shadow-md drop-shadow-md transition duration-200 ease-in-out focus:-outline-offset-2 focus-visible:outline focus-visible:outline-2 focus-visible:outline-dark enabled:hover:bg-opacity-80 disabled:bg-dark disabled:bg-opacity-50"
  >
    <component v-if="!loading" :is="leftInnerIcon" class="size-5" />
    <PhCircleNotch v-else class="size-5 animate-spin" aria-hidden="true" />
    <span class="block pb-1">{{ loading ? loadingText : text }}</span>
  </button>
</template>

<style scoped></style>

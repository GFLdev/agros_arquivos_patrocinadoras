<script setup lang="ts">
import { nextTick, onUpdated, ref } from 'vue'
import { PhCaretDown, PhPencil, PhTrash } from '@phosphor-icons/vue'

defineProps({
  title: String,
  admin: Boolean,
  first: Boolean,
  last: Boolean,
  editHandler: Function,
  deleteHandler: Function,
})

const open = ref<boolean>(false)
const contentHeight = ref('0px')

const toggle = () => {
  open.value = !open.value
}

// Calcula altura do conteúdo dinamicamente
async function adjustHeight() {
  await nextTick()
  const content = document.querySelector(`.accordion-content`) as HTMLElement
  if (content) {
    contentHeight.value = open.value ? `${content.scrollHeight}px` : '0px'
  }
}

onUpdated(() => {
  adjustHeight()
})
</script>

<template>
  <section
    class="max-h-fit w-full overflow-hidden border-x-2 border-b-2 border-primary bg-primary"
    :class="`${first ? 'rounded-t-md' : ''} ${last ? 'rounded-b-md' : ''}`"
  >
    <div class="cursor-pointer bg-primary text-center text-white shadow-lg drop-shadow-lg">
      <div class="grid w-full grid-cols-3 items-center justify-between px-4 py-2" @click.prevent="toggle">
        <div class="justify-self-start">
          <PhCaretDown class="size-5 transition-all duration-300" :class="open ? 'rotate-180' : ''" />
        </div>
        <div class="mb-1 select-none justify-self-center">{{ title }}</div>
        <div class="flex flex-row gap-x-2 justify-self-end" v-if="admin">
          <button
            @click.prevent="
              (e) => {
                e.stopPropagation()
                editHandler !== undefined && editHandler()
              }
            "
            class="z-10 rounded-lg p-1 text-green text-opacity-80 transition-all duration-200 ease-in-out hover:bg-green hover:bg-opacity-70 hover:text-white"
          >
            <PhPencil class="size-5" weight="fill" />
          </button>
          <button
            @click.prevent="
              (e) => {
                e.stopPropagation()
                deleteHandler !== undefined && deleteHandler()
              }
            "
            class="z-10 rounded-lg p-1 text-red text-opacity-80 transition-all duration-200 ease-in-out hover:bg-red hover:bg-opacity-70 hover:text-white"
          >
            <PhTrash class="size-5" weight="fill" />
          </button>
        </div>
      </div>
    </div>
    <div
      class="-z-10 overflow-hidden bg-white transition-all duration-300 ease-in-out"
      :style="{ maxHeight: contentHeight }"
    >
      <div
        class="accordion-content relative flex w-full flex-col items-center justify-center gap-y-4 px-8 py-4 transition-all duration-300 ease-in-out"
        :class="`${open ? 'animate-enter-from-top' : 'animate-exit-to-top'} ${last ? 'rounded-b-md' : ''}`"
      >
        <slot name="content"></slot>
      </div>
    </div>
  </section>
</template>

<style scoped></style>

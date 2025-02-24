<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, onMounted, ref } from 'vue'
import { PhCaretDown, PhPencil, PhTrash } from '@phosphor-icons/vue'

const props = defineProps({
  title: String,
  admin: Boolean,
  first: Boolean,
  last: Boolean,
  editHandler: Function,
  deleteHandler: Function,
  contentHeight: Number,
})

const open = ref<boolean>(false)
const maxContentHeight = ref<number>(0)

const contentNode = ref<Node | null>(null)
const observer = ref<MutationObserver | null>(null)

const heightModel = defineModel<number | null>()

const toggle = () => {
  open.value = !open.value
}

// Calcula altura do conteúdo dinamicamente
async function adjustHeight() {
  await nextTick()
  if (contentNode.value) {
    const el = contentNode.value as HTMLElement
    const childHeight: number = computed(() => props.contentHeight).value ?? 0
    maxContentHeight.value = open.value ? el.scrollHeight + childHeight : 0
    heightModel.value = maxContentHeight.value
  }
}

onMounted(() => {
  observer.value = new MutationObserver(async () => {
    await adjustHeight()
  })
  if (contentNode.value) {
    observer.value.observe(contentNode.value, {
      childList: true,
      subtree: true,
      attributes: true,
      characterData: true,
    })
  }
})

onBeforeUnmount(() => {
  if (observer.value) {
    observer.value.disconnect()
  }
})
</script>

<template>
  <section
    class="max-h-fit w-full overflow-hidden border-x-2 border-b-2 border-primary bg-primary"
    :class="`${first ? 'rounded-t-md' : ''} ${last ? 'rounded-b-md' : ''}`"
  >
    <div class="cursor-pointer bg-primary text-center text-white shadow-lg drop-shadow-lg">
      <div class="grid w-full grid-flow-col grid-cols-5 items-center gap-8 px-4 py-2" @click.prevent="toggle">
        <div class="col-span-1 justify-self-start">
          <PhCaretDown class="size-5 transition-all duration-300" :class="open ? 'rotate-180' : ''" />
        </div>
        <div class="col-span-3 mb-1 select-none justify-self-center text-wrap text-center font-light">
          <p>{{ title }}</p>
        </div>
        <div class="col-span-1 flex flex-row gap-x-2 justify-self-end" v-if="admin">
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
      :style="{ maxHeight: `${maxContentHeight}px` }"
    >
      <div
        ref="contentNode"
        class="relative flex w-full flex-col items-center justify-center gap-y-8 px-8 py-4 transition-all duration-300 ease-in-out"
        :class="`${open ? 'animate-enter-from-top' : 'animate-exit-to-top'} ${last ? 'rounded-b-md' : ''}`"
      >
        <slot name="content"></slot>
      </div>
    </div>
  </section>
</template>

<style scoped></style>

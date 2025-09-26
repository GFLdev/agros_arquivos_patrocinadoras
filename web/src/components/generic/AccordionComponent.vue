<script setup lang="ts">
import { computed, type ModelRef, nextTick, onBeforeUnmount, onMounted, type PropType, type Ref, ref } from 'vue'
import { PhCaretDown, PhPencil, PhTrash } from '@phosphor-icons/vue'

const props: {
  readonly title: string
  readonly admin: boolean
  readonly first: boolean
  readonly last: boolean
  readonly editHandler: () => unknown
  readonly deleteHandler: () => unknown
  readonly contentHeight: number
} = defineProps({
  title: {
    type: String,
    required: true,
  },
  admin: {
    type: Boolean,
    default: false,
  },
  first: {
    type: Boolean,
    default: true,
  },
  last: {
    type: Boolean,
    default: true,
  },
  editHandler: {
    type: Function as PropType<() => unknown>,
    default: (): void => {},
  },
  deleteHandler: {
    type: Function as PropType<() => unknown>,
    default: (): void => {},
  },
  contentHeight: {
    type: Number,
    default: 0,
  },
})

// Estado de abertura do accordion
const open: Ref<boolean> = ref<boolean>(false)

// Observador
const observer: Ref<MutationObserver | undefined> = ref<MutationObserver>()

// Estado para cálculo da altura máxima do conteúdo interno
const maxContentHeight: Ref<number> = ref<number>(0)

// Nó do conteúdo interno do accordion
const contentNode: Ref<Node | undefined> = ref<Node>()

// Altura máxima do conteúdo interno do accordion
const heightModel: ModelRef<number | undefined> = defineModel<number>()

/**
 * Alterna o estado do valor 'open' entre true e false.
 *
 * @return {void} Não retorna nenhum valor.
 */
function toggle(): void {
  open.value = !open.value
}

/**
 * Ajusta dinamicamente a altura de um contêiner com base no conteúdo e no estado.
 * Esta função é assíncrona e garante que as atualizações do DOM sejam tratadas corretamente
 * antes de ajustar a altura.
 *
 * @return {Promise<void>} Uma promessa que é resolvida quando o ajuste de altura é concluído.
 */
async function adjustHeight(): Promise<void> {
  await nextTick()
  if (contentNode.value) {
    const el = contentNode.value as HTMLElement
    const childHeight: number = computed(() => props.contentHeight).value ?? 0
    maxContentHeight.value = open.value ? el.scrollHeight + childHeight : 0
    heightModel.value = maxContentHeight.value
  }
}

// Fechar observador
onBeforeUnmount((): void => {
  if (observer.value) observer.value.disconnect()
})

// Definir observador na primeira renderização
onMounted(() => {
  observer.value = new MutationObserver(async (): Promise<void> => {
    await adjustHeight()
  })

  // Observar quaisquer mudanças no nó do conteúdo interno do accordion
  if (contentNode.value) {
    observer.value.observe(contentNode.value, {
      childList: true,
      subtree: true,
      attributes: true,
      characterData: true,
    })
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
                editHandler()
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
                deleteHandler()
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

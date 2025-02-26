<script setup lang="ts">
import { PhCircleNotch, PhPencil, PhTrash } from '@phosphor-icons/vue'
import { type PropType, type Ref, ref } from 'vue'
import { fileIcon } from '@/utils/file.ts'

const props: {
  readonly name: string
  readonly lastModified: number
  readonly admin: boolean
  readonly mimeType: string
  readonly editHandler: () => unknown
  readonly deleteHandler: () => unknown
  readonly downloadHandler: () => unknown
} = defineProps({
  name: {
    type: String,
    required: true,
  },
  lastModified: {
    type: Number,
    required: true,
  },
  admin: {
    type: Boolean,
    default: false,
  },
  mimeType: {
    type: String,
    default: 'application/octet-stream',
  },
  editHandler: {
    type: Function as PropType<() => unknown>,
    default: (): void => {},
  },
  deleteHandler: {
    type: Function as PropType<() => unknown>,
    default: (): void => {},
  },
  downloadHandler: {
    type: Function as PropType<() => unknown>,
    default: (): void => {},
  },
})

// Status de carregamento/download
const downloading: Ref<boolean> = ref<boolean>(false)

/**
 * Converte um timestamp Unix fornecido em uma string de data formatada.
 *
 * @param {number} unix_ts - O timestamp Unix (em segundos) a ser convertido.
 * @return {string} Uma string de data formatada no formato "DD/MM/YYYY, HH:mm".
 */
function getFormattedDate(unix_ts: number): string {
  const date: Date = new Date(unix_ts * 1000)
  const year: string = date.getFullYear().toString()
  const month: string = (date.getMonth() + 1).toString().padStart(2, '0')
  const day: string = date.getDate().toString().padStart(2, '0')
  const hour: string = date.getHours().toString().padStart(2, '0')
  const min: string = date.getMinutes().toString().padStart(2, '0')
  return day + '/' + month + '/' + year + ', ' + hour + ':' + min
}

/**
 * Lida com a execução assíncrona do processo de download.
 * O método garante que o estado de download seja gerenciado adequadamente
 * antes e depois da invocação do manipulador de download fornecido.
 *
 * @return {Promise<void>} Uma promise que é resolvida quando o processo de download é concluído.
 */
async function downloadWrapper(): Promise<void> {
  downloading.value = true
  await props.downloadHandler()
  downloading.value = false
}
</script>

<template>
  <div class="flex w-full flex-row items-center justify-between gap-4">
    <div
      class="flex flex-row items-center gap-4 transition-all duration-200 ease-in-out"
      :class="`${!downloading ? 'cursor-pointer hover:text-gray' : ''}`"
      @click.prevent="downloading ? void 0 : downloadWrapper()"
    >
      <component v-if="!downloading" class="size-5" :is="fileIcon(mimeType)" />
      <PhCircleNotch v-else class="size-5 animate-spin" />
      <div class="pb-1">{{ name }}</div>
    </div>
    <div class="flex flex-row items-center gap-4 pb-1">
      <span class="text-sm font-light">Atualizado: {{ getFormattedDate(lastModified) }}</span>
      <div class="flex flex-row gap-x-2 justify-self-end" v-if="admin">
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
</template>

<style scoped></style>

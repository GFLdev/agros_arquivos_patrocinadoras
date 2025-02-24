<script setup lang="ts">
import { PhCircleNotch, PhFileArrowDown, PhPencil, PhTrash } from '@phosphor-icons/vue'
import { ref } from 'vue'
import { fileIcon } from '@/utils/file.ts'

const props = defineProps({
  name: {
    type: String,
    required: true,
  },
  lastModified: {
    type: Number,
    required: true,
  },
  admin: Boolean,
  mimeType: String,
  editHandler: Function,
  deleteHandler: Function,
  downloadHandler: Function,
})

const downloading = ref<boolean>(false)

function getFormattedDate(unix_ts: number) {
  const date = new Date(unix_ts * 1000)
  const year = date.getFullYear().toString()
  const month = (date.getMonth() + 1).toString().padStart(2, '0')
  const day = date.getDate().toString().padStart(2, '0')
  const hour = date.getHours().toString().padStart(2, '0')
  const min = date.getMinutes().toString().padStart(2, '0')
  return day + '/' + month + '/' + year + ', ' + hour + ':' + min
}

async function downloadFile() {
  if (props.downloadHandler === undefined) return

  downloading.value = true
  await props.downloadHandler()
  downloading.value = false
}
</script>

<template>
  <div class="flex w-full flex-row items-center justify-between gap-4">
    <div
      class="flex flex-row items-center gap-4 transition-all duration-200 ease-in-out"
      :class="`${downloadHandler !== undefined || !downloading ? 'cursor-pointer hover:text-gray' : ''}`"
      @click.prevent="downloading ? void 0 : downloadFile()"
    >
      <component v-if="!downloading" class="size-5" :is="fileIcon(mimeType ?? '')" />
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
</template>

<style scoped></style>

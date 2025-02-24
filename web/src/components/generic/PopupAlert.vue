<script setup lang="ts">
import { PhCheckCircle, PhInfo, PhWarningCircle, PhXCircle } from '@phosphor-icons/vue'
import { computed, type PropType, ref, watch } from 'vue'
import { AlertType } from '@/@types/Enumerations.ts'

const props = defineProps({
  text: {
    type: String,
    required: true,
  },
  duration: {
    type: Number,
    default: 3000,
  },
  type: Number as PropType<AlertType>,
})

const closedCalled = ref<boolean>(false)

const showModel = defineModel<boolean>()

const icons = {
  [AlertType.Success]: PhCheckCircle,
  [AlertType.Error]: PhXCircle,
  [AlertType.Warning]: PhWarningCircle,
  [AlertType.Info]: PhInfo,
}

const styles = {
  [AlertType.Success]: `text-white bg-green`,
  [AlertType.Error]: `text-white bg-red`,
  [AlertType.Warning]: `text-white bg-orange`,
  [AlertType.Info]: `text-white bg-secondary`,
}

// Função para fechar popup
function close() {
  closedCalled.value = true
  showModel.value = false
}

watch(
  () => showModel.value,
  () => setTimeout(close, computed(() => props.duration).value),
)
</script>

<template>
  <div
    class="fixed top-4 z-50 flex w-full flex-col items-center justify-center"
    :class="`${showModel ? 'animate-enter-from-top' : closedCalled ? 'animate-exit-to-top' : 'hidden'}`"
  >
    <div
      class="inline-flex h-fit w-max max-w-96 items-center gap-x-2 rounded-md px-5 py-2.5 shadow-md drop-shadow-md"
      :class="`${type !== undefined ? styles[type] : styles[AlertType.Info]}`"
    >
      <component :is="type !== undefined ? icons[type] : icons[AlertType.Info]" class="size-6" weight="fill" />
      <p class="font-lato">{{ text }}</p>
    </div>
  </div>
</template>

<style scoped></style>

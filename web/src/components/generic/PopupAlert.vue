<script setup lang="ts">
import { PhCheckCircle, PhInfo, PhWarningCircle, PhXCircle } from '@phosphor-icons/vue'
import { computed, type ModelRef, type PropType, type Ref, ref, watch } from 'vue'
import { AlertType } from '@/@types/Enumerations.ts'

const props: {
  readonly text: string
  readonly duration: number
  readonly type: AlertType
} = defineProps({
  text: {
    type: String,
    required: true,
  },
  duration: {
    type: Number,
    default: 3000,
  },
  type: {
    type: Number as PropType<AlertType>,
    default: AlertType.Info,
  },
})

// Ícones para cada tipo de alerta
const icons = {
  [AlertType.Success]: PhCheckCircle,
  [AlertType.Error]: PhXCircle,
  [AlertType.Warning]: PhWarningCircle,
  [AlertType.Info]: PhInfo,
}

// Estilos (tailwind) para cada tipo de alerta
const styles = {
  [AlertType.Success]: `text-white bg-green`,
  [AlertType.Error]: `text-white bg-red`,
  [AlertType.Warning]: `text-white bg-orange`,
  [AlertType.Info]: `text-white bg-secondary`,
}

// Estado para caso o alerta já tenha sido fechado. Usado para não ter
// o bug visual de animação do alerta fechando mesmo que nunca tenha sido aberto
const closedCalled: Ref<boolean> = ref<boolean>(false)

// Estado de visibilidade
const showModel: ModelRef<boolean | undefined> = defineModel<boolean>()

/**
 * Fecha a visualização do modelo atual e atualiza os indicadores de status.
 *
 * @return {void} Este método não retorna nenhum valor.
 */
function close(): void {
  closedCalled.value = true
  showModel.value = false
}

// Fechar alerta dada a sua duração
watch(
  (): boolean | undefined => showModel.value,
  (): number => setTimeout(close, computed((): number => props.duration).value),
)
</script>

<template>
  <div
    class="fixed top-4 z-50 flex w-full flex-col items-center justify-center"
    :class="`${showModel ? 'animate-enter-from-top' : closedCalled ? 'animate-exit-to-top' : 'hidden'}`"
  >
    <div
      class="inline-flex h-fit w-max max-w-96 items-center gap-x-2 rounded-md px-5 py-2.5 shadow-md drop-shadow-md"
      :class="styles[type]"
    >
      <component :is="icons[type]" class="size-6" weight="fill" />
      <p class="font-lato">{{ text }}</p>
    </div>
  </div>
</template>

<style scoped></style>

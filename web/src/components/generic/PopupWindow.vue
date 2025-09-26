<script setup lang="ts">
import { PhX } from '@phosphor-icons/vue'
import { type ModelRef, type Ref, ref } from 'vue'

defineProps({
  title: {
    type: String,
    required: true,
  },
})

// Estado para caso a janela já tenha sido fechada. Usado para não ter
// o bug visual de animação da janela fechando mesmo que nunca tenha sido aberta
const closedCalled: Ref<boolean> = ref<boolean>(false)

// Estado de visibilidade
const showModel: ModelRef<boolean | undefined> = defineModel<boolean>()

/**
 * Fecha o modelo atual atualizando as variáveis de estado relevantes.
 *
 * @return {void} Este método não retorna nenhum valor.
 */
function close(): void {
  closedCalled.value = true
  showModel.value = false
}
</script>

<template>
  <div
    class="fixed inset-0 z-40 h-screen w-full cursor-pointer bg-black bg-opacity-50 transition-all duration-1000"
    :class="`${showModel ? 'animate-open-with-fade' : closedCalled ? 'animate-close-with-fade' : 'hidden'}`"
    @click="close"
    v-bind="$attrs"
  ></div>
  <section
    class="fixed left-1/2 top-1/2 z-40 flex h-fit max-h-screen w-screen -translate-x-1/2 -translate-y-1/2 transform flex-col gap-4 bg-white px-8 pb-8 pt-4 shadow-lg drop-shadow-lg sm:w-fit sm:min-w-[28rem] sm:rounded-md"
    :class="`${showModel ? 'animate-open-with-fade' : closedCalled ? 'animate-close-with-fade' : 'hidden'}`"
    v-bind="$attrs"
  >
    <div class="relative">
      <button class="absolute -right-4 top-0 transition-colors duration-200 hover:text-red" @click="close">
        <PhX size="24" />
      </button>
      <h2 class="text-center text-xl font-light">{{ title.toUpperCase() }}</h2>
    </div>
    <div>
      <slot />
    </div>
  </section>
</template>

<style scoped></style>

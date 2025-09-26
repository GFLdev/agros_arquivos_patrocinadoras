<script setup lang="ts">
import { computed, type ComputedRef, getCurrentInstance, type ModelRef, type ShallowRef, useTemplateRef } from 'vue'
import { PhFolderOpen } from '@phosphor-icons/vue'
import { isFileEmpty } from '@/utils/validate.ts'

defineProps({
  label: {
    type: String,
    required: true,
  },
  placeholder: {
    type: String,
    required: true,
  },
  leftInnerIcon: {
    type: Object,
    default: PhFolderOpen,
  },
  disabled: {
    type: Boolean,
    default: false,
  },
  required: {
    type: Boolean,
    default: false,
  },
})

// Elemento do tipo input
const inputElement: Readonly<ShallowRef<HTMLInputElement | null>> = useTemplateRef('input-element')

// UID da instância do componente
const uid: number = getCurrentInstance()!.uid

// Computada para identificar se o model, que contém o arquivo, é vazio
const fileValue: ComputedRef<File | null> = computed((): File | null => {
  return fileModel.value ? (isFileEmpty(fileModel.value) ? null : fileModel.value) : null
})

// Model contendo o arquivo
const fileModel: ModelRef<File | undefined> = defineModel<File>()

/**
 * Lida com o evento de alteração de um elemento de entrada de arquivo, atribuindo o primeiro arquivo selecionado ao
 * valor do modelo.
 *
 * @param {Event} e - O evento disparado pela mudança no elemento de entrada de arquivo.
 * @return {void} Não retorna nenhum valor.
 */
function handleFileChange(e: Event): void {
  e.preventDefault()
  const target: HTMLInputElement = e.target as HTMLInputElement
  fileModel.value = target.files ? target.files[0] : new File([], '')
}

/**
 * Dispara o evento de clique do elemento de entrada de arquivo ao ser invocado, permitindo que o diálogo de seleção de
 * arquivo seja aberto.
 *
 * @param {Event | KeyboardEvent} e - O evento que dispara esta função. Pode ser um evento geral ou um evento de
 * teclado. Se for um evento de teclado, apenas as teclas Enter ou Espaço irão disparar a ação.
 * @return {void} Esta função não retorna nenhum valor.
 */
function triggerFileInput(e: Event | KeyboardEvent): void {
  e.preventDefault()
  if (e instanceof KeyboardEvent) {
    if (e.key === 'Enter' || e.key === ' ') {
      inputElement.value?.click()
    }
  } else {
    inputElement.value?.click()
  }
}
</script>

<template>
  <div class="relative mt-2 grid grid-cols-1">
    <input
      :id="uid.toString()"
      ref="input-element"
      type="file"
      :name="uid.toString()"
      :multiple="false"
      class="hidden"
      @change="handleFileChange"
      :required="required"
      :disabled="disabled"
    />
    <label
      :for="uid.toString()"
      @keydown="triggerFileInput"
      @click.prevent="triggerFileInput"
      tabindex="0"
      class="peer col-start-1 row-start-1 block w-full cursor-pointer rounded-md bg-white py-2 pr-3 text-base text-gray outline outline-1 -outline-offset-1 outline-gray placeholder:text-gray focus:text-dark focus:outline focus:outline-2 focus:-outline-offset-2 focus:outline-dark sm:text-sm/6"
    >
      <span class="absolute -top-2 left-2 inline-block rounded-lg bg-white px-1 text-xs font-medium"
        >{{ label }}<b v-if="required">&nbsp;*</b></span
      >
      <span class="absolute top-1.5 col-start-1 row-start-1 pl-10">
        {{ fileValue ? fileValue.name : placeholder }}
      </span>
      <component :is="leftInnerIcon" class="col-start-1 row-start-1 ml-3 size-5 self-center" aria-hidden="true" />
    </label>
  </div>
</template>

<style scoped>
*:disabled {
  @apply opacity-50;
}

*:enabled {
  @apply opacity-100;
}
</style>

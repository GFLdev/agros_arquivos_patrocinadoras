<script setup lang="ts">
import { computed, getCurrentInstance, useTemplateRef } from 'vue'

defineProps({
  label: {
    type: String,
    required: true,
  },
  placeholder: {
    type: String,
    required: true,
  },
  leftInnerIcon: Object,
  disabled: Boolean,
  required: Boolean,
})

const inputElement = useTemplateRef('input-element')
const model = defineModel<File>()
const uid = getCurrentInstance()!.uid

// Computada para identificar se o model é vazio
const fileValue = computed(() => {
  return model.value && model.value.size === 0 && model.value.name === '' ? null : model.value
})

// Evento de captura do arquivo
function handleFileChange(e: Event) {
  e.preventDefault()
  const target = e.target as HTMLInputElement
  model.value = target.files ? target.files[0] : new File([], '')
}

// Dispara o clique no input quando a label estiver focada e o usuário pressionar Enter ou Espaço
function triggerFileInput(e: Event | KeyboardEvent) {
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
      <p class="absolute top-1 col-start-1 row-start-1" :class="`${leftInnerIcon ? 'pl-10' : 'pl-3'}`">
        {{ fileValue ? fileValue.name : placeholder }}
      </p>
      <component
        v-if="leftInnerIcon"
        :is="leftInnerIcon"
        class="col-start-1 row-start-1 ml-3 size-5 self-center"
        aria-hidden="true"
      />
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

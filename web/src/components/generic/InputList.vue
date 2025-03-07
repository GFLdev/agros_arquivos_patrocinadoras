<script setup lang="ts">
import { getCurrentInstance, watch } from 'vue'

const props = defineProps({
  label: {
    type: String,
    required: true,
  },
  values: {
    type: Map<string, string>,
    required: true,
  },
  selected: {
    type: String,
    required: true,
  },
  leftInnerIcon: Object,
  disabled: Boolean,
  required: Boolean,
})

const model = defineModel<string>()
const uid = getCurrentInstance()!.uid

watch(
  () => props.selected,
  (newValue) => {
    if (model.value !== newValue) {
      model.value = newValue
    }
  },
  { immediate: true }, // Executa mesmo na primeira renderização
)
</script>

<template>
  <div class="relative mt-2 grid grid-cols-1">
    <select
      :id="uid.toString()"
      :name="uid.toString()"
      class="peer col-start-1 row-start-1 block w-full rounded-md bg-white py-[9.5px] pr-3 text-base text-gray outline outline-1 -outline-offset-1 outline-gray placeholder:text-gray focus:text-dark focus:outline focus:outline-2 focus:-outline-offset-2 focus:outline-dark sm:text-sm/6"
      :class="`${leftInnerIcon ? 'pl-10' : 'pl-3'}`"
      v-model="model"
      :required="required"
      :disabled="disabled"
    >
      <option v-for="[key, value] of values" :key="key" :value="key" :disabled="key === selected">
        {{ value }}
      </option>
    </select>
    <label
      :for="uid.toString()"
      class="absolute -top-2 left-2 inline-block rounded-lg bg-white px-1 text-xs font-medium text-gray peer-focus:text-dark"
      >{{ label }}<b v-if="required">&nbsp;*</b></label
    >
    <component
      v-if="leftInnerIcon"
      :is="leftInnerIcon"
      class="pointer-events-none col-start-1 row-start-1 ml-3 size-5 self-center text-gray peer-focus:text-dark"
      aria-hidden="true"
    />
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

<script setup lang="ts">
import { PhEye, PhEyeClosed } from '@phosphor-icons/vue'
import { getCurrentInstance, onMounted, ref, watchEffect } from 'vue'

const props = defineProps({
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
  showable: Boolean,
  matched: Boolean,
})

const model = defineModel<string | null | undefined>()
const visible = ref<boolean>(false)
const uid = getCurrentInstance()!.uid

onMounted(() => {
  watchEffect(() => {
    if (props.disabled) {
      visible.value = false
    }
  })
})
</script>

<template>
  <div class="relative mt-2 flex">
    <div class="relative -mr-px grid grow grid-cols-1">
      <input
        :id="uid.toString()"
        :name="uid.toString()"
        :class="`peer col-start-1 row-start-1 block w-full bg-white py-1.5 ${leftInnerIcon ? 'pl-10' : 'pl-3'} ${showable ? 'rounded-l-md' : 'rounded-md'} pr-3 text-base text-gray outline outline-1 -outline-offset-1 outline-gray placeholder:text-base placeholder:text-gray focus:text-dark focus:outline focus:outline-2 focus:-outline-offset-2 focus:outline-dark sm:text-sm/6 sm:placeholder:text-sm/6`"
        :placeholder="placeholder"
        v-model="model"
        :disabled="disabled"
        :type="showable ? (visible ? 'text' : 'password') : 'password'"
      />
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
    <button
      type="button"
      v-if="showable"
      class="flex shrink-0 items-center gap-x-1.5 rounded-r-md bg-white px-3 py-2 text-sm font-semibold text-dark outline outline-1 -outline-offset-1 outline-gray transition-colors duration-200 hover:bg-gray hover:bg-opacity-30 focus:relative focus:outline focus:outline-2 focus:-outline-offset-2 focus:outline-dark"
      :disabled="disabled"
      @click="visible = !visible"
    >
      <PhEyeClosed v-if="visible" class="-ml-0.5 size-4 text-dark" aria-hidden="true" />
      <PhEye v-else class="-ml-0.5 size-4 text-dark" aria-hidden="true" />
    </button>
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

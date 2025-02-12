<script setup lang="ts">
import { ref } from 'vue'

const toggleShow = ref<boolean>(false)

const toggle = () => {
  toggleShow.value = !toggleShow.value
}
</script>

<template>
  <section>
    <div @click="toggle">
      <slot name="header"></slot>
    </div>
    <div :class="`${toggleShow ? 'block' : 'hidden'}`">
      <slot name="content"></slot>
    </div>
  </section>
</template>

<style scoped></style>

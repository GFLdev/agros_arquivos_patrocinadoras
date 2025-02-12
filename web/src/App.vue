<script setup lang="ts">
import { RouterView } from 'vue-router'
import { onMounted } from 'vue'
import { useAuthStore } from '@/stores/authStore.ts'

const authStore = useAuthStore()

onMounted(async () => {
  await authStore.getSession()
})
</script>

<template>
  <RouterView />
</template>

<style scoped></style>

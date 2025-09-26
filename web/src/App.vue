<script setup lang="ts">
import { RouterView } from 'vue-router'
import { onMounted } from 'vue'
import { useAuthStore } from '@/stores/authStore.ts'

// Pinia Store
const authStore: ReturnType<typeof useAuthStore> = useAuthStore()

// Obter sessão, para todos os componentes da SPA
onMounted(authStore.getSession)
</script>

<template>
  <RouterView />
</template>

<style scoped></style>

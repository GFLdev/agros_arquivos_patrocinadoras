<script setup lang="ts">
import LoginForm from '@/components/login/LoginForm.vue'
import Logo from '/logo.svg'
</script>

<template>
  <main
    class="flex h-screen w-screen flex-col items-center justify-center gap-8 bg-gradient-90 from-secondary to-primary"
  >
    <img :src="Logo" class="w-full max-w-[300px] select-none px-4" alt="Logo Agros" />
    <LoginForm />
  </main>
</template>

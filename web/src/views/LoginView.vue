<script setup lang="ts">
import LoginForm from '@/components/LoginForm.vue'
import Logo from '/logo.svg'
</script>

<template>
  <main
    class="from-secondary to-primary flex flex-col gap-8 h-screen w-screen items-center justify-center bg-gradient-90"
  >
      <img :src="Logo" class="max-w-[300px] w-full select-none px-4"  alt="Logo Agros"/>
      <LoginForm />
  </main>
</template>

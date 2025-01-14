<script setup lang="ts">
import LoginForm from '@/components/LoginForm.vue'
import Logo from '/logo.svg'
</script>

<template>
  <main id="login">
    <div class="flex flex-col items-center gap-8">
      <v-img :src="Logo" class="w-[300px]"></v-img>
      <LoginForm class="justify-self-center" />
    </div>
  </main>
</template>

<style scoped>
#login {
  @apply bg-gradient-90
  from-secondary
  to-primary
  h-screen
  w-screen
  flex
  justify-center
  items-center
  !important;
}
</style>

<script setup lang="ts">
// import Logo from '/logo.svg'
import { ref } from 'vue'

const inst = ref<string | null>(null)

// FIXME: Retirar
inst.value = 'ufv'.toUpperCase()
</script>

<template>
  <header class="flex w-screen items-center justify-center bg-gradient-90 from-secondary to-primary">
    <section class="grid w-full max-w-[1620px] grid-cols-2 items-center justify-between px-8 py-4 sm:grid-cols-3">
      <div class="hidden justify-self-start sm:block">
<!--        <v-img :src="Logo" class="w-[200px]"></v-img>-->
      </div>
      <div class="justify-self-start text-center sm:justify-self-center">
        <h1 class="text-xl text-agros-white sm:text-2xl">{{ inst }} | Espaço Patrocinadora</h1>
      </div>
      <div class="justify-self-end">
<!--        <v-btn id="exit-btn" :prepend-icon="PhSignOut" variant="plain">Sair</v-btn>-->
      </div>
    </section>
  </header>
</template>

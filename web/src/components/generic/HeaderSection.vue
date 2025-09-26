<script setup lang="ts">
import Logo from '/logo.svg'
import { PhSignOut } from '@phosphor-icons/vue'
import { useAuthStore } from '@/stores/authStore.ts'
defineProps({
  title: {
    type: String,
    default: 'Espaço Patrocinadora',
  },
})

// Pinia Store
const authStore: ReturnType<typeof useAuthStore> = useAuthStore()
</script>

<template>
  <header class="sticky top-0 z-10 animate-open-with-fade bg-gradient-90 from-secondary to-primary font-lato">
    <section
      class="grid h-full w-full max-w-[1620px] grid-cols-2 items-center justify-between px-8 py-4 sm:grid-cols-3"
    >
      <div class="hidden justify-self-start sm:block">
        <img :src="Logo" class="w-[200px]" alt="Logo Agros" />
      </div>
      <div class="text-center sm:justify-self-center">
        <h1 class="text-2xl text-white sm:text-3xl">{{ title }}</h1>
      </div>
      <div class="justify-self-end">
        <button
          type="button"
          class="focus-visible:outline-offset inline-flex items-center gap-x-2 rounded-md px-3 py-1.5 text-base text-red shadow-sm transition duration-200 ease-in-out focus:-outline-offset-2 focus-visible:outline focus-visible:outline-2 focus-visible:outline-red enabled:hover:bg-red enabled:hover:bg-opacity-70 enabled:hover:text-white"
          @click="authStore.logout"
        >
          <PhSignOut weight="bold" />
          <span class="mb-1 font-semibold">Sair</span>
        </button>
      </div>
    </section>
  </header>
</template>

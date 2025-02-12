<script setup lang="ts">
import Header from '@/components/generic/HeaderSection.vue'
import AccordionComponent from '@/components/generic/AccordionComponent.vue'
import { onBeforeMount, ref } from 'vue'
import type { CategModel, FileModel, UserModel } from '@/@types/Responses.ts'
import { getAllCategories, getAllFiles, getAllUsers } from '@/services/queries.ts'

const users = ref<UserModel[] | null>(null)
const categs = ref<Map<string, CategModel[] | null>>(new Map<string, CategModel[] | null>())
const files = ref<Map<string, FileModel[] | null>>(new Map<string, FileModel[] | null>())

const fetched = ref<Set<string>>(new Set<string>())

onBeforeMount(async () => {
  // Obter todos os usuários
  users.value = await getAllUsers()
})

// Função para obter valores das categorias e armazenar em categs
async function handleGetCategories(userId: string) {
  if (fetched.value.has(userId)) {
    return
  }

  const data: CategModel[] | null = await getAllCategories(userId)
  categs.value.set(userId, data)
  fetched.value.add(userId)
}

// Função para obter valores dos arquivos e armazenar em files
async function handleGetFiles(userId: string, categId: string) {
  if (fetched.value.has(categId)) {
    return
  }

  const data: FileModel[] | null = await getAllFiles(userId, categId)
  files.value.set(userId, data)
  fetched.value.add(categId)
}
</script>

<template>
  <Header />
  <main
    class="flex w-screen flex-col items-center justify-center gap-8 bg-gray bg-opacity-5 px-8 py-4 font-lato sm:px-16 sm:py-8"
  >
    <!-- Seção de introdução -->
    <section class="w-fit max-w-3xl rounded-lg bg-white px-10 py-8 text-center shadow-xl drop-shadow-xl sm:text-lg">
      <h2 class="mb-4 text-2xl">Página de Administrador</h2>
      <p>
        Lorem ipsum dolor sit amet, consectetur adipiscing elit. Integer mauris est, fringilla a malesuada vel, molestie
        pretium mi. Pellentesque habitant morbi tristique senectus et netus et malesuada fames ac turpis egestas.
      </p>
    </section>
    <!-- TODO: Accordion de usuários -->
    <AccordionComponent v-for="user in users" :key="user.user_id">
      <template #header>
        <div :onmouseover="() => handleGetCategories(user.user_id)">{{ user.name }}</div>
      </template>
      <template #content>
        <!-- TODO: Accordion de categorias -->
        <AccordionComponent v-for="categ in categs.get(user.user_id)" :key="categ.categ_id">
          <template #header>
            <div :onmouseover="() => handleGetFiles(user.user_id, categ.categ_id)">{{ categ.name }}</div>
          </template>
          <template #content>
            <!-- TODO: Accordion de arquivos -->
            <AccordionComponent v-for="file in files.get(categ.categ_id)" :key="file.file_id">
              <template #header>{{ file.name }}</template>
              <template #content></template>
            </AccordionComponent>
            <!-- TODO: Botão de criar arquivo -->
          </template>
        </AccordionComponent>
        <!-- TODO: Botão de criar categoria -->
      </template>
    </AccordionComponent>
    <!-- TODO: Botão de criar usuário -->
  </main>
</template>

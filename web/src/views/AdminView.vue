<script setup lang="ts">
import Header from '@/components/generic/HeaderSection.vue'
import AccordionComponent from '@/components/generic/AccordionComponent.vue'
import { onBeforeMount, ref } from 'vue'
import type { CategModel, FileModel, UserModel } from '@/@types/Responses.ts'
import { getAllCategories, getAllFiles, getAllUsers } from '@/services/queries.ts'
import AddButton from '@/components/admin/AddButton.vue'

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

// Função para abrir janela de atualização de usuário
async function handleUpdateUser(userId: string) {
  console.log('update user ' + userId)
}

// Função para abrir janela de atualização de categoria
async function handleUpdateCategory(userId: string, categId: string) {
  console.log('update category ' + userId + ' ' + categId)
}

// Função para abrir janela de exclusão de usuário
async function handleDeleteUser(userId: string) {
  console.log('delete user ' + userId)
}

// Função para abrir janela de exclusão de categoria
async function handleDeleteCategory(userId: string, categId: string) {
  console.log('delete category ' + userId + ' ' + categId)
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
    <section class="flex w-full max-w-3xl flex-col items-center justify-center">
      <!-- Accordion de usuários -->
      <AccordionComponent
        v-for="(user, u_index) in users"
        class="z-20"
        :key="user.user_id"
        :title="user.name"
        :admin="true"
        :first="u_index === 0"
        :last="!!users && u_index == users.length - 1"
        :edit-handler="() => handleUpdateUser(user.user_id)"
        :delete-handler="() => handleDeleteUser(user.user_id)"
        :onmouseover="() => handleGetCategories(user.user_id)"
      >
        <template #content>
          <!-- Accordion de categorias -->
          <AccordionComponent
            v-for="(categ, c_index) in categs.get(user.user_id)"
            class="z-10"
            :key="categ.categ_id"
            :title="categ.name"
            :admin="true"
            :first="c_index === 0"
            :last="!!users && c_index == users.length - 1"
            :edit-handler="() => handleUpdateCategory(user.user_id, categ.categ_id)"
            :delete-handler="() => handleDeleteCategory(user.user_id, categ.categ_id)"
            :onmouseover="() => handleGetFiles(user.user_id, categ.categ_id)"
          >
            <template #content>
              <!-- TODO: Container de arquivos -->
              <div v-for="file in files.get(categ.categ_id)" :key="file.file_id">
                <span>{{ file.name }}</span>
              </div>
              <!-- Botão de criar arquivo -->
              <AddButton text="Novo arquivo" />
            </template>
          </AccordionComponent>
          <!-- Botão de criar categoria -->
          <AddButton text="Nova categoria" />
        </template>
      </AccordionComponent>
    </section>
    <!-- Botão de criar usuário -->
    <AddButton text="Novo usuário" class="w-full max-w-3xl" />
  </main>
</template>

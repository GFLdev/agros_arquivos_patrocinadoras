<script setup lang="ts">
import Header from '@/components/generic/HeaderSection.vue'
import { PhEnvelopeSimple, PhMapPin, PhPhone } from '@phosphor-icons/vue'
import AccordionComponent from '@/components/generic/AccordionComponent.vue'
import { onBeforeMount, ref } from 'vue'
import type { CategModel, FileModel, UserModel } from '@/@types/Responses.ts'
import { getAllCategories, getAllFiles, getUserById } from '@/services/queries.ts'
import { useAuthStore } from '@/stores/authStore.ts'
import FileItem from '@/components/generic/FileItem.vue'
import { downloadFile } from '@/utils/file.ts'

const authStore = useAuthStore()

const user = ref<UserModel>({
  user_id: '',
  username: '',
  name: '',
  password: '',
  updated_at: '',
})
const categs = ref<CategModel[]>([])
const files = ref<Record<string, FileModel[]>>({})

const fetched = ref<Set<string>>(new Set<string>())

onBeforeMount(() => {
  // Obter dados do usuário
  getUserById(authStore.user?.id ?? '').then((res) => {
    user.value = res.data ?? {
      user_id: '',
      username: '',
      name: '',
      password: '',
      updated_at: '',
    }
  })
  // Obter todas as categorias
  getAllCategories(authStore.user?.id ?? '').then((res) => {
    categs.value = res.data ?? []
  })
})

// Função para obter valores dos arquivos e armazenar em files
function handleGetFiles(categId: string) {
  if (fetched.value.has(categId) || !authStore.user) {
    return
  }
  getAllFiles(authStore.user.id, categId).then((res) => {
    files.value[categId] = res.data ?? []
    fetched.value.add(categId)
  })
}
</script>

<template>
  <Header :title="user?.name" />
  <main
    class="flex w-full flex-col items-center justify-center gap-12 bg-gray bg-opacity-5 px-8 py-4 font-lato sm:px-16 sm:py-8"
  >
    <section class="w-fit max-w-3xl rounded-lg bg-white px-10 py-8 text-center shadow-xl drop-shadow-xl sm:text-lg">
      <p>
        Nesta página estão disponíveis informações sobre os planos de previdência e de saúde do Agros, em atendimento à
        Resolução Normativa 389, da Agência Nacional de Saúde Suplementar (ANS).
      </p>
    </section>
    <div class="w-full max-w-3xl">
      <section v-if="categs && categs.length > 0" class="w-full">
        <AccordionComponent
          v-for="(categ, c_index) in categs.sort((a, b) => a.name.localeCompare(b.name))"
          class="z-10"
          :key="categ.categ_id"
          :title="categ.name"
          :admin="false"
          :first="c_index === 0"
          :last="!!categs && c_index === categs.length - 1"
          :onmouseover="() => handleGetFiles(categ.categ_id)"
        >
          <template #content>
            <div v-if="files[categ.categ_id] && files[categ.categ_id].length > 0" class="w-full">
              <FileItem
                v-for="file in files[categ.categ_id].sort((a, b) => a.name.localeCompare(b.name))"
                :key="file.file_id"
                :admin="false"
                :name="file.name"
                :mime-type="file.mimetype"
                :last-modified="file.updated_at"
                :download-handler="async () => await downloadFile(user.user_id, categ.categ_id, file.file_id)"
              />
            </div>
            <p v-else>Nenhum arquivo encontrado</p>
          </template>
        </AccordionComponent>
      </section>
      <p v-else class="w-full text-center">Nenhuma categoria encontrada</p>
    </div>
    <!-- Divisor -->
    <div class="relative w-full max-w-4xl">
      <div class="absolute inset-0 flex items-center" aria-hidden="true">
        <div class="w-full border-t border-gray border-opacity-30" />
      </div>
      <div class="relative flex justify-center">
        <span class="bg-[ text-gray-500 bg-gray px-2 text-sm" />
      </div>
    </div>
    <!-- Seção de dúvidas -->
    <section class="relative rounded-lg bg-white px-10 py-8 shadow-xl drop-shadow-xl lg:static">
      <div class="w-full max-w-xl text-center">
        <!-- Heading -->
        <h2 class="text-pretty text-2xl font-semibold tracking-tight text-dark">Dúvidas?</h2>
        <!-- Descrição -->
        <p class="mt-6 text-dark sm:text-lg/8">
          Em caso de dúvida sobre as informações aqui apresentadas, entre em contato com o Agros.
        </p>
        <!-- Itens -->
        <dl
          class="mt-6 grid grid-cols-1 space-y-4 text-dark sm:grid-cols-3 sm:items-start sm:justify-between sm:space-y-0 sm:text-base/7"
        >
          <div
            class="flex items-center justify-center gap-4 transition-colors duration-200 hover:text-primary sm:flex-col sm:gap-2"
          >
            <dt class="flex-none">
              <span class="sr-only">Endereço</span>
              <PhMapPin class="size-6" aria-hidden="true" />
            </dt>
            <dd>Av. Purdue, s/n<br />Campus da UFV - Viçosa</dd>
          </div>
          <div
            class="flex items-center justify-center gap-4 transition-colors duration-200 hover:text-primary sm:flex-col sm:gap-2"
          >
            <dt class="flex-none">
              <span class="sr-only">Telefone</span>
              <PhPhone class="size-6" aria-hidden="true" />
            </dt>
            <dd>
              <a>(31) 3899-6550<br />ramal 6539</a>
            </dd>
          </div>
          <div
            class="flex items-center justify-center gap-4 transition-colors duration-200 hover:text-primary sm:flex-col sm:gap-2"
          >
            <dt class="flex-none">
              <span class="sr-only">Email</span>
              <PhEnvelopeSimple class="size-6" aria-hidden="true" />
            </dt>
            <dd>
              <a href="mailto:dge@agros.org.br">dge@agros.org.br</a>
            </dd>
          </div>
        </dl>
      </div>
    </section>
  </main>
</template>

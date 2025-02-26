<script setup lang="ts">
import Header from '@/components/generic/HeaderSection.vue'
import { PhEnvelopeSimple, PhMapPin, PhPhone } from '@phosphor-icons/vue'
import AccordionComponent from '@/components/generic/AccordionComponent.vue'
import { onBeforeMount, type Ref, ref } from 'vue'
import type { CategModel, FileModel, GetAllResponse } from '@/@types/Responses.ts'
import { getAllCategories, getAllFiles } from '@/services/queries.ts'
import { useAuthStore } from '@/stores/authStore.ts'
import FileItem from '@/components/generic/FileItem.vue'
import { downloadFile } from '@/utils/file.ts'
import PopupAlert from '@/components/generic/PopupAlert.vue'
import { Alert, codeToAlertType } from '@/utils/modals.ts'
import { AlertType } from '@/@types/Enumerations.ts'

// Pinia Store
const authStore: ReturnType<typeof useAuthStore> = useAuthStore()

// Array das categorias
const categs: Ref<CategModel[]> = ref<CategModel[]>([])

// Mapeamento dos arquivos
const files: Ref<Record<string, FileModel[]>> = ref<Record<string, FileModel[]>>({})

// Set para impedir múltiplas requisições no banco
const fetched: Ref<Set<string>> = ref<Set<string>>(new Set<string>())

// Alerta
const alert: Ref<Alert> = ref<Alert>(new Alert())

/**
 * Lida com a obtenção de dados das categorias e gerencia tentativas em caso de erros.
 *
 * @param {number} [retries=0] - O número atual de tentativas para obter os dados das categorias.
 * @param {number} [maxRetry=5] - O número máximo de tentativas para obter os dados das categorias.
 * @param {number} [timeout=5000] - A duração do tempo limite (em milissegundos) antes de tentar novamente a
 * solicitação.
 * @return {Promise<void>} Uma promise que é resolvida quando os dados das categorias são obtidos com sucesso ou todas
 * as tentativas são esgotadas.
 */
async function handleGetCategories(retries: number = 0, maxRetry: number = 5, timeout: number = 5000): Promise<void> {
  if (retries >= maxRetry) return
  try {
    const res: GetAllResponse<CategModel> = await getAllCategories(authStore.user?.id ?? '')
    if (res.code >= 400) {
      alert.value.handleAlert(res.message, codeToAlertType(res.code))
      setTimeout((): Promise<void> => handleGetCategories(retries + 1, timeout), timeout)
    } else {
      categs.value = res.data ?? []
    }
  } catch {
    alert.value.handleAlert('Erro desconhecido. Tente novamente mais tarde.', AlertType.Error)
    setTimeout((): Promise<void> => handleGetCategories(retries + 1, timeout), timeout)
  }
}

/**
 * Busca e lida com os arquivos de um usuário e categoria específicos.
 *
 * @param {string} categId - O ID da categoria sob a qual os arquivos estão organizados.
 * @return {Promise<void>} Uma promise que é resolvida assim que os arquivos forem processados ou rejeitada caso ocorra
 * um erro durante a busca.
 */
async function handleGetFiles(categId: string): Promise<void> {
  if (fetched.value.has(categId) || !authStore.user) {
    return
  }

  try {
    const res: GetAllResponse<FileModel> = await getAllFiles(authStore.user.id, categId)
    if (res.code >= 400) {
      alert.value.handleAlert(res.message, codeToAlertType(res.code))
    } else {
      files.value[categId] = res.data ?? []
      fetched.value.add(categId)
    }
  } catch {
    alert.value.handleAlert('Erro desconhecido. Tente novamente mais tarde.', AlertType.Error)
  }
}

// Obter todas as categorias
onBeforeMount(handleGetCategories)
</script>

<template>
  <Header :title="authStore.user?.name" />
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
                :download-handler="
                  async () => await downloadFile(authStore.user?.id ?? '', categ.categ_id, file.file_id)
                "
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
  <!-- Alerta -->
  <PopupAlert :text="alert.text" :type="alert.type" :duration="alert.duration" v-model="alert.show" />
</template>

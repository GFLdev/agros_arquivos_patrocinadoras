<script setup lang="ts">
import Header from '@/components/generic/HeaderSection.vue'
import AccordionComponent from '@/components/generic/AccordionComponent.vue'
import { onBeforeMount, type Ref, ref } from 'vue'
import type { CategModel, FileModel, GetAllResponse, UserModel } from '@/@types/Responses.ts'
import { getAllCategories, getAllFiles, getAllUsers } from '@/services/queries.ts'
import AddButton from '@/components/admin/AddButton.vue'
import CreateUserForm from '@/components/admin/CreateUserForm.vue'
import DeleteUserForm from '@/components/admin/DeleteUserForm.vue'
import CreateCategoryForm from '@/components/admin/CreateCategoryForm.vue'
import DeleteCategoryForm from '@/components/admin/DeleteCategoryForm.vue'
import CreateFileForm from '@/components/admin/CreateFileForm.vue'
import FileItem from '@/components/generic/FileItem.vue'
import DeleteFileForm from '@/components/admin/DeleteFileForm.vue'
import UpdateUserForm from '@/components/admin/UpdateUserForm.vue'
import UpdateCategoryForm from '@/components/admin/UpdateCategoryForm.vue'
import { downloadFile } from '@/utils/file.ts'
import UpdateFileForm from '@/components/admin/UpdateFileForm.vue'
import { AlertType } from '@/@types/Enumerations.ts'
import PopupAlert from '@/components/generic/PopupAlert.vue'
import { Alert, codeToAlertType } from '@/utils/modals.ts'

// Dados do banco
const users: Ref<UserModel[]> = ref<UserModel[]>([])
const categs: Ref<Record<string, CategModel[]>> = ref<Record<string, CategModel[]>>({})
const files: Ref<Record<string, FileModel[]>> = ref<Record<string, FileModel[]>>({})

// Set para impedir múltiplas requisições no banco
const fetched: Ref<Set<string>> = ref<Set<string>>(new Set<string>())

// ID's selecionados
const selectedUser: Ref<UserModel | undefined> = ref<UserModel>()
const selectedCateg: Ref<CategModel | undefined> = ref<CategModel>()
const selectedFile: Ref<FileModel | undefined> = ref<FileModel>()

// Mapeamento para o formulário de atualização
const selectedUsersRec: Ref<Record<string, string>> = ref<Record<string, string>>({})
const selectedCategsRec: Ref<Record<string, string>> = ref<Record<string, string>>({})

// Estados para abrir janelas
const showCreateUser: Ref<boolean> = ref<boolean>(false)
const showCreateCateg: Ref<boolean> = ref<boolean>(false)
const showCreateFile: Ref<boolean> = ref<boolean>(false)

const showUpdateUser: Ref<boolean> = ref<boolean>(false)
const showUpdateCateg: Ref<boolean> = ref<boolean>(false)
const showUpdateFile: Ref<boolean> = ref<boolean>(false)

const showDeleteUser: Ref<boolean> = ref<boolean>(false)
const showDeleteCateg: Ref<boolean> = ref<boolean>(false)
const showDeleteFile: Ref<boolean> = ref<boolean>(false)

// Hack para animação dos accordions aninhados
const contentHeight: Ref<number> = ref<number>(0)

// Alerta
const alert: Ref<Alert> = ref<Alert>(new Alert())

/**
 * Lida com a obtenção de dados de usuários e gerencia tentativas em caso de erros.
 *
 * @param {number} [retries=0] - O número atual de tentativas para obter os dados do usuário.
 * @param {number} [maxRetry=5] - O número máximo de tentativas para obter os dados do usuário.
 * @param {number} [timeout=5000] - A duração do tempo limite (em milissegundos) antes de tentar novamente a
 * solicitação.
 * @return {Promise<void>} Uma promise que é resolvida quando os dados do usuário são obtidos com sucesso ou todas as
 * tentativas são esgotadas.
 */
async function handleGetUsers(retries: number = 0, maxRetry: number = 5, timeout: number = 5000): Promise<void> {
  if (retries >= maxRetry) return
  try {
    const res: GetAllResponse<UserModel> = await getAllUsers()
    if (res.code >= 400) {
      alert.value.handleAlert(res.message, codeToAlertType(res.code))
      setTimeout((): Promise<void> => handleGetUsers(retries + 1, timeout), timeout)
    } else {
      users.value = res.data ?? []
    }
  } catch {
    alert.value.handleAlert('Erro desconhecido. Tente novamente mais tarde.', AlertType.Error)
    setTimeout((): Promise<void> => handleGetUsers(retries + 1, timeout), timeout)
  }
}

/**
 * Lida com a recuperação de categorias para um usuário específico.
 *
 * @param {string} userId - O identificador único do usuário para o qual as categorias serão recuperadas.
 * @param {boolean} [reset=false] - Indica se deve redefinir e buscar as categorias, mesmo que elas já tenham sido
 * buscadas.
 * @return {Promise<void>} Uma promise que é resolvida assim que as categorias são recuperadas e armazenadas ou um erro
 * é tratado.
 */
async function handleGetCategories(userId: string, reset: boolean = false): Promise<void> {
  if (!reset && fetched.value.has(userId)) {
    return
  }

  try {
    const res: GetAllResponse<CategModel> = await getAllCategories(userId)
    if (res.code >= 400) {
      alert.value.handleAlert(res.message, codeToAlertType(res.code))
    } else {
      categs.value[userId] = res.data ?? []
      fetched.value.add(userId)
    }
  } catch {
    alert.value.handleAlert('Erro desconhecido. Tente novamente mais tarde.', AlertType.Error)
  }
}

/**
 * Busca e lida com os arquivos de um usuário e categoria específicos.
 *
 * @param {string} userId - O ID do usuário para quem os arquivos estão sendo recuperados.
 * @param {string} categId - O ID da categoria sob a qual os arquivos estão organizados.
 * @param {boolean} [reset=false] - Indica se deve forçar uma nova solicitação para buscar os arquivos, ignorando
 * qualquer dado em cache.
 * @return {Promise<void>} Uma promise que é resolvida assim que os arquivos forem processados ou rejeitada caso ocorra
 * um erro durante a busca.
 */
async function handleGetFiles(userId: string, categId: string, reset: boolean = false): Promise<void> {
  if (!reset && fetched.value.has(categId)) {
    return
  }

  try {
    const res: GetAllResponse<FileModel> = await getAllFiles(userId, categId)
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

// Obter todos os usuários antes da primeira renderização
onBeforeMount(handleGetUsers)
</script>

<template>
  <Header />
  <main
    class="flex h-full w-full animate-open-with-fade flex-col items-center justify-center gap-8 bg-gray bg-opacity-5 px-8 py-4 font-lato sm:px-16 sm:py-8"
  >
    <!-- Seção de introdução -->
    <section class="w-fit max-w-3xl rounded-lg bg-white px-10 py-8 text-center shadow-xl drop-shadow-xl sm:text-lg">
      <h2 class="mb-4 text-2xl font-light">Página de Administrador</h2>
      <p class="mb-4 font-light">
        Bem-vindo à página de administração! Aqui você pode gerenciar usuários, categorias e arquivos do sistema de
        forma prática e centralizada. Utilize as ferramentas abaixo para executar ações como criação, edição e remoção,
        garantindo a organização e eficiência da plataforma.
      </p>
    </section>
    <section class="flex w-full max-w-3xl flex-col items-center justify-center">
      <!-- Accordion de usuários -->
      <AccordionComponent
        v-for="(user, u_index) in users.sort((a, b) => a.name.localeCompare(b.name))"
        class="z-20"
        :content-height="contentHeight"
        :key="user.user_id"
        :title="user.name"
        :admin="true"
        :first="u_index === 0"
        :last="u_index === users.length - 1"
        :edit-handler="
          () => {
            selectedUser = user
            showUpdateUser = true
          }
        "
        :delete-handler="
          () => {
            selectedUser = user
            showDeleteUser = true
          }
        "
        :onmouseover="() => handleGetCategories(user.user_id)"
      >
        <template #content>
          <div v-if="categs[user.user_id] && categs[user.user_id].length > 0" class="w-full">
            <!-- Accordion de categorias -->
            <AccordionComponent
              v-for="(categ, c_index) in categs[user.user_id].sort((a, b) => a.name.localeCompare(b.name))"
              class="z-10"
              v-model="contentHeight"
              :key="categ.categ_id"
              :title="categ.name"
              :admin="true"
              :first="c_index === 0"
              :last="c_index === categs[user.user_id].length - 1"
              :edit-handler="
                () => {
                  selectedUser = user
                  selectedCateg = categ
                  users.map((u) => (selectedUsersRec[u.user_id] = u.name))
                  showUpdateCateg = true
                }
              "
              :delete-handler="
                () => {
                  selectedUser = user
                  selectedCateg = categ
                  showDeleteCateg = true
                }
              "
              :onmouseover="() => handleGetFiles(user.user_id, categ.categ_id)"
            >
              <template #content>
                <!-- Container de arquivos -->
                <div v-if="files[categ.categ_id] && files[categ.categ_id].length > 0" class="w-full">
                  <FileItem
                    v-for="file in files[categ.categ_id].sort((a, b) => a.name.localeCompare(b.name))"
                    :key="file.file_id"
                    :admin="true"
                    :name="file.name"
                    :mime-type="file.mimetype"
                    :last-modified="file.updated_at"
                    :edit-handler="
                      () => {
                        selectedUser = user
                        selectedCateg = categ
                        selectedFile = file
                        categs[user.user_id].map((c) => (selectedCategsRec[c.categ_id] = c.name))
                        showUpdateFile = true
                      }
                    "
                    :delete-handler="
                      () => {
                        selectedUser = user
                        selectedCateg = categ
                        selectedFile = file
                        showDeleteFile = true
                      }
                    "
                    :download-handler="() => downloadFile(user.user_id, categ.categ_id, file.file_id)"
                  />
                </div>
                <!-- Botão de criar arquivo -->
                <AddButton
                  text="Novo arquivo"
                  :on-click="
                    () => {
                      selectedUser = user
                      selectedCateg = categ
                      showCreateFile = true
                    }
                  "
                />
              </template>
            </AccordionComponent>
          </div>
          <!-- Botão de criar categoria -->
          <AddButton
            text="Nova categoria"
            :on-click="
              () => {
                selectedUser = user
                showCreateCateg = true
              }
            "
          />
        </template>
      </AccordionComponent>
    </section>
    <!-- Botão de criar usuário -->
    <AddButton text="Novo usuário" :on-click="() => (showCreateUser = true)" class="w-full max-w-3xl" />
  </main>
  <!-- Janelas de usuário -->
  <CreateUserForm v-model="showCreateUser" @submitted="handleGetUsers" />
  <UpdateUserForm v-if="selectedUser" v-model="showUpdateUser" :user="selectedUser" @submitted="handleGetUsers" />
  <DeleteUserForm v-if="selectedUser" v-model="showDeleteUser" :user="selectedUser" @submitted="handleGetUsers" />
  <!-- Janelas de categoria -->
  <CreateCategoryForm
    v-if="selectedUser"
    v-model="showCreateCateg"
    :user="selectedUser"
    @submitted="() => handleGetCategories(selectedUser?.user_id ?? '', true)"
  />
  <UpdateCategoryForm
    v-if="selectedCateg"
    v-model="showUpdateCateg"
    :categ="selectedCateg"
    :users="selectedUsersRec"
    @submitted="() => users.map((u) => handleGetCategories(u.user_id, true))"
  />
  <DeleteCategoryForm
    v-if="selectedCateg"
    v-model="showDeleteCateg"
    :categ="selectedCateg"
    @submitted="() => handleGetCategories(selectedUser?.user_id ?? '', true)"
  />
  <!-- Janelas de arquivo -->
  <CreateFileForm
    v-if="selectedCateg"
    v-model="showCreateFile"
    :categ="selectedCateg"
    @submitted="() => handleGetFiles(selectedCateg?.user_id ?? '', selectedCateg?.categ_id ?? '', true)"
  />
  <UpdateFileForm
    v-if="selectedCateg && selectedFile"
    v-model="showUpdateFile"
    :categ="selectedCateg"
    :file="selectedFile"
    :categs="selectedCategsRec"
    @submitted="() => categs[selectedCateg?.user_id ?? ''].map((c) => handleGetFiles(c.user_id, c.categ_id, true))"
  />
  <DeleteFileForm
    v-if="selectedCateg && selectedFile"
    v-model="showDeleteFile"
    :categ="selectedCateg"
    :file="selectedFile"
    @submitted="() => handleGetFiles(selectedCateg?.user_id ?? '', selectedCateg?.categ_id ?? '', true)"
  />
  <!-- Alerta -->
  <PopupAlert :text="alert.text" :type="alert.type" :duration="alert.duration" v-model="alert.show" />
</template>

<script setup lang="ts">
import Header from '@/components/generic/HeaderSection.vue'
import AccordionComponent from '@/components/generic/AccordionComponent.vue'
import { onBeforeMount, ref } from 'vue'
import type { CategModel, FileModel, UserModel } from '@/@types/Responses.ts'
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
import { codeToAlertType } from '@/utils/modals.ts'

// Dados do banco
const users = ref<UserModel[]>([])
const categs = ref<Record<string, CategModel[]>>({})
const files = ref<Record<string, FileModel[]>>({})

// Set para impedir múltiplas requisições no banco
const fetched = ref<Set<string>>(new Set<string>())

// ID's selecionados
const selectedUser = ref<UserModel>({
  user_id: '',
  username: '',
  name: '',
  password: '',
  updated_at: '',
})
const selectedCateg = ref<CategModel>({
  categ_id: '',
  user_id: '',
  name: '',
  updated_at: '',
})
const selectedFile = ref<FileModel>({
  file_id: '',
  categ_id: '',
  name: '',
  extension: '',
  mimetype: '',
  blob: '',
  updated_at: 0,
})

// Mapeamento para o formulário de atualização
const selectedUsersMap = ref<Map<string, string>>(new Map<string, string>())
const selectedCategsMap = ref<Map<string, string>>(new Map<string, string>())

// Estados para abrir janelas
const showCreateUser = ref<boolean>(false)
const showCreateCateg = ref<boolean>(false)
const showCreateFile = ref<boolean>(false)

const showUpdateUser = ref<boolean>(false)
const showUpdateCateg = ref<boolean>(false)
const showUpdateFile = ref<boolean>(false)

const showDeleteUser = ref<boolean>(false)
const showDeleteCateg = ref<boolean>(false)
const showDeleteFile = ref<boolean>(false)

// Hack para animação dos accordions aninhados
const contentHeight = ref<number>(0)

// Alerta
const showAlert = ref<boolean>(false)
const alertType = ref<AlertType>(AlertType.Info)
const alertText = ref<string>('')
const alertDuration = ref<number>(3000)

// Gerenciar alerta
function handleAlert(text: string, type: AlertType = AlertType.Info, duration: number = 3000) {
  alertText.value = text
  alertType.value = type
  alertDuration.value = duration
  showAlert.value = true
}

// Função para obter valores dos usuários e armazenar em users
async function handleGetUsers() {
  try {
    const res = await getAllUsers()
    users.value = res.data ?? []
    if (res.code >= 400) {
      handleAlert(res.message, codeToAlertType(res.code))
    }
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
  } catch (_) {
    handleAlert('Erro desconhecido. Tente novamente mais tarde.', AlertType.Error)
  }
}

// Função para obter valores das categorias e armazenar em categs
async function handleGetCategories(userId: string, reset = false) {
  if (!reset && fetched.value.has(userId)) {
    return
  }

  try {
    const res = await getAllCategories(userId)
    categs.value[userId] = res.data ?? []
    fetched.value.add(userId)
    if (res.code >= 400) {
      handleAlert(res.message, codeToAlertType(res.code))
    }
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
  } catch (_) {
    handleAlert('Erro desconhecido. Tente novamente mais tarde.', AlertType.Error)
  }
}

// Função para obter valores dos arquivos e armazenar em files
async function handleGetFiles(userId: string, categId: string, reset = false) {
  if (!reset && fetched.value.has(categId)) {
    return
  }

  try {
    const res = await getAllFiles(userId, categId)
    files.value[categId] = res.data ?? []
    fetched.value.add(categId)
    if (res.code >= 400) {
      handleAlert(res.message, codeToAlertType(res.code))
    }
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
  } catch (_) {
    handleAlert('Erro desconhecido. Tente novamente mais tarde.', AlertType.Error)
  }
}

onBeforeMount(() => {
  // Obter todos os usuários
  handleGetUsers()
})
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
                  users.map((u) => selectedUsersMap.set(u.user_id, u.name))
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
                        categs[user.user_id].map((c) => selectedCategsMap.set(c.categ_id, c.name))
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
  <UpdateUserForm v-model="showUpdateUser" :user="selectedUser" @submitted="handleGetUsers" />
  <DeleteUserForm v-model="showDeleteUser" :user="selectedUser" @submitted="handleGetUsers" />
  <!-- Janelas de categoria -->
  <CreateCategoryForm
    v-model="showCreateCateg"
    :user="selectedUser"
    @submitted="() => handleGetCategories(selectedUser.user_id, true)"
  />
  <UpdateCategoryForm
    v-model="showUpdateCateg"
    :categ="selectedCateg"
    :users="selectedUsersMap"
    @submitted="() => users.map((u) => handleGetCategories(u.user_id, true))"
  />
  <DeleteCategoryForm
    v-model="showDeleteCateg"
    :categ="selectedCateg"
    @submitted="() => handleGetCategories(selectedUser.user_id, true)"
  />
  <!-- Janelas de arquivo -->
  <CreateFileForm
    v-model="showCreateFile"
    :categ="selectedCateg"
    @submitted="() => handleGetFiles(selectedCateg.user_id, selectedCateg.categ_id, true)"
  />
  <UpdateFileForm
    v-model="showUpdateFile"
    :categ="selectedCateg"
    :file="selectedFile"
    :categs="selectedCategsMap"
    @submitted="() => categs[selectedCateg.user_id].map((c) => handleGetFiles(c.user_id, c.categ_id, true))"
  />
  <DeleteFileForm
    v-model="showDeleteFile"
    :categ="selectedCateg"
    :file="selectedFile"
    @submitted="() => handleGetFiles(selectedCateg.user_id, selectedCateg.categ_id, true)"
  />
  <!-- Alerta -->
  <PopupAlert :text="alertText" :type="alertType" :duration="alertDuration" v-model="showAlert" />
</template>

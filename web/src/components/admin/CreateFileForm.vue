<script setup lang="ts">
import { PhFile, PhFilePlus, PhFolderOpen, PhXCircle } from '@phosphor-icons/vue'
import InputText from '@/components/generic/InputText.vue'
import SubmitButton from '@/components/generic/SubmitButton.vue'
import CancelButton from '@/components/generic/CancelButton.vue'
import PopupWindow from '@/components/generic/PopupWindow.vue'
import type { FileRequest } from '@/@types/Requests.ts'
import { type PropType, ref, watch, watchEffect } from 'vue'
import InputFile from '@/components/generic/InputFile.vue'
import { fileIcon, toBase64 } from '@/utils/file.ts'
import { createFile } from '@/services/queries.ts'
import type { CategModel } from '@/@types/Responses.ts'
import { AlertType } from '@/@types/Enumerations.ts'
import PopupAlert from '@/components/generic/PopupAlert.vue'
import { codeToAlertType } from '@/utils/modals.ts'

defineProps({
  categ: {
    type: Object as PropType<CategModel>,
    required: true,
  },
})

const emits = defineEmits(['submitted'])

// Formulário
const name = ref<string>('')
const extension = ref<string>('')
const inFile = ref<File>(new File([], ''))

// Validações
const loading = ref<boolean>(false)
const filled = ref<boolean>(false)
const formValid = ref<boolean>(false)

// Estado de visibilidade
const showModel = defineModel<boolean>()

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

// Função para abrir janela de criação de arquivo
async function handleCreateFile(userId: string, categId: string) {
  if (!formValid.value) {
    handleAlert('Campos necessários não preenchidos', AlertType.Warning)
    return
  }

  extension.value = '.' + (inFile.value.name.split('.').pop() ?? 'bin')
  loading.value = true

  const body: FileRequest = {
    name: name.value,
    extension: extension.value,
    mimetype: inFile.value.type,
    content: (await toBase64(inFile.value)).split(',')[1], // apenas os bytes base64
  }

  try {
    const res = await createFile(userId, categId, body)
    handleAlert(res.message, codeToAlertType(res.code))
    emits('submitted')
    showModel.value = false
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
  } catch (_) {
    handleAlert('Erro desconhecido. Tente novamente mais tarde.', AlertType.Error)
  } finally {
    loading.value = false
  }
}

// Função para verificar se 'file' está vazio
function isFileEmpty(): boolean {
  return inFile.value.size === 0 && inFile.value.name === ''
}

// Função para reset das variáveis reativas
function reset() {
  loading.value = false
  name.value = ''
  extension.value = ''
  inFile.value = new File([], '')
}

watchEffect(() => {
  const n = name.value.length > 0
  const f = !isFileEmpty()
  filled.value = n && f
  formValid.value = filled.value
})

watch(
  () => showModel.value,
  () => {
    if (!showModel.value) {
      // Timeout para homogeneidade com as transições
      setTimeout(reset, 200)
    }
  },
)
</script>

<template>
  <PopupWindow :title="`Criar novo arquivo`" v-model="showModel">
    <form
      class="flex flex-col gap-4 space-y-4 px-8 py-4"
      @submit.prevent="() => handleCreateFile(categ.user_id, categ.categ_id)"
    >
      <div class="flex w-full flex-col gap-4">
        <!-- Campos do formulário -->
        <InputText
          placeholder="Nome do arquivo"
          label="Nome"
          v-model="name"
          :left-inner-icon="PhFile"
          :required="true"
        />
        <InputFile
          placeholder="Selecione o arquivo"
          label="Arquivo"
          v-model="inFile"
          :left-inner-icon="isFileEmpty() ? PhFolderOpen : fileIcon(inFile.type)"
          :required="true"
        />
        <!-- Avisos -->
        <div v-if="!formValid" class="w-full text-center text-sm font-light">
          <p v-if="!filled">Preencha todos os campos</p>
        </div>
      </div>
      <div class="flex w-full flex-row items-center justify-end gap-4">
        <!-- Botões -->
        <CancelButton
          text="Cancelar"
          :on-click="() => (showModel = false)"
          :disabled="loading"
          :left-inner-icon="PhXCircle"
        />
        <SubmitButton
          text="Criar"
          loading-text="Criando"
          :loading="loading"
          :disabled="loading || !formValid"
          :left-inner-icon="PhFilePlus"
        />
      </div>
    </form>
  </PopupWindow>
  <PopupAlert :text="alertText" :type="alertType" :duration="alertDuration" v-model="showAlert" />
</template>

<style scoped></style>

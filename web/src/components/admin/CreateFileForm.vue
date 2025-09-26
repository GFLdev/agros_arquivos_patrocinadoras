<script setup lang="ts">
import { PhFile, PhFilePlus, PhFolderOpen, PhXCircle } from '@phosphor-icons/vue'
import InputText from '@/components/generic/InputText.vue'
import SubmitButton from '@/components/generic/SubmitButton.vue'
import CancelButton from '@/components/generic/CancelButton.vue'
import PopupWindow from '@/components/generic/PopupWindow.vue'
import type { FileRequest } from '@/@types/Requests.ts'
import { type EmitFn, type ModelRef, type PropType, type Ref, ref, watch, watchEffect } from 'vue'
import InputFile from '@/components/generic/InputFile.vue'
import { fileIcon, toBase64 } from '@/utils/file.ts'
import { createFile } from '@/services/queries.ts'
import type { CategModel, QueryResponse } from '@/@types/Responses.ts'
import { AlertType } from '@/@types/Enumerations.ts'
import PopupAlert from '@/components/generic/PopupAlert.vue'
import { Alert, codeToAlertType, TRANSITION_DURATION } from '@/utils/modals.ts'
import { isFileEmpty } from '@/utils/validate.ts'

defineProps({
  categ: {
    type: Object as PropType<CategModel>,
    required: true,
  },
})

// Emissores
const emits: EmitFn<'submitted'[]> = defineEmits(['submitted'])

// Formulário
const formName: Ref<string> = ref<string>('')
const formExtension: Ref<string> = ref<string>('')
const formFile: Ref<File> = ref<File>(new File([], ''))

// Status de carregamento
const isLoading: Ref<boolean> = ref<boolean>(false)

// Validações
const isFilled: Ref<boolean> = ref<boolean>(false)
const isFormValid: Ref<boolean> = ref<boolean>(false)

// Alerta
const alert: Ref<Alert> = ref<Alert>(new Alert())

// Estado de visibilidade
const showModel: ModelRef<boolean | undefined> = defineModel<boolean>()

/**
 * Redefine o estado das variáveis para seus valores iniciais.
 *
 * @return {void} Sem valor de retorno.
 */
function reset(): void {
  isLoading.value = false
  formName.value = ''
  formExtension.value = ''
  formFile.value = new File([], '')
}

/**
 * Gerencia a criação de um novo arquivo validando a entrada, preparando os dados do arquivo e fazendo uma solicitação à
 * API.
 *
 * @param {string} userId - O ID do usuário que está criando o arquivo.
 * @param {string} categId - O ID da categoria onde o arquivo será criado.
 * @return {Promise<void>} Uma promessa que é resolvida quando o processo de criação do arquivo é concluído.
 */
async function handleCreateFile(userId: string, categId: string): Promise<void> {
  isLoading.value = true
  if (!isFormValid.value) {
    alert.value.handleAlert('Campos necessários não preenchidos', AlertType.Warning)
    return
  }

  formExtension.value = '.' + (formFile.value.name.split('.').pop() ?? 'bin')

  const body: FileRequest = {
    name: formName.value,
    extension: formExtension.value,
    mimetype: formFile.value.type.trim() !== '' ? formFile.value.type : 'application/octet-stream',
    content: await toBase64(formFile.value),
  }

  try {
    const res: QueryResponse = await createFile(userId, categId, body)
    alert.value.handleAlert(res.message, codeToAlertType(res.code))
    emits('submitted')
    showModel.value = false
  } catch {
    alert.value.handleAlert('Erro desconhecido. Tente novamente mais tarde.', AlertType.Error)
  } finally {
    isLoading.value = false
  }
}

// Timeout para homogeneidade com as transições
watch(
  (): boolean | undefined => showModel.value,
  (): void => {
    if (!showModel.value) setTimeout(reset, TRANSITION_DURATION)
  },
)

// Validações
watchEffect((): void => {
  // Verificar preenchimento
  const n: boolean = formName.value.length > 0
  const f: boolean = !isFileEmpty(formFile.value)
  isFilled.value = n && f

  isFormValid.value = isFilled.value
})
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
          v-model="formName"
          :left-inner-icon="PhFile"
          :required="true"
        />
        <InputFile
          placeholder="Selecione o arquivo"
          label="Arquivo"
          v-model="formFile"
          :left-inner-icon="isFileEmpty(formFile) ? PhFolderOpen : fileIcon(formFile.type)"
          :required="true"
        />
        <!-- Avisos -->
        <div v-if="!isFormValid" class="w-full text-center text-sm font-light">
          <p v-if="!isFilled">Preencha todos os campos</p>
        </div>
      </div>
      <div class="flex w-full flex-row items-center justify-end gap-4">
        <!-- Botões -->
        <CancelButton
          text="Cancelar"
          :on-click="() => (showModel = false)"
          :disabled="isLoading"
          :left-inner-icon="PhXCircle"
        />
        <SubmitButton
          text="Criar"
          loading-text="Criando"
          :loading="isLoading"
          :disabled="isLoading || !isFormValid"
          :left-inner-icon="PhFilePlus"
        />
      </div>
    </form>
  </PopupWindow>
  <PopupAlert :text="alert.text" :type="alert.type" :duration="alert.duration" v-model="alert.show" />
</template>

<style scoped></style>

<script setup lang="ts">
import { PhFile, PhFolder, PhFolderOpen, PhPencil, PhXCircle } from '@phosphor-icons/vue'
import InputText from '@/components/generic/InputText.vue'
import SubmitButton from '@/components/generic/SubmitButton.vue'
import CancelButton from '@/components/generic/CancelButton.vue'
import PopupWindow from '@/components/generic/PopupWindow.vue'
import type { UpdateFileRequest } from '@/@types/Requests.ts'
import { computed, type EmitFn, type ModelRef, onMounted, type PropType, type Ref, ref, watch, watchEffect } from 'vue'
import type { CategModel, FileModel, QueryResponse } from '@/@types/Responses.ts'
import InputList from '@/components/generic/InputList.vue'
import { updateFile } from '@/services/queries.ts'
import InputFile from '@/components/generic/InputFile.vue'
import { fileIcon, toBase64 } from '@/utils/file.ts'
import PopupAlert from '@/components/generic/PopupAlert.vue'
import { AlertType } from '@/@types/Enumerations.ts'
import { Alert, codeToAlertType, TRANSITION_DURATION } from '@/utils/modals.ts'
import { isFileEmpty } from '@/utils/validate.ts'

const props: {
  readonly categ: CategModel
  readonly file: FileModel
  readonly categs: Record<string, string>
} = defineProps({
  categ: {
    type: Object as PropType<CategModel>,
    required: true,
  },
  file: {
    type: Object as PropType<FileModel>,
    required: true,
  },
  categs: {
    type: Object as PropType<Record<string, string>>,
    required: true,
  },
})

// Emissores
const emits: EmitFn<'submitted'[]> = defineEmits(['submitted'])

// Formulário
const formCateg: Ref<string> = ref<string>(props.categ.categ_id)
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
  formCateg.value = computed(() => props.categ.categ_id).value
}

/**
 * Lida com a atualização de um arquivo, validando as entradas do formulário, processando os dados do arquivo,
 * fazendo uma solicitação à API.
 *
 * @param {string} userId - O identificador do usuário que está realizando a atualização.
 * @param {string} categId - O identificador da categoria à qual o arquivo pertence.
 * @param {string} fileId - O identificador do arquivo a ser atualizado.
 * @return {Promise<void>} Uma promessa que é resolvida quando a operação de atualização for concluída.
 */
async function handleUpdateFile(userId: string, categId: string, fileId: string): Promise<void> {
  isLoading.value = true
  if (!isFormValid.value) {
    alert.value.handleAlert('Campos necessários não preenchidos', AlertType.Warning)
    return
  }

  let mimetype: string = ''
  if (!isFileEmpty(formFile.value)) {
    formExtension.value = '.' + (formFile.value.name.split('.').pop() ?? 'bin')
    mimetype = formFile.value.type.trim() !== '' ? formFile.value.type : 'application/octet-stream'
  }

  const body: Partial<UpdateFileRequest> = {
    categ_id: formCateg.value,
    name: formName.value,
    extension: formExtension.value,
    mimetype: mimetype,
    content: await toBase64(formFile.value),
  }

  try {
    const res: QueryResponse = await updateFile(userId, categId, fileId, body)
    alert.value.handleAlert(res.message, codeToAlertType(res.code))
    emits('submitted')
    showModel.value = false
  } catch {
    alert.value.handleAlert('Erro desconhecido. Tente novamente mais tarde.', AlertType.Error)
  } finally {
    isLoading.value = false
  }
}

// Selecionar a categoria atual na lista suspensa, na primeira renderização
onMounted(() => {
  formCateg.value = props.categ.categ_id
})

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
  const c = formCateg.value !== props.categ.categ_id
  const n = formName.value.length > 0
  const f = !isFileEmpty(formFile.value)
  isFilled.value = c || n || f

  isFormValid.value = isFilled.value
})
</script>

<template>
  <PopupWindow :title="`Editar arquivo`" v-model="showModel">
    <form
      class="flex flex-col gap-4 space-y-4 px-8 py-4"
      @submit.prevent="() => handleUpdateFile(categ.user_id, file.categ_id, file.file_id)"
    >
      <div class="flex w-full flex-col gap-4">
        <!-- Campos do formulário -->
        <InputList
          :values="categs"
          label="Categoria"
          v-model="formCateg"
          :selected="categ.categ_id"
          :left-inner-icon="PhFolder"
          :required="false"
        />
        <InputText
          :placeholder="file.name"
          label="Nome"
          v-model="formName"
          :left-inner-icon="PhFile"
          :required="false"
        />
        <InputFile
          placeholder="Selecione o arquivo"
          label="Arquivo"
          v-model="formFile"
          :left-inner-icon="isFileEmpty(formFile) ? PhFolderOpen : fileIcon(formFile.type)"
          :required="false"
        />
        <!-- Avisos -->
        <div v-if="!isFormValid" class="w-full text-center text-sm font-light">
          <p v-if="!isFilled">Preencha ao menos um campo</p>
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
          text="Editar"
          loading-text="Editando"
          :loading="isLoading"
          :disabled="isLoading || !isFormValid"
          :left-inner-icon="PhPencil"
        />
      </div>
    </form>
  </PopupWindow>
  <PopupAlert :text="alert.text" :type="alert.type" :duration="alert.duration" v-model="alert.show" />
</template>

<style scoped></style>

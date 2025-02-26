<script setup lang="ts">
import { PhFolder, PhFolderPlus, PhXCircle } from '@phosphor-icons/vue'
import InputText from '@/components/generic/InputText.vue'
import SubmitButton from '@/components/generic/SubmitButton.vue'
import CancelButton from '@/components/generic/CancelButton.vue'
import PopupWindow from '@/components/generic/PopupWindow.vue'
import type { CategRequest } from '@/@types/Requests.ts'
import { type EmitFn, type ModelRef, type PropType, type Ref, ref, watch, watchEffect } from 'vue'
import { createCategory } from '@/services/queries.ts'
import type { QueryResponse, UserModel } from '@/@types/Responses.ts'
import { AlertType } from '@/@types/Enumerations.ts'
import { Alert, codeToAlertType, TRANSITION_DURATION } from '@/utils/modals.ts'
import PopupAlert from '@/components/generic/PopupAlert.vue'

defineProps({
  user: {
    type: Object as PropType<UserModel>,
    required: true,
  },
})

// Emissores
const emits: EmitFn<'submitted'[]> = defineEmits(['submitted'])

// Formulário
const formName: Ref<string> = ref<string>('')

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
}

/**
 * Lida com a criação de uma categoria enviando os dados fornecidos para o servidor.
 *
 * @param {string} userId - O ID do usuário que está criando a categoria.
 * @return {Promise<void>} Uma Promise que é resolvida quando o processo de criação é concluído.
 */
async function handleCreateCategory(userId: string): Promise<void> {
  isLoading.value = true
  if (!isFormValid.value) {
    alert.value.handleAlert('Campos necessários não preenchidos', AlertType.Warning)
    return
  }

  const body: CategRequest = {
    name: formName.value,
  }

  try {
    const res: QueryResponse = await createCategory(userId, body)
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
  isFilled.value = (formName.value ?? '').length > 0

  isFormValid.value = isFilled.value
})
</script>

<template>
  <PopupWindow :title="`Criar nova categoria`" v-model="showModel">
    <form class="flex flex-col gap-4 space-y-4 px-8 py-4" @submit.prevent="() => handleCreateCategory(user.user_id)">
      <div class="flex w-full flex-col gap-4">
        <!-- Campos do formulário -->
        <InputText
          placeholder="Nome da categoria"
          label="Nome"
          v-model="formName"
          :left-inner-icon="PhFolder"
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
          :left-inner-icon="PhFolderPlus"
        />
      </div>
    </form>
  </PopupWindow>
  <PopupAlert :text="alert.text" :type="alert.type" :duration="alert.duration" v-model="alert.show" />
</template>

<style scoped></style>

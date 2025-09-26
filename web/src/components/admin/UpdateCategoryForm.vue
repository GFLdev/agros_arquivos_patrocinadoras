<script setup lang="ts">
import { PhFolder, PhPencil, PhUser, PhXCircle } from '@phosphor-icons/vue'
import InputText from '@/components/generic/InputText.vue'
import SubmitButton from '@/components/generic/SubmitButton.vue'
import CancelButton from '@/components/generic/CancelButton.vue'
import PopupWindow from '@/components/generic/PopupWindow.vue'
import type { UpdateCategRequest } from '@/@types/Requests.ts'
import { computed, type EmitFn, type ModelRef, onMounted, type PropType, type Ref, ref, watch, watchEffect } from 'vue'
import type { CategModel, QueryResponse } from '@/@types/Responses.ts'
import InputList from '@/components/generic/InputList.vue'
import { updateCategory } from '@/services/queries.ts'
import PopupAlert from '@/components/generic/PopupAlert.vue'
import { AlertType } from '@/@types/Enumerations.ts'
import { Alert, codeToAlertType, TRANSITION_DURATION } from '@/utils/modals.ts'

const props: {
  readonly categ: CategModel
  readonly users: Record<string, string>
} = defineProps({
  categ: {
    type: Object as PropType<CategModel>,
    required: true,
  },
  users: {
    type: Object as PropType<Record<string, string>>,
    required: true,
  },
})

// Emissores
const emits: EmitFn<'submitted'[]> = defineEmits(['submitted'])

// Formulário
const formUser: Ref<string> = ref<string>(props.categ.user_id)
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
  formUser.value = computed(() => props.categ.user_id).value
}

/**
 * Lida com a atualização de uma categoria enviando as informações atualizadas para o servidor
 * e gerenciando a validação do formulário.
 *
 * @param {string} userId - O ID do usuário que está realizando a atualização.
 * @param {string} categId - O ID da categoria a ser atualizada.
 * @return {Promise<void>} Uma promessa que é resolvida quando o processo de atualização da categoria é concluído.
 */
async function handleUpdateCategory(userId: string, categId: string): Promise<void> {
  isLoading.value = true
  if (!isFormValid.value) {
    alert.value.handleAlert('Campos necessários não preenchidos', AlertType.Warning)
    return
  }

  const body: Partial<UpdateCategRequest> = {
    user_id: formUser.value,
    name: formName.value,
  }

  try {
    const res: QueryResponse = await updateCategory(userId, categId, body)
    alert.value.handleAlert(res.message, codeToAlertType(res.code))
    emits('submitted')
    showModel.value = false
  } catch {
    alert.value.handleAlert('Erro desconhecido. Tente novamente mais tarde.', AlertType.Error)
  } finally {
    isLoading.value = false
  }
}

// Selecionar o usuário atual na lista suspensa, na primeira renderização
onMounted((): void => {
  formUser.value = props.categ.user_id
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
  const u = formUser.value !== props.categ.user_id
  const n = formName.value.length > 0
  isFilled.value = u || n

  isFormValid.value = isFilled.value
})
</script>

<template>
  <PopupWindow :title="`Editar categoria`" v-model="showModel">
    <form
      class="flex flex-col gap-4 space-y-4 px-8 py-4"
      @submit.prevent="() => handleUpdateCategory(categ.user_id, categ.categ_id)"
    >
      <div class="flex w-full flex-col gap-4">
        <!-- Campos do formulário -->
        <InputList
          :values="users"
          label="Usuário"
          v-model="formUser"
          :selected="categ.user_id"
          :left-inner-icon="PhUser"
          :required="false"
        />
        <InputText
          :placeholder="categ.name"
          label="Nome"
          v-model="formName"
          :left-inner-icon="PhFolder"
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

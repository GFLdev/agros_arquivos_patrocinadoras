<script setup lang="ts">
import { PhIdentificationBadge, PhPassword, PhPencil, PhUser, PhXCircle } from '@phosphor-icons/vue'
import InputPassword from '@/components/generic/InputPassword.vue'
import InputText from '@/components/generic/InputText.vue'
import SubmitButton from '@/components/generic/SubmitButton.vue'
import CancelButton from '@/components/generic/CancelButton.vue'
import PopupWindow from '@/components/generic/PopupWindow.vue'
import type { UserRequest } from '@/@types/Requests.ts'
import type { QueryResponse, UserModel } from '@/@types/Responses.ts'
import { type EmitFn, type ModelRef, type PropType, type Ref, ref, watch, watchEffect } from 'vue'
import { checkPasswords, validatePassword, validateUsername } from '@/utils/validate.ts'
import { updateUser } from '@/services/queries.ts'
import PopupAlert from '@/components/generic/PopupAlert.vue'
import { AlertType } from '@/@types/Enumerations.ts'
import { Alert, codeToAlertType, TRANSITION_DURATION } from '@/utils/modals.ts'

defineProps({
  user: {
    type: Object as PropType<UserModel>,
    required: true,
  },
})

// Emissores
const emits: EmitFn<'submitted'[]> = defineEmits(['submitted'])

// Formulário
const formUsername: Ref<string> = ref<string>('')
const formName: Ref<string> = ref<string>('')
const formPasswd: Ref<string> = ref<string>('')
const formConfirmPasswd: Ref<string> = ref<string>('')

// Status de carregamento
const isLoading: Ref<boolean> = ref<boolean>(false)

// Validações
const isFilled: Ref<boolean> = ref<boolean>(false)
const isPasswdMatched: Ref<boolean> = ref<boolean>(false)
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
  formUsername.value = ''
  formName.value = ''
  formPasswd.value = ''
  formConfirmPasswd.value = ''
}

/**
 * Lida com a atualização de um usuário enviando uma solicitação para atualizar os detalhes do usuário
 * e gerenciando a validação do formulário.
 *
 * @param {string} userId - O identificador único do usuário a ser atualizado.
 * @return {Promise<void>} Uma promessa que é resolvida quando o processo de atualização do usuário é concluído.
 */
async function handleUpdateUser(userId: string): Promise<void> {
  isLoading.value = true
  if (!isFormValid.value) {
    alert.value.handleAlert('Campos necessários não preenchidos', AlertType.Warning)
    return
  }

  const body: Partial<UserRequest> = {
    username: formUsername.value,
    name: formName.value,
    password: formPasswd.value,
  }

  try {
    const res: QueryResponse = await updateUser(userId, body)
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
  const u: boolean = (formUsername.value ?? '').length > 0
  const n: boolean = (formName.value ?? '').length > 0
  const p: boolean = (formPasswd.value ?? '').length > 0
  const cp: boolean = (formConfirmPasswd.value ?? '').length > 0
  isFilled.value = u || n || p || cp

  // Verificar senhas e tamanho dos parâmetros
  isPasswdMatched.value = true
  let validLengths: boolean = true
  if (u) validLengths &&= validateUsername(formUsername.value)
  if (p || cp) {
    isPasswdMatched.value = checkPasswords(formPasswd.value, formConfirmPasswd.value)
    validLengths &&= validatePassword(formPasswd.value)
  }

  isFormValid.value = isFilled.value && isPasswdMatched.value && validLengths
})
</script>

<template>
  <PopupWindow :title="`Editar usuário`" v-model="showModel">
    <form class="flex flex-col gap-4 space-y-4 px-8 py-4" @submit.prevent="() => handleUpdateUser(user.user_id)">
      <div class="flex w-full flex-col gap-4">
        <!-- Campos do formulário -->
        <InputText
          :placeholder="user.username"
          label="Usuário"
          v-model="formUsername"
          :left-inner-icon="PhIdentificationBadge"
          :required="false"
        />
        <InputText
          :placeholder="user.name"
          label="Nome Completo"
          v-model="formName"
          :left-inner-icon="PhUser"
          :required="false"
        />
        <InputPassword
          placeholder="&#9679;&#9679;&#9679;&#9679;&#9679;"
          label="Senha"
          v-model="formPasswd"
          :left-inner-icon="PhPassword"
          :showable="true"
          :required="false"
        />
        <InputPassword
          placeholder="&#9679;&#9679;&#9679;&#9679;&#9679;"
          label="Confirmar Senha"
          v-model="formConfirmPasswd"
          :left-inner-icon="PhPassword"
          :showable="true"
          :required="!!formPasswd"
        />
        <!-- Avisos -->
        <div v-if="!isFormValid" class="w-full text-center text-sm font-light">
          <p v-if="!isFilled">Preencha ao menos um campo</p>
          <p v-if="(formUsername?.length ?? 0) > 0 && !validateUsername(formUsername)">
            O usuário deve ter entre 4 e 16 caracteres
          </p>
          <p v-if="(formPasswd?.length ?? 0) > 0 && !validatePassword(formPasswd)">
            A senha deve ter pelo menos 4 caracteres
          </p>
          <p v-else-if="(formPasswd?.length ?? 0) > 0 && !isPasswdMatched">As senhas devem sem iguais</p>
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

<script setup lang="ts">
import { PhIdentificationBadge, PhPassword, PhUser, PhUserPlus, PhXCircle } from '@phosphor-icons/vue'
import InputPassword from '@/components/generic/InputPassword.vue'
import InputText from '@/components/generic/InputText.vue'
import SubmitButton from '@/components/generic/SubmitButton.vue'
import CancelButton from '@/components/generic/CancelButton.vue'
import PopupWindow from '@/components/generic/PopupWindow.vue'
import type { UserRequest } from '@/@types/Requests.ts'
import { type EmitFn, type ModelRef, type Ref, ref, watch, watchEffect } from 'vue'
import { checkPasswords, validatePassword, validateUsername } from '@/utils/validate.ts'
import { createUser } from '@/services/queries.ts'
import PopupAlert from '@/components/generic/PopupAlert.vue'
import { AlertType } from '@/@types/Enumerations.ts'
import { Alert, codeToAlertType, TRANSITION_DURATION } from '@/utils/modals.ts'
import type { QueryResponse } from '@/@types/Responses.ts'

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
const isMatched: Ref<boolean> = ref<boolean>(false)
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
 * Lida com a criação de um novo usuário, validando os dados do formulário, enviando uma requisição de rede
 * e atualizando a interface do usuário com base na resposta.
 *
 * @return {Promise<void>} Uma promessa que é resolvida quando o processo de criação do usuário é concluído.
 */
async function handleCreateUser(): Promise<void> {
  isLoading.value = true
  if (!isFormValid.value) {
    alert.value.handleAlert('Campos necessários não preenchidos', AlertType.Warning)
    return
  }

  const body: UserRequest = {
    username: formUsername.value,
    name: formName.value,
    password: formPasswd.value,
  }

  try {
    const res: QueryResponse = await createUser(body)
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
  const u: boolean = formUsername.value.length > 0
  const n: boolean = formName.value.length > 0
  const p: boolean = formPasswd.value.length > 0
  const cp: boolean = formConfirmPasswd.value.length > 0
  isFilled.value = u && n && p && cp

  // Verificar senhas e tamanho dos parâmetros
  isMatched.value = checkPasswords(formPasswd.value, formConfirmPasswd.value)
  const validLengths: boolean = validateUsername(formUsername.value) && validatePassword(formPasswd.value)

  isFormValid.value = isFilled.value && isMatched.value && validLengths
})
</script>

<template>
  <PopupWindow :title="`Criar novo usuário`" v-model="showModel">
    <form class="flex flex-col gap-4 space-y-4 px-8 py-4" @submit.prevent="handleCreateUser">
      <div class="flex w-full flex-col gap-4">
        <!-- Campos do formulário -->
        <InputText
          placeholder="Nome de usuário"
          label="Usuário"
          v-model="formUsername"
          :left-inner-icon="PhIdentificationBadge"
          :required="true"
        />
        <InputText
          placeholder="Nome Completo"
          label="Nome Completo"
          v-model="formName"
          :left-inner-icon="PhUser"
          :required="true"
        />
        <InputPassword
          placeholder="&#9679;&#9679;&#9679;&#9679;&#9679;"
          label="Senha"
          v-model="formPasswd"
          :left-inner-icon="PhPassword"
          :showable="true"
          :required="true"
        />
        <InputPassword
          placeholder="&#9679;&#9679;&#9679;&#9679;&#9679;"
          label="Confirmar Senha"
          v-model="formConfirmPasswd"
          :left-inner-icon="PhPassword"
          :showable="true"
          :required="true"
        />
        <!-- Avisos -->
        <div v-if="!isFormValid" class="w-full text-center text-sm font-light">
          <p v-if="!isFilled">Preencha todos os campos</p>
          <p v-if="!validateUsername(formUsername)">O usuário deve ter entre 4 e 16 caracteres</p>
          <p v-if="!validatePassword(formPasswd)">A senha deve ter pelo menos 4 caracteres</p>
          <p v-else-if="!isMatched">As senhas devem sem iguais</p>
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
          :left-inner-icon="PhUserPlus"
        />
      </div>
    </form>
  </PopupWindow>
  <PopupAlert :text="alert.text" :type="alert.type" :duration="alert.duration" v-model="alert.show" />
</template>

<style scoped></style>

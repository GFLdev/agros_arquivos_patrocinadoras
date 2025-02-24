<script setup lang="ts">
import { PhIdentificationBadge, PhPassword, PhPencil, PhUser, PhXCircle } from '@phosphor-icons/vue'
import InputPassword from '@/components/generic/InputPassword.vue'
import InputText from '@/components/generic/InputText.vue'
import SubmitButton from '@/components/generic/SubmitButton.vue'
import CancelButton from '@/components/generic/CancelButton.vue'
import PopupWindow from '@/components/generic/PopupWindow.vue'
import type { UserRequest } from '@/@types/Requests.ts'
import type { UserModel } from '@/@types/Responses.ts'
import { type PropType, ref, watch, watchEffect } from 'vue'
import { validatePassword, validateUsername } from '@/utils/validate.ts'
import { updateUser } from '@/services/queries.ts'
import PopupAlert from '@/components/generic/PopupAlert.vue'
import { AlertType } from '@/@types/Enumerations.ts'
import { codeToAlertType } from '@/utils/modals.ts'

defineProps({
  user: {
    type: Object as PropType<UserModel>,
    required: true,
  },
})

const emits = defineEmits(['submitted'])

// Formulário
const username = ref<string>('')
const name = ref<string>('')
const passwd = ref<string>('')
const confirmPasswd = ref<string>('')

// Validações
const loading = ref<boolean>(false)
const filled = ref<boolean>(false)
const matched = ref<boolean>(false)
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

// Função para abrir janela de atualização de usuário
async function handleUpdateUser(userId: string) {
  if (!formValid.value) {
    handleAlert('Campos necessários não preenchidos', AlertType.Warning)
    return
  }
  loading.value = true

  const body: UserRequest = {
    username: username.value,
    name: name.value,
    password: passwd.value,
  }

  try {
    const res = await updateUser(userId, body)
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

// Função para verificar se as senhas são iguais
function checkPasswords(): boolean {
  return passwd.value === confirmPasswd.value
}

// Função para reset das variáveis reativas
function reset() {
  loading.value = false
  username.value = ''
  name.value = ''
  passwd.value = ''
  confirmPasswd.value = ''
}

watchEffect(() => {
  const u = username.value.length > 0
  const n = name.value.length > 0
  const p = passwd.value.length > 0
  const cp = confirmPasswd.value.length > 0

  filled.value = u || n || p || cp
  matched.value = true

  let validLengths = true
  if (p || cp) {
    matched.value = checkPasswords()
    validLengths &&= validatePassword(passwd.value)
  }
  if (u) {
    validLengths &&= validateUsername(username.value)
  }

  formValid.value = filled.value && matched.value && validLengths
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
  <PopupWindow :title="`Editar usuário`" v-model="showModel">
    <form class="flex flex-col gap-4 space-y-4 px-8 py-4" @submit.prevent="() => handleUpdateUser(user.user_id)">
      <div class="flex w-full flex-col gap-4">
        <!-- Campos do formulário -->
        <InputText
          :placeholder="user.username"
          label="Usuário"
          v-model="username"
          :left-inner-icon="PhIdentificationBadge"
          :required="false"
        />
        <InputText
          :placeholder="user.name"
          label="Nome Completo"
          v-model="name"
          :left-inner-icon="PhUser"
          :required="false"
        />
        <InputPassword
          placeholder="&#9679;&#9679;&#9679;&#9679;&#9679;"
          label="Senha"
          v-model="passwd"
          :left-inner-icon="PhPassword"
          :showable="true"
          :required="false"
        />
        <InputPassword
          placeholder="&#9679;&#9679;&#9679;&#9679;&#9679;"
          label="Confirmar Senha"
          v-model="confirmPasswd"
          :left-inner-icon="PhPassword"
          :showable="true"
          :required="!!passwd"
        />
        <!-- Avisos -->
        <div v-if="!formValid" class="w-full text-center text-sm font-light">
          <p v-if="!filled">Preencha ao menos um campo</p>
          <p v-if="username.length > 0 && !validateUsername(username)">O usuário deve ter entre 4 e 16 caracteres</p>
          <p v-if="passwd.length > 0 && !validatePassword(passwd)">A senha deve ter pelo menos 4 caracteres</p>
          <p v-else-if="passwd.length > 0 && !matched">As senhas devem sem iguais</p>
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
          text="Editar"
          loading-text="Editando"
          :loading="loading"
          :disabled="loading || !formValid"
          :left-inner-icon="PhPencil"
        />
      </div>
    </form>
  </PopupWindow>
  <PopupAlert :text="alertText" :type="alertType" :duration="alertDuration" v-model="showAlert" />
</template>

<style scoped></style>

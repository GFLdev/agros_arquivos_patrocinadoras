<script setup lang="ts">
import PopupWindow from '@/components/generic/PopupWindow.vue'
import { PhUserMinus, PhXCircle } from '@phosphor-icons/vue'
import SubmitButton from '@/components/generic/SubmitButton.vue'
import CancelButton from '@/components/generic/CancelButton.vue'
import { type PropType, ref } from 'vue'
import { deleteUser } from '@/services/queries.ts'
import type { UserModel } from '@/@types/Responses.ts'
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

// Validações
const loading = ref<boolean>(false)

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

// Função para abrir janela de exclusão de usuário
async function handleDeleteUser(userId: string) {
  loading.value = true

  try {
    const res = await deleteUser(userId)
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
</script>

<template>
  <PopupWindow :title="`Excluir usuário`" v-model="showModel">
    <form class="flex flex-col gap-4 space-y-4 px-8 py-4" @submit.prevent="() => handleDeleteUser(user.user_id)">
      <div class="w-full text-center font-light">
        <p>
          Deseja realmente excluir o usuário<br /><b>{{ user.name }}</b
          >?
        </p>
        <br />
        <p>Obs: Esta ação não pode ser desfeita.</p>
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
          text="Excluir"
          loading-text="Excluindo"
          :loading="loading"
          :disabled="loading"
          :left-inner-icon="PhUserMinus"
        />
      </div>
    </form>
  </PopupWindow>
  <PopupAlert :text="alertText" :type="alertType" :duration="alertDuration" v-model="showAlert" />
</template>

<style scoped></style>

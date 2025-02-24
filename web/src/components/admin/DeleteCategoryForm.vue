<script setup lang="ts">
import PopupWindow from '@/components/generic/PopupWindow.vue'
import { PhFolderMinus, PhXCircle } from '@phosphor-icons/vue'
import SubmitButton from '@/components/generic/SubmitButton.vue'
import CancelButton from '@/components/generic/CancelButton.vue'
import { type PropType, ref } from 'vue'
import { deleteCategory } from '@/services/queries.ts'
import type { CategModel } from '@/@types/Responses.ts'
import PopupAlert from '@/components/generic/PopupAlert.vue'
import { AlertType } from '@/@types/Enumerations.ts'
import { codeToAlertType } from '@/utils/modals.ts'

defineProps({
  categ: {
    type: Object as PropType<CategModel>,
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

// Função para abrir janela de exclusão de categoria
async function handleDeleteCateg(userId: string, categId: string) {
  loading.value = true

  try {
    const res = await deleteCategory(userId, categId)
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
  <PopupWindow :title="`Excluir categoria`" v-model="showModel">
    <form
      class="flex flex-col gap-4 space-y-4 px-8 py-4"
      @submit.prevent="() => handleDeleteCateg(categ.user_id, categ.categ_id)"
    >
      <div class="w-full text-center font-light">
        <p>
          Deseja realmente excluir a categoria<br /><b>{{ categ.name }}</b
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
          :left-inner-icon="PhFolderMinus"
        />
      </div>
    </form>
  </PopupWindow>
  <PopupAlert :text="alertText" :type="alertType" :duration="alertDuration" v-model="showAlert" />
</template>

<style scoped></style>

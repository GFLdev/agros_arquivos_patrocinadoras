<script setup lang="ts">
import PopupWindow from '@/components/generic/PopupWindow.vue'
import { PhUserMinus, PhXCircle } from '@phosphor-icons/vue'
import SubmitButton from '@/components/generic/SubmitButton.vue'
import CancelButton from '@/components/generic/CancelButton.vue'
import { type EmitFn, type ModelRef, type PropType, type Ref, ref } from 'vue'
import { deleteUser } from '@/services/queries.ts'
import type { QueryResponse, UserModel } from '@/@types/Responses.ts'
import PopupAlert from '@/components/generic/PopupAlert.vue'
import { AlertType } from '@/@types/Enumerations.ts'
import { Alert, codeToAlertType } from '@/utils/modals.ts'

defineProps({
  user: {
    type: Object as PropType<UserModel>,
    required: true,
  },
})

// Emissores
const emits: EmitFn<'submitted'[]> = defineEmits(['submitted'])

// Status de carregamento
const isLoading: Ref<boolean> = ref<boolean>(false)

// Alerta
const alert: Ref<Alert> = ref<Alert>(new Alert())

// Estado de visibilidade
const showModel: ModelRef<boolean | undefined> = defineModel<boolean>()

/**
 * Lida com a lógica para excluir um usuário específico.
 *
 * @param {string} userId - O identificador único do usuário a ser excluído.
 * @return {Promise<void>} Retorna uma promessa resolvida com void.
 */
async function handleDeleteUser(userId: string): Promise<void> {
  isLoading.value = true

  try {
    const res: QueryResponse = await deleteUser(userId)
    alert.value.handleAlert(res.message, codeToAlertType(res.code))
    emits('submitted')
    showModel.value = false
  } catch {
    alert.value.handleAlert('Erro desconhecido. Tente novamente mais tarde.', AlertType.Error)
  } finally {
    isLoading.value = false
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
          :disabled="isLoading"
          :left-inner-icon="PhXCircle"
        />
        <SubmitButton
          text="Excluir"
          loading-text="Excluindo"
          :loading="isLoading"
          :disabled="isLoading"
          :left-inner-icon="PhUserMinus"
        />
      </div>
    </form>
  </PopupWindow>
  <PopupAlert :text="alert.text" :type="alert.type" :duration="alert.duration" v-model="alert.show" />
</template>

<style scoped></style>

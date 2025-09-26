<script setup lang="ts">
import PopupWindow from '@/components/generic/PopupWindow.vue'
import { PhFolderMinus, PhXCircle } from '@phosphor-icons/vue'
import SubmitButton from '@/components/generic/SubmitButton.vue'
import CancelButton from '@/components/generic/CancelButton.vue'
import { type EmitFn, type ModelRef, type PropType, type Ref, ref } from 'vue'
import { deleteCategory } from '@/services/queries.ts'
import type { CategModel, QueryResponse } from '@/@types/Responses.ts'
import PopupAlert from '@/components/generic/PopupAlert.vue'
import { AlertType } from '@/@types/Enumerations.ts'
import { Alert, codeToAlertType } from '@/utils/modals.ts'

defineProps({
  categ: {
    type: Object as PropType<CategModel>,
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
 * Lida com a exclusão de uma categoria específica para um determinado usuário.
 *
 * @param {string} userId - O identificador único do usuário.
 * @param {string} categId - O identificador único da categoria a ser excluída.
 * @return {Promise<void>} Uma promessa que é resolvida quando o processo de exclusão da categoria é concluído.
 */
async function handleDeleteCateg(userId: string, categId: string): Promise<void> {
  isLoading.value = true

  try {
    const res: QueryResponse = await deleteCategory(userId, categId)
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
          :disabled="isLoading"
          :left-inner-icon="PhXCircle"
        />
        <SubmitButton
          text="Excluir"
          loading-text="Excluindo"
          :loading="isLoading"
          :disabled="isLoading"
          :left-inner-icon="PhFolderMinus"
        />
      </div>
    </form>
  </PopupWindow>
  <PopupAlert :text="alert.text" :type="alert.type" :duration="alert.duration" v-model="alert.show" />
</template>

<style scoped></style>

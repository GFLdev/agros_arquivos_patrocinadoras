<script setup lang="ts">
import PopupWindow from '@/components/generic/PopupWindow.vue'
import { PhFileMinus, PhXCircle } from '@phosphor-icons/vue'
import SubmitButton from '@/components/generic/SubmitButton.vue'
import CancelButton from '@/components/generic/CancelButton.vue'
import { type EmitFn, type ModelRef, type PropType, type Ref, ref } from 'vue'
import { deleteFile } from '@/services/queries.ts'
import type { CategModel, FileModel, QueryResponse } from '@/@types/Responses.ts'
import PopupAlert from '@/components/generic/PopupAlert.vue'
import { AlertType } from '@/@types/Enumerations.ts'
import { Alert, codeToAlertType } from '@/utils/modals.ts'

defineProps({
  categ: {
    type: Object as PropType<CategModel>,
    required: true,
  },
  file: {
    type: Object as PropType<FileModel>,
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
 * Lida com a exclusão de um arquivo para um usuário e categoria específicos.
 *
 * @param {string} userId O identificador único do usuário solicitando a exclusão do arquivo.
 * @param {string} categId O identificador único da categoria à qual o arquivo pertence.
 * @param {string} fileId O identificador único do arquivo a ser excluído.
 * @return {Promise<void>} Uma promessa que é resolvida quando o processo de exclusão do arquivo é concluído.
 */
async function handleDeleteFile(userId: string, categId: string, fileId: string): Promise<void> {
  isLoading.value = true

  try {
    const res: QueryResponse = await deleteFile(userId, categId, fileId)
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
  <PopupWindow :title="`Excluir arquivo`" v-model="showModel">
    <form
      class="flex flex-col gap-4 space-y-4 px-8 py-4"
      @submit.prevent="() => handleDeleteFile(categ.user_id, file.categ_id, file.file_id)"
    >
      <div class="w-full text-center font-light">
        <p>
          Deseja realmente excluir o arquivo<br /><b>{{ file.name }}</b
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
          :left-inner-icon="PhFileMinus"
        />
      </div>
    </form>
  </PopupWindow>
  <PopupAlert :text="alert.text" :type="alert.type" :duration="alert.duration" v-model="alert.show" />
</template>

<style scoped></style>

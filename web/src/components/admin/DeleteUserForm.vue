<script setup lang="ts">
import PopupWindow from '@/components/generic/PopupWindow.vue'
import { PhUserMinus, PhXCircle } from '@phosphor-icons/vue'
import SubmitButton from '@/components/generic/SubmitButton.vue'
import CancelButton from '@/components/generic/CancelButton.vue'
import { ref } from 'vue'
import apiClient from '@/services/axios.ts'
import router from '@/router'
import type { AxiosError } from 'axios'

defineProps({
  userId: {
    type: String,
    required: true,
  },
})

// Validações
const loading = ref<boolean>(false)

// Estado de visibilidade
const showModel = defineModel<boolean>()

// Função para abrir janela de exclusão de usuário
async function handleDeleteUser(user: string) {
  loading.value = true

  try {
    const res = await apiClient.delete('/auth/user/' + user)
    const message = res.data?.message
    console.log(message)
    router.go(0)
  } catch (e: unknown) {
    const error = e as AxiosError
    if (error.response && error.response.status === 401) {
      // TODO: Exiba uma mensagem amigável para o usuário
      console.error('Credenciais inválidas')
    } else {
      // TODO: Tratar erros de login e exibir mensagens relevantes ao usuário
      console.error('Erro ao excluir usuário:', error.message || error)
    }
  }
}
</script>

<template>
  <PopupWindow :title="`Excluir usuário`" v-model="showModel">
    <form class="flex flex-col gap-4 space-y-4 overflow-scroll px-8 py-4">
      <div class="w-full text-center">
        <p>Deseja realmente excluir este usuário?</p>
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
          :on-click="() => handleDeleteUser(userId)"
          :loading="loading"
          :disabled="loading"
          :left-inner-icon="PhUserMinus"
        />
      </div>
    </form>
  </PopupWindow>
</template>

<style scoped></style>

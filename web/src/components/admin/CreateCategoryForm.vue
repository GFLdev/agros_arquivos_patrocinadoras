<script setup lang="ts">
import { PhFolder, PhFolderPlus, PhXCircle } from '@phosphor-icons/vue'
import InputText from '@/components/generic/InputText.vue'
import SubmitButton from '@/components/generic/SubmitButton.vue'
import CancelButton from '@/components/generic/CancelButton.vue'
import PopupWindow from '@/components/generic/PopupWindow.vue'
import type { CategRequest } from '@/@types/Requests.ts'
import apiClient from '@/services/axios.ts'
import type { AxiosError } from 'axios'
import { ref, watchEffect } from 'vue'
import router from '@/router'
defineProps({
  userId: {
    type: String,
    required: true,
  },
})

// Formulário
const name = ref<string | null>(null)

// Validações
const loading = ref<boolean>(false)
const filled = ref<boolean>(false)
const formValid = ref<boolean>(false)

// Estado de visibilidade
const showModel = defineModel<boolean>()

// Função para abrir janela de criação de usuário
async function handleCreateCategory(user: string) {
  if (!name.value) {
    // TODO: Manipular erro
    return
  }
  loading.value = true

  const body: CategRequest = {
    name: name.value,
  }

  try {
    const res = await apiClient.post(`/auth/user/${user}/category`, JSON.stringify(body))
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
      console.error('Erro ao criar categoria:', error.message || error)
    }
  }
}

watchEffect(() => {
  filled.value = !!name.value && name.value.length > 0
  formValid.value = filled.value
})
</script>

<template>
  <PopupWindow :title="`Criar nova categoria`" v-model="showModel">
    <form class="flex flex-col gap-4 space-y-4 overflow-scroll px-8 py-4">
      <div class="flex w-full flex-col gap-4">
        <!-- Campos do formulário -->
        <InputText
          placeholder="Nome da categoria"
          label="Nome"
          v-model="name"
          :left-inner-icon="PhFolder"
          :required="true"
        />
        <!-- Avisos -->
        <div v-if="!formValid" class="w-full text-center text-sm font-light">
          <p v-if="!formValid">Preencha todos os campos</p>
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
          text="Criar"
          loading-text="Criando"
          :on-click="() => handleCreateCategory(userId)"
          :loading="loading"
          :disabled="loading || !formValid"
          :left-inner-icon="PhFolderPlus"
        />
      </div>
    </form>
  </PopupWindow>
</template>

<style scoped></style>

<script setup lang="ts">
import { PhFolder, PhPencil, PhUser, PhXCircle } from '@phosphor-icons/vue'
import InputText from '@/components/generic/InputText.vue'
import SubmitButton from '@/components/generic/SubmitButton.vue'
import CancelButton from '@/components/generic/CancelButton.vue'
import PopupWindow from '@/components/generic/PopupWindow.vue'
import type { UpdateCategRequest } from '@/@types/Requests.ts'
import { computed, onMounted, type PropType, ref, watch, watchEffect } from 'vue'
import type { CategModel } from '@/@types/Responses.ts'
import InputList from '@/components/generic/InputList.vue'
import { updateCategory } from '@/services/queries.ts'
import PopupAlert from '@/components/generic/PopupAlert.vue'
import { AlertType } from '@/@types/Enumerations.ts'
import { codeToAlertType } from '@/utils/modals.ts'

const props = defineProps({
  categ: {
    type: Object as PropType<CategModel>,
    required: true,
  },
  users: {
    type: Map<string, string>,
    required: true,
  },
})

const emits = defineEmits(['submitted'])

// Formulário
const selectedUser = ref<string>(props.categ.user_id)
const name = ref<string>('')

// Validações
const loading = ref<boolean>(false)
const filled = ref<boolean>(false)
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

// Função para abrir janela de atualização de categoria
async function handleUpdateCategory(userId: string, categId: string) {
  if (!formValid.value) {
    handleAlert('Campos necessários não preenchidos', AlertType.Warning)
    return
  }
  loading.value = true

  const body: UpdateCategRequest = {
    user_id: selectedUser.value,
    name: name.value,
  }

  try {
    const res = await updateCategory(userId, categId, body)
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

// Função para reset das variáveis reativas
function reset() {
  loading.value = false
  name.value = ''
  selectedUser.value = computed(() => props.categ.user_id).value
}

onMounted(() => {
  selectedUser.value = props.categ.user_id
})

watchEffect(() => {
  const u = selectedUser.value !== props.categ.user_id
  const n = name.value.length > 0
  filled.value = u || n
  formValid.value = filled.value
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
  <PopupWindow :title="`Editar categoria`" v-model="showModel">
    <form
      class="flex flex-col gap-4 space-y-4 px-8 py-4"
      @submit.prevent="() => handleUpdateCategory(categ.user_id, categ.categ_id)"
    >
      <div class="flex w-full flex-col gap-4">
        <!-- Campos do formulário -->
        <InputList
          :values="users"
          label="Usuário"
          v-model="selectedUser"
          :selected="categ.user_id"
          :left-inner-icon="PhUser"
          :required="false"
        />
        <InputText
          :placeholder="categ.name"
          label="Nome"
          v-model="name"
          :left-inner-icon="PhFolder"
          :required="false"
        />
        <!-- Avisos -->
        <div v-if="!formValid" class="w-full text-center text-sm font-light">
          <p v-if="!filled">Preencha ao menos um campo</p>
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

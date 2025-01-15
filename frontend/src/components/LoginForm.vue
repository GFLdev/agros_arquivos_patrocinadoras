<script setup lang="ts">
import { ref } from 'vue'
import {
  PhEyeClosed,
  PhEye,
  PhSignIn,
  PhCircleNotch,
  PhIdentificationCard,
  PhPassword,
} from '@phosphor-icons/vue'

const visible = ref<boolean>(false)
const username = ref<string | null>()
const passwd = ref<string | null>()
const loading = ref<boolean>(false)

// Função para gerenciar o envio do formulário para autenticação do usuário
async function handleSignIn(): Promise<void> {
  visible.value = false
  loading.value = true
  // TODO: Criar lógica
}
</script>

<template>
  <div
    class="flex w-full flex-col items-center justify-center gap-6 rounded-none bg-white px-8 py-4 drop-shadow-2xl sm:mx-16 sm:w-fit sm:rounded-lg sm:px-16 sm:py-8"
  >
    <p class="text-center text-lg font-light text-agros-gray-dark">
      Por favor, preencha os campos para fazer o login.
    </p>
    <section class="flex w-full flex-col justify-center gap-6 px-8">
      <!-- Input de usuário -->
      <div class="relative mt-2 grid grid-cols-1">
        <input
          type="text"
          name="username"
          id="username"
          class="text-dark outline-gray placeholder:text-gray focus:outline-dark peer col-start-1 row-start-1 block w-full rounded-md bg-white py-1.5 pl-10 pr-3 text-base outline outline-1 -outline-offset-1 focus:outline focus:outline-2 focus:-outline-offset-2 sm:text-sm/6"
          placeholder="Nome de usuário"
          v-model="username"
          required
          :disabled="loading"
        />
        <label
          for="username"
          class="text-gray peer-focus:text-dark absolute -top-2 left-2 inline-block rounded-lg bg-white px-1 text-xs font-medium"
          >Usuário <b>*</b></label
        >
        <PhIdentificationCard
          class="text-gray peer-focus:text-dark pointer-events-none col-start-1 row-start-1 ml-3 size-5 self-center"
          aria-hidden="true"
        />
      </div>
      <!-- Input de senha -->
      <div class="relative mt-2 flex">
        <div class="relative -mr-px grid grow grid-cols-1">
          <input
            name="passwd"
            id="passwd"
            class="text-dark outline-gray placeholder:text-gray focus:outline-dark peer col-start-1 row-start-1 block w-full rounded-l-md bg-white pr-3 pl-10 py-1.5 text-base outline outline-1 -outline-offset-1 focus:outline focus:outline-2 focus:-outline-offset-2 sm:text-sm/6"
            placeholder="&#9679;&#9679;&#9679;&#9679;&#9679;&#9679;&#9679;&#9679;"
            v-model="passwd"
            :disabled="loading"
            :type="visible ? 'text' : 'password'"
          />
          <label
            for="passwd"
            class="text-gray peer-focus:text-dark absolute -top-2 left-2 inline-block rounded-lg bg-white px-1 text-xs font-medium"
            >Senha <b>*</b></label
          >
          <PhPassword
            class="text-gray peer-focus:text-dark pointer-events-none col-start-1 row-start-1 ml-3 size-5 self-center"
            aria-hidden="true"
          />
        </div>
        <button
          type="button"
          class="text-dark outline-gray hover:bg-gray focus:outline-dark flex shrink-0 items-center gap-x-1.5 rounded-r-md bg-white px-3 py-2 text-sm font-semibold outline outline-1 -outline-offset-1 transition-colors duration-200 hover:bg-opacity-30 focus:relative focus:outline focus:outline-2 focus:-outline-offset-2"
          :disabled="loading"
          @click="visible = !visible"
        >
          <PhEyeClosed v-if="visible" class="text-dark -ml-0.5 size-4" aria-hidden="true" />
          <PhEye v-else class="text-dark -ml-0.5 size-4" aria-hidden="true" />
        </button>
      </div>
    </section>
    <!-- Botão de login -->
    <button
      type="button"
      :class="
        'text-dark focus-visible:outline-dark enabled:hover:bg-dark ' +
        ' focus-visible:outline-offset inline-flex items-center gap-x-2 rounded-md px-3 py-1.5 text-base shadow-sm' +
        ' transition duration-200 ease-in-out focus:-outline-offset-2 focus-visible:outline focus-visible:outline-2' +
        ' enabled:hover:text-white' +
        (loading ? ' disabled:bg-dark disabled:text-white' : '')
      "
      :disabled="!username || !passwd || loading"
      @click="handleSignIn"
    >
      <PhSignIn v-if="!loading" class="size-5" aria-hidden="true" />
      <PhCircleNotch v-else class="size-5 animate-spin" aria-hidden="true" />
      <span v-if="!loading">Entrar</span>
      <span v-else>Entrar</span>
    </button>
  </div>
</template>

<style scoped>
*:disabled {
  @apply opacity-50;
}

*:enabled {
  @apply opacity-100;
}
</style>

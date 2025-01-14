<script setup lang="ts">
import { PhDoorOpen, PhEye, PhEyeClosed, PhIdentificationBadge, PhLock } from '@phosphor-icons/vue'
import { ref } from 'vue'

const visible = ref<boolean>(false)
const username = ref<string | null>()
const passwd = ref<string | null>()
</script>

<template>
  <v-list
    class="flex flex-col items-center px-8 py-4 justify-center max-w-xl rounded-lg"
    lines="two"
    elevation="5"
  >
    <v-list-item class="text-center text-lg">
      Por favor, preencha os campos para fazer o login.
    </v-list-item>
    <section class="flex flex-col justify-center px-8 py-4 w-full">
      <v-text-field
        class="login-form-text-field"
        label="Nome de usuário"
        :prepend-icon="PhIdentificationBadge"
        density="comfortable"
        variant="underlined"
        v-model="username"
      ></v-text-field>
      <v-text-field
        class="login-form-text-field"
        label="Senha"
        :prepend-icon="PhLock"
        :append-inner-icon="visible ? PhEyeClosed : PhEye"
        :type="visible ? 'text' : 'password'"
        density="comfortable"
        variant="underlined"
        @click:append-inner="visible = !visible"
        v-model="passwd"
      ></v-text-field>
    </section>
    <v-btn
      id="login-form-btn"
      :prepend-icon="PhDoorOpen"
      variant="plain"
      :disabled="!username || !passwd">Entrar</v-btn>
  </v-list>
</template>

<style scoped>
.login-form-text-field {
  @apply w-full
  !important;
}

#login-form-btn {
  @apply w-fit
  !important;
}
</style>

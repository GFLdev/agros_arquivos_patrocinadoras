<script setup lang="ts">
import Header from '@/components/HeaderSection.vue'
import { PhAt, PhMapPin, PhPhone } from '@phosphor-icons/vue'
</script>

<template>
  <Header />
  <main class="flex w-screen flex-col items-center justify-center gap-8 px-16 py-8">
    <section class="w-fit max-w-5xl rounded-lg px-4 py-2 text-center text-lg">
      <p>
        Nesta página estão disponíveis informações sobre os planos de previdência e de saúde do
        Agros, em atendimento à Resolução Normativa 389, da Agência Nacional de Saúde Suplementar
        (ANS).
      </p>
    </section>
    <div class="w-full max-w-5xl">
      <section class="w-full">
        <div v-for="i in 3" :key="i" class="w-full">
          <h2 :id="`accordion-heading-${i}`">
            <button
              type="button"
              class="text-gray-500 border-gray-200 flex w-full items-center justify-between gap-3 border-b py-5 text-right font-medium"
            >
              <span>Title</span>
            </button>
          </h2>
          <div :id="`accordion-body-${i}`" class="hidden">
            <div class="border-gray-200 border-b py-5">
              <p class="text-gray-500 mb-2">Content</p>
            </div>
          </div>
        </div>
      </section>
    </div>
    <section
      class="flex w-fit max-w-3xl flex-col items-center justify-center rounded-lg px-4 py-2 text-center"
    >
      <h2 class="text-agros-gray-dark text-xl font-bold">Dúvidas?</h2>
      <ul>
        <li>
          <p>
            Em caso de dúvida sobre as informações aqui apresentadas, entre em contato com o Agros.
          </p>
        </li>
        <li>
          <p>
            <PhMapPin />
            Av. Purdue, s/n, Campus da UFV - Viçosa
          </p>
        </li>
        <li>
          <p>
            <PhPhone />
            (31) 3899-6550 - ramal 6539
          </p>
        </li>
        <li>
          <p>
            <PhAt />
            dge@agros.org.br
          </p>
        </li>
      </ul>
    </section>
  </main>
</template>

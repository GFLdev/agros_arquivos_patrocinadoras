<script setup lang="ts">
import Header from '@/components/generic/HeaderSection.vue'
import { PhEnvelopeSimple, PhMapPin, PhPhone } from '@phosphor-icons/vue'
</script>

<template>
  <Header />
  <main
    class="flex w-screen flex-col items-center justify-center gap-8 bg-gray bg-opacity-5 px-8 py-4 font-lato sm:px-16 sm:py-8"
  >
    <section class="w-fit max-w-3xl rounded-lg bg-white px-10 py-8 text-center shadow-xl drop-shadow-xl sm:text-lg">
      <p>
        Nesta página estão disponíveis informações sobre os planos de previdência e de saúde do Agros, em atendimento à
        Resolução Normativa 389, da Agência Nacional de Saúde Suplementar (ANS).
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
    <!-- Divisor -->
    <div class="relative my-8 w-full">
      <div class="absolute inset-0 flex items-center" aria-hidden="true">
        <div class="w-full border-t border-gray border-opacity-30" />
      </div>
      <div class="relative flex justify-center">
        <span class="bg-[ text-gray-500 bg-gray px-2 text-sm">Continue</span>
      </div>
    </div>
    <!-- Seção de dúvidas -->
    <section class="relative rounded-lg bg-white px-10 py-8 shadow-xl drop-shadow-xl lg:static">
      <div class="w-full max-w-xl text-center">
        <!-- Heading -->
        <h2 class="text-pretty text-2xl font-semibold tracking-tight text-dark">Dúvidas?</h2>
        <!-- Descrição -->
        <p class="mt-6 text-dark sm:text-lg/8">
          Em caso de dúvida sobre as informações aqui apresentadas, entre em contato com o Agros.
        </p>
        <!-- Itens -->
        <dl
          class="mt-6 grid grid-cols-1 space-y-4 text-dark sm:grid-cols-3 sm:items-start sm:justify-between sm:space-y-0 sm:text-base/7"
        >
          <div
            class="flex items-center justify-center gap-4 transition-colors duration-200 hover:text-primary sm:flex-col sm:gap-2"
          >
            <dt class="flex-none">
              <span class="sr-only">Address</span>
              <PhMapPin class="size-6" aria-hidden="true" />
            </dt>
            <dd>Av. Purdue, s/n<br />Campus da UFV - Viçosa</dd>
          </div>
          <div
            class="flex items-center justify-center gap-4 transition-colors duration-200 hover:text-primary sm:flex-col sm:gap-2"
          >
            <dt class="flex-none">
              <span class="sr-only">Telephone</span>
              <PhPhone class="size-6" aria-hidden="true" />
            </dt>
            <dd>
              <a>(31) 3899-6550<br />ramal 6539</a>
            </dd>
          </div>
          <div
            class="flex items-center justify-center gap-4 transition-colors duration-200 hover:text-primary sm:flex-col sm:gap-2"
          >
            <dt class="flex-none">
              <span class="sr-only">Email</span>
              <PhEnvelopeSimple class="size-6" aria-hidden="true" />
            </dt>
            <dd>
              <a href="mailto:dge@agros.org.br">dge@agros.org.br</a>
            </dd>
          </div>
        </dl>
      </div>
    </section>
  </main>
</template>

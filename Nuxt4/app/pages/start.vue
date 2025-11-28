<!-- eslint-disable vue/multi-word-component-names -->
<script setup>
const slides = [
  {
    id: 1,
    title: 'Привет, это пробный проект',
    description: 'Начало пути',
  },
  {
    id: 2,
    title: 'Создан чтобы был',
    description: 'Не люблю писать',
  },
  {
    id: 3,
    title: 'Вот хоть котиков нате',
    description: 'Завершение.',
    image:
      'https://st4.depositphotos.com/38981760/39647/v/450/depositphotos_396475506-stock-illustration-cat-illustration-set-outline-silhouette.jpg',
  },
]
</script>

<template>
  <div>
    <UCarousel
      orientation="vertical"
      :items="slides"
      wheel-gestures
      :ui="{ item: 'basis-full', container: 'h-screen' }"
      class="h-full"
    >
      <template #default="{ item }">
        <div
          class="h-screen w-screen flex flex-col items-center justify-center p-4"
          :style="{ backgroundImage: `url(${item.image})` }"
        >
          <h1 class="text-4xl font-bold mb-4">{{ item.title }}</h1>
          <div v-if="item.description === 'Завершение.'" class="mt-8">
            <UButton size="xl" to="/"> Начать </UButton>
          </div>
          <p v-if="!item.description === 'Завершение.'" class="text-lg">
            {{ item.description }}
          </p>
        </div>
      </template>
    </UCarousel>
  </div>
</template>

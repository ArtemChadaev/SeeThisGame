<script setup lang="ts">
import type { NavigationMenuItem } from '@nuxt/ui'

const items = ref<NavigationMenuItem[]>([
  [
    {
      label: 'Настройки',
      icon: 'i-lucide-settings',
      children: [
        {
          label: 'Основные',
          icon: 'i-lucide-sliders-horizontal',
          to: 'general'
        },
        {
          label: 'Клан',
          icon: 'i-lucide-flag',
          to: 'clan'
        },
        {
          label: 'Безопасность',
          icon: 'i-lucide-brick-wall-fire',
        },
        {
          label: 'Устройства',
          icon: 'i-lucide-cog',
        }
      ]
    },

    {
      label: 'карточки',
      icon: 'i-lucide-id-card-lanyard',
      children: [
        {
          label: 'Основные',
          icon: 'i-lucide-drama',
        },
        {
          label: 'Любимые',
          icon: 'i-lucide-heart',
        },
        {
          label: 'Предметы',
          icon: 'i-lucide-tool-case',
        },
        {
          label: 'Статистика',
          icon: 'i-lucide-chart-column-stacked',
        },
      ]
    },
    {
    label: 'Игры',
    icon: 'i-lucide-gamepad-2',
    children: [
      {
        label: 'Руководство',
        icon: 'i-lucide-book',
        to: ''
      },
      {
        label: 'Достижения',
        icon: 'i-lucide-file-badge',
      },
      {
        label: 'Управление',
        icon: 'i-lucide-keyboard',
      },
      {
        label: 'Предложения',
        icon: 'i-lucide-notebook',
      }
    ]
  },
  ],
  [
    {
      label: 'GitHub',
      icon: 'i-simple-icons-github',
      to: 'https://github.com/ChadaevArtem/nuxt',
      target: '_blank'
    }
  ]
])
</script>

<template>
  <UContainer class="flex -mx-10">
    <UNavigationMenu orientation="vertical" :items="items" class="w-50" />
    <NuxtPage />
  </UContainer>
</template>

<script setup lang="ts">
import type { TabsItem } from '@nuxt/ui'
import RegisterForm from '~/pages/login/RegisterForm.vue'
import LoginForm from '~/pages/login/LoginForm.vue'

const route = useRoute()
const userStore = useUserStore()
const goBack = () => {

  userStore.getMySettings()

  const page =
    typeof route.query.previous === 'string' ? route.query.previous : '/'

  navigateTo(page)
}
const itemsTabs = ref<TabsItem[]>([
  {
    label: 'login',
    slot: 'login' as const,
  },
  {
    label: 'register',
    slot: 'register' as const,
  },
])
</script>

<template>
  <UTabs :items="itemsTabs" class="w-full">
    <template #login>
      <LoginForm @complete="goBack" />
    </template>

    <template #register>
      <RegisterForm @complete="goBack" />
    </template>
  </UTabs>
</template>

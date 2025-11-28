<template>
  <UApp class="w-screen h-screen bg-neutral-50 dark:bg-neutral-950">
    <UContainer class="py-2">
      <AppHeader />
    </UContainer>

    <UContainer class="mt-8 max-w-6xl">
      <UCard>
        <slot />
      </UCard>
    </UContainer>
  </UApp>
</template>
<script setup lang="ts">
// TODO: Сделать чёрную тему

const tokenStore = useTokenStore()
const userStore = useUserStore()
onMounted(() => {
  // Не вошли нет id почему-то, но сам токен есть
  if (!userStore.isLogged && !!tokenStore.refreshToken)
    userStore.getMySettings()

})

</script>

<script setup lang="ts">
import { useClipboard } from '@vueuse/core'
const userStore = useUserStore()

const userid = userStore.id?.toString() ?? ''
const {copy} = useClipboard({source: userid})

const state = {
  name: undefined,
}
//TODO: Разобраться и вспомнить что я намудрил
// const onSubmit = async () => {
//   try {
//     const response = await $fetch<{
//       refreshToken: string
//       accessToken: string
//     }>('/auth/sign-in', {
//       method: 'POST',
//       body: event.data,
//     })
//
//     tokenStore.setTokens(response.refreshToken, response.accessToken)
//
//     toast.add({
//       title: 'Успех',
//       description: 'Вы успешно вошли в систему',
//       color: 'success',
//     })
//
//     emit('complete')
//   } catch (error: unknown) {
//     console.error(error)
//     if (error instanceof FetchError) {
//       toast.add({
//         title: 'Ошибка',
//         description: error.data?.error || 'Ошибка авторизации',
//         color: 'error',
//       })
//     }
//
//   }
// }
</script>

<template>
  <UContainer class="m-2 ml-12">
    <UForm
      class="space-y-2 w-full flex flex-col gap-2"
      @submit="onSubmit"
    >
      <UAvatar :src="!!userStore.icon ? userStore.icon : 'https://github.com/benjamincanac.png'" size="3xl" class="self-center mb-0" />
      <p class="text-sm text-gray-400 cursor-pointer self-center hover:text-info" @click="() => {copy()}">id:{{userStore.id}}</p>
      <UInput v-model="state.name" :placeholder="!!userStore.name ? userStore.name : 'Имя'" />
      <UButton label="Сохранить" type="submit" class="mt-2 self-end" />
    </UForm>
  </UContainer>
</template>

<style scoped></style>

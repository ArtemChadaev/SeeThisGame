<script setup lang="ts">
import * as z from 'zod'
import type { FormSubmitEvent } from '@nuxt/ui'
import { FetchError } from 'ofetch'

const tokenStore = useTokenStore()

const emit = defineEmits(['complete'])

const schema = z.object({
  email: z.string().email('Некорректный email'),
  password: z
    .string()
    .min(8, 'Минимум 8 символов')
    .max(32, 'Максимум 32 символа'),
})

type Schema = z.output<typeof schema>

const state = reactive<Partial<Schema>>({
  email: undefined,
  password: undefined,
})

const toast = useToast()

const onSubmit = async (event: FormSubmitEvent<Schema>) => {
  try {
    const response = await $fetch<{
      refreshToken: string
      accessToken: string
    }>('/auth/sign-in', {
      method: 'POST',
      body: event.data,
    })

    tokenStore.setTokens(response.refreshToken, response.accessToken)

    toast.add({
      title: 'Успех',
      description: 'Вы успешно вошли в систему',
      color: 'success',
    })

    emit('complete')
  } catch (error: unknown) {
    console.error(error)
    if (error instanceof FetchError) {
      toast.add({
        title: 'Ошибка',
        description: error.data?.error || 'Ошибка авторизации',
        color: 'error',
      })
    }
  }
}
</script>

<template>
  <UForm
    :schema="schema"
    :state="state"
    class="space-y-2 w-full flex flex-col gap-2"
    @submit="onSubmit"
  >
    <UFormField label="Email" name="email">
      <UInput v-model="state.email" placeholder="Email" class="w-full" />
    </UFormField>
    <UFormField label="Password" name="password">
      <UInput
        v-model="state.password"
        placeholder="Password"
        type="password"
        class="w-full mb-[72px]"
      />
    </UFormField>
    <UButton label="login" type="submit" class="mt-2 self-end" size="lg" />
  </UForm>
</template>

<style scoped></style>

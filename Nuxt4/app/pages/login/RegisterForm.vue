<script setup lang="ts">
import type { FormSubmitEvent } from '@nuxt/ui'
import * as z from 'zod'
import { FetchError } from 'ofetch'

const tokenStore = useTokenStore()
const emit = defineEmits(['complete'])
const toast = useToast()

const schema = z
  .object({
    email: z.string().email('Некорректный email'),
    password: z
      .string()
      .min(8, 'Минимум 8 символов')
      .max(32, 'Максимум 32 символа'),
    passwordConfirm: z.string(),
  })
  .refine((data) => data.password === data.passwordConfirm, {
    message: 'пароли не совпадают',
    path: ['passwordConfirm'],
  })

type Schema = z.output<typeof schema>

const state = reactive({
  email: '',
  password: '',
  passwordConfirm: '',
})

async function onSubmit(event: FormSubmitEvent<Schema>) {
  try {
    const response = await $fetch<{
      refreshToken: string
      accessToken: string
    }>('/auth/sign-up', {
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
    console.log(error)
    if (error instanceof FetchError) {
      if (
        error.response &&
        error.response.status === 409 &&
        error.data.error === 'email_exist'
      ) {
        toast.add({
          title: 'Ошибка',
          description: 'email уже используется',
          color: 'error',
        })
        return
      }
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
        class="w-full"
      />
    </UFormField>
    <UFormField label="PasswordConfirm" name="passwordConfirm">
      <UInput
        v-model="state.passwordConfirm"
        placeholder="PasswordConfirm"
        type="password"
        class="w-full"
      />
    </UFormField>
    <UButton type="submit" label="register" class="mt-2 self-end" size="lg" />
  </UForm>
</template>

<style scoped></style>

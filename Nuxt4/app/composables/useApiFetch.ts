// Получаем тип для опций напрямую из глобального типа $Fetch
type ApiFetchOptions = Parameters<typeof $fetch>[1];
// TODO: Сделать блок при перенаправлении с кода 400 invalid-token и других
export const useApiFetch = async <T>(url: string, options: ApiFetchOptions = {}) => {
  const tokenStore = useTokenStore();
  const isAuthRoute = url.includes('/auth/');

  if (!isAuthRoute && tokenStore.refreshToken) {
    if (tokenStore.isTokenExpired || !tokenStore.accessToken) {
      try {
        await tokenStore.updateToken();
      } catch (error) {
        console.error('Критическая ошибка: не удалось обновить токен.', error);
        throw new Error('Не удалось выполнить запрос из-за ошибки аутентификации.');
      }
    }

    const headers = new Headers(options.headers);

    if (tokenStore.accessToken) {
      headers.set('Authorization', `Bearer ${tokenStore.accessToken}`);
    }

    options.headers = headers;
  }

  // Теперь ошибки точно не будет, так как мы используем родной тип Nuxt
  return $fetch<T>(url, options);
};
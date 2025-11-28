# Dockerfile

# --- Стадия сборки ---
# Используем официальный образ Bun как основу для сборки
FROM oven/bun:1.2.22-alpine AS build

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем файлы package.json и bun.lock для установки зависимостей
COPY package.json bun.lock ./

# Устанавливаем зависимости, используя --frozen-lockfile для строгого соответствия lock-файлу
RUN bun install --frozen-lockfile

# Копируем остальные файлы проекта в рабочую директорию
COPY . .

# Собираем Nuxt приложение для продакшена
# Команда 'build' берется из твоего package.json
RUN bun run build

# --- Стадия продакшена ---
# Используем легковесный образ Node.js Alpine для запуска приложения
FROM node:20-alpine AS production

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем собранный проект (директорию .output) из стадии сборки
COPY --from=build /app/.output ./.output

# Устанавливаем переменные окружения для Nuxt сервера
# HOST=0.0.0.0 нужен для доступности сервера извне контейнера
ENV HOST=0.0.0.0
# Стандартный порт Nuxt
ENV PORT=3000
# Ты можешь добавить сюда другие переменные окружения, необходимые твоему приложению
# Например: ENV NUXT_PUBLIC_API_URL=https://api.example.com

# Сообщаем Docker, что контейнер будет слушать указанный порт
EXPOSE 3000

# Команда для запуска Nuxt сервера в продакшен режиме
# Nuxt 3/4 использует .output/server/index.mjs как точку входа
CMD ["node", ".output/server/index.mjs"]
# Интеграция с ИИ через n8n

## 1. Цель

Автоматизация процессов генерации контента и обработки данных с помощью различных ИИ-моделей. n8n выступает в роли связующего звена (middleware) между бэкендом приложения и ИИ-сервисами.

## 2. Схема работы

1.  **Бэкенд (Go)** отправляет HTTP-запрос (webhook) в n8n с данными для обработки.
2.  **n8n** получает запрос и запускает соответствующий рабочий процесс (workflow).
3.  **Рабочий процесс n8n:**
    - Обрабатывает входящие данные.
    - Формирует запрос для ИИ-модели.
    - Отправляет запрос к соответствующему ИИ-сервису.
    - Получает результат от ИИ.
    - (Опционально) Постобрабатывает результат.
    - Отправляет результат обратно на бэкенд или сохраняет в хранилище.

### 2.1. API n8n для взаимодействия с бэкендом

n8n предоставляет несколько способов взаимодействия:

#### Webhook Endpoints (Production URL)
```
POST https://your-n8n-instance.com/webhook/{webhook-path}
POST https://your-n8n-instance.com/webhook-test/{webhook-path}
```

#### n8n REST API
```
Base URL: https://your-n8n-instance.com/api/v1/
Authentication: X-N8N-API-KEY: your_api_key
```

**Основные эндпоинты:**

1. **Запуск workflow через webhook:**
   ```
   POST /webhook/generate-image
   Content-Type: application/json
   
   {
     "prompt": "fantasy character, warrior",
     "model": "stable-diffusion",
     "size": "1024x1024"
   }
   ```

2. **Запуск workflow через API (с получением execution ID):**
   ```
   POST /api/v1/workflows/{workflowId}/execute
   X-N8N-API-KEY: your_api_key
   Content-Type: application/json
   
   {
     "data": {
       "prompt": "fantasy character, warrior"
     }
   }
   
   Response:
   {
     "data": {
       "executionId": "12345"
     }
   }
   ```

3. **Проверка статуса выполнения:**
   ```
   GET /api/v1/executions/{executionId}
   X-N8N-API-KEY: your_api_key
   
   Response:
   {
     "data": {
       "id": "12345",
       "finished": true,
       "mode": "webhook",
       "status": "success",
       "data": {
         "resultData": {
           "runData": {...}
         }
       }
     }
   }
   ```

4. **Получение результата выполнения:**
   ```
   GET /api/v1/executions/{executionId}/results
   X-N8N-API-KEY: your_api_key
   ```

#### Пример асинхронной работы с ComfyUI

**Workflow в n8n для генерации изображения:**

1. **Webhook Trigger** → получает запрос от бэкенда
2. **ComfyUI: Submit Job** → отправляет задачу в ComfyUI
   ```
   POST http://localhost:8188/prompt
   {
     "prompt": {...workflow_definition...}
   }
   Response: {"prompt_id": "abc123"}
   ```

3. **Set Variable** → сохраняет `prompt_id`
4. **Wait** → ждет 2-5 секунд (опционально)
5. **Loop: Check Status** → проверяет готовность
   ```
   GET http://localhost:8188/history/{prompt_id}
   ```

6. **Condition** → если готово, переходит к следующему шагу
7. **ComfyUI: Get Result** → получает изображение
   ```
   GET http://localhost:8188/view?filename={filename}
   ```

8. **Upload to S3** → загружает в хранилище
9. **Webhook Response** → возвращает URL изображения

**Два подхода к асинхронности:**

**Подход 1: Синхронный (n8n ждет результата)**
```javascript
// Go Backend
response := httpClient.Post("https://n8n.com/webhook/generate-image", data)
imageUrl := response.Data.ImageUrl // получаем готовый результат
```

**Подход 2: Асинхронный (с callback)**
```javascript
// Go Backend - Шаг 1: Запуск задачи
response := httpClient.Post("https://n8n.com/webhook/generate-image-async", {
  "prompt": "warrior",
  "callback_url": "https://backend.com/api/image-ready"
})
taskId := response.Data.TaskId

// n8n workflow отправит результат на callback_url когда готово
// Go Backend - Шаг 2: Получение результата через callback
func HandleImageReady(w http.ResponseWriter, r *http.Request) {
  var result ImageResult
  json.NewDecoder(r.Body).Decode(&result)
  // Сохраняем результат в БД
  db.SaveImage(result.TaskId, result.ImageUrl)
}
```

**Подход 3: Polling (опрос статуса)**
```javascript
// Go Backend - Шаг 1: Запуск задачи
response := httpClient.Post("https://n8n.com/api/v1/workflows/123/execute", data)
executionId := response.Data.ExecutionId

// Шаг 2: Периодическая проверка статуса
for {
  status := httpClient.Get("https://n8n.com/api/v1/executions/" + executionId)
  if status.Data.Finished {
    result := httpClient.Get("https://n8n.com/api/v1/executions/" + executionId + "/results")
    break
  }
  time.Sleep(5 * time.Second)
}
```

#### Пример workflow с несколькими API для одной функции

**Функция: Генерация изображения персонажа с проверкой качества**

```
1. Webhook Trigger
   ↓
2. ComfyUI: Generate Image (POST /prompt)
   → Получаем prompt_id
   ↓
3. Loop: Check ComfyUI Status (GET /history/{prompt_id})
   → Ждем пока status = "completed"
   ↓
4. ComfyUI: Download Image (GET /view?filename=...)
   ↓
5. OpenAI Vision: Analyze Quality (POST /v1/chat/completions)
   → Проверяем соответствие промпту
   ↓
6. Condition: Quality Check
   → Если качество низкое, возвращаемся к шагу 2
   → Если качество хорошее, продолжаем
   ↓
7. Stability AI: Upscale (POST /v1/generation/esrgan-v1-x2plus/image-to-image/upscale)
   ↓
8. MinIO: Upload to S3 (PUT /bucket/image.png)
   ↓
9. PostgreSQL: Save Metadata
   ↓
10. Webhook Response / Callback
```

## 3. Функции ИИ и соответствующие API

### 3.1. Генерация изображений

**Функции:**
- Генерация изображений персонажей по описанию
- Создание изображений локаций и окружения
- Генерация событий и сцен
- Создание UI элементов и иконок
- Генерация вариаций существующих изображений
- Upscaling и улучшение качества изображений
- Inpainting и outpainting изображений

**API:**
- **[OpenAI DALL-E 3 API](https://platform.openai.com/docs/api-reference/images)** - `https://api.openai.com/v1/images/generations`
- **[Stability AI API](https://platform.stability.ai/docs/api-reference)** - `https://api.stability.ai/v1/generation/`
- **[Midjourney API (неофициальный)](https://github.com/erictik/midjourney-api)** - через прокси-сервис
- **[ComfyUI API](https://github.com/comfyanonymous/ComfyUI)** - `http://localhost:8188/` (локальный)
- **[Replicate API](https://replicate.com/docs/reference/http)** - `https://api.replicate.com/v1/predictions`
- **[Leonardo.AI API](https://docs.leonardo.ai/reference/introduction)** - `https://cloud.leonardo.ai/api/rest/v1/`

### 3.2. Генерация текста и диалогов

**Функции:**
- Генерация диалогов NPC
- Создание описаний предметов и локаций
- Генерация квестов и сюжетных линий
- Создание имен персонажей
- Генерация лора и истории мира
- Перевод текстов
- Улучшение и редактирование текстов

**API:**
- **[OpenAI GPT-4 API](https://platform.openai.com/docs/api-reference/chat)** - `https://api.openai.com/v1/chat/completions`
- **[Anthropic Claude API](https://docs.anthropic.com/claude/reference/getting-started-with-the-api)** - `https://api.anthropic.com/v1/messages`
- **[Google Gemini API](https://ai.google.dev/api/rest)** - `https://generativelanguage.googleapis.com/v1/models/`
- **[Cohere API](https://docs.cohere.com/reference/chat)** - `https://api.cohere.ai/v1/chat`
- **[Hugging Face Inference API](https://huggingface.co/docs/api-inference/index)** - `https://api-inference.huggingface.co/models/`
- **[Together AI API](https://docs.together.ai/reference/inference)** - `https://api.together.xyz/v1/chat/completions`
- **[Ollama API](https://github.com/ollama/ollama/blob/main/docs/api.md)** - `http://localhost:11434/api/` (локальный)

### 3.3. Генерация и обработка голоса

**Функции:**
- Синтез речи для озвучки диалогов
- Клонирование голосов персонажей
- Преобразование текста в речь (TTS)
- Распознавание речи (STT)
- Генерация звуковых эффектов

**API:**
- **[ElevenLabs API](https://elevenlabs.io/docs/api-reference/text-to-speech)** - `https://api.elevenlabs.io/v1/text-to-speech/`
- **[OpenAI TTS API](https://platform.openai.com/docs/api-reference/audio/createSpeech)** - `https://api.openai.com/v1/audio/speech`
- **[OpenAI Whisper API](https://platform.openai.com/docs/api-reference/audio/createTranscription)** - `https://api.openai.com/v1/audio/transcriptions`
- **[Google Cloud Text-to-Speech](https://cloud.google.com/text-to-speech/docs/reference/rest)** - `https://texttospeech.googleapis.com/v1/text:synthesize`
- **[Azure Speech Services](https://learn.microsoft.com/en-us/azure/ai-services/speech-service/rest-text-to-speech)** - `https://<region>.tts.speech.microsoft.com/cognitiveservices/v1`
- **[Play.ht API](https://docs.play.ht/reference/api-getting-started)** - `https://api.play.ht/api/v2/tts`

### 3.4. Генерация музыки и звуков

**Функции:**
- Генерация фоновой музыки
- Создание звуковых эффектов
- Генерация амбиентных звуков

**API:**
- **[Suno AI API](https://suno.com/)** - через неофициальные обертки
- **[Mubert API](https://docs.mubert.com/api/)** - `https://api-b2b.mubert.com/v2/`
- **[Soundraw API](https://soundraw.io/)** - коммерческий доступ
- **[AudioCraft (Meta)](https://github.com/facebookresearch/audiocraft)** - локальное развертывание

### 3.5. Модерация контента

**Функции:**
- Проверка текста на токсичность
- Модерация изображений
- Фильтрация неприемлемого контента
- Детекция NSFW контента

**API:**
- **[OpenAI Moderation API](https://platform.openai.com/docs/api-reference/moderations)** - `https://api.openai.com/v1/moderations`
- **[Perspective API (Google)](https://developers.perspectiveapi.com/s/)** - `https://commentanalyzer.googleapis.com/v1alpha1/comments:analyze`
- **[Azure Content Safety](https://learn.microsoft.com/en-us/azure/ai-services/content-safety/)** - `https://<endpoint>.cognitiveservices.azure.com/contentsafety/`
- **[Sightengine API](https://sightengine.com/docs/)** - `https://api.sightengine.com/1.0/check.json`

### 3.6. Embeddings и семантический поиск

**Функции:**
- Создание векторных представлений текста
- Семантический поиск по контенту
- Кластеризация контента
- Поиск похожих элементов

**API:**
- **[OpenAI Embeddings API](https://platform.openai.com/docs/api-reference/embeddings)** - `https://api.openai.com/v1/embeddings`
- **[Cohere Embed API](https://docs.cohere.com/reference/embed)** - `https://api.cohere.ai/v1/embed`
- **[Voyage AI API](https://docs.voyageai.com/reference/embeddings-api)** - `https://api.voyageai.com/v1/embeddings`
- **[Hugging Face Sentence Transformers](https://huggingface.co/sentence-transformers)** - через Inference API

### 3.7. Анализ изображений и Vision

**Функции:**
- Описание изображений
- Распознавание объектов на изображениях
- OCR (распознавание текста)
- Анализ сцен и контекста

**API:**
- **[OpenAI GPT-4 Vision API](https://platform.openai.com/docs/guides/vision)** - `https://api.openai.com/v1/chat/completions`
- **[Anthropic Claude Vision](https://docs.anthropic.com/claude/docs/vision)** - `https://api.anthropic.com/v1/messages`
- **[Google Cloud Vision API](https://cloud.google.com/vision/docs/reference/rest)** - `https://vision.googleapis.com/v1/images:annotate`
- **[Azure Computer Vision](https://learn.microsoft.com/en-us/azure/ai-services/computer-vision/)** - `https://<endpoint>.cognitiveservices.azure.com/vision/`

### 3.8. Fine-tuning и кастомизация моделей

**Функции:**
- Дообучение моделей на специфичных данных
- Создание кастомных моделей для игры
- Адаптация моделей под стиль игры

**API:**
- **[OpenAI Fine-tuning API](https://platform.openai.com/docs/api-reference/fine-tuning)** - `https://api.openai.com/v1/fine_tuning/jobs`
- **[Hugging Face AutoTrain](https://huggingface.co/autotrain)** - веб-интерфейс и API
- **[Replicate Training API](https://replicate.com/docs/guides/fine-tune-a-language-model)** - `https://api.replicate.com/v1/trainings`

## 4. Используемые технологии

- **Платформа автоматизации:** [n8n](https://n8n.io/)
- **Векторная БД:** [Qdrant](https://qdrant.tech/) или [Pinecone](https://www.pinecone.io/) для хранения embeddings
- **Хранилище:** S3-совместимое хранилище (например, MinIO) для сгенерированных файлов
- **Кэширование:** Redis для кэширования результатов ИИ-запросов
- **Очереди:** RabbitMQ или Redis для управления очередью задач

## 5. Этапы разработки

1.  **Развертывание инфраструктуры:**
    - Установка n8n (через Docker)
    - Настройка локальных ИИ-сервисов (ComfyUI, Ollama)
    - Развертывание хранилища и БД

2.  **Создание рабочих процессов (Workflows) в n8n:**
    - Настройка Webhook-триггеров для каждого типа запросов
    - Создание узлов для обработки данных и формирования промптов
    - Настройка HTTP-узлов для отправки запросов к различным API
    - Добавление узлов для сохранения результатов
    - Настройка обработки ошибок и retry-логики

3.  **Интеграция с API:**
    - Настройка аутентификации для каждого API
    - Создание переиспользуемых credentials в n8n
    - Настройка rate limiting и квот

4.  **Тестирование:**
    - Отладка каждого workflow с тестовыми данными
    - Нагрузочное тестирование
    - Тестирование обработки ошибок

5.  **Интеграция с бэкендом:**
    - Настройка отправки запросов из Go-приложения
    - Реализация асинхронной обработки
    - Настройка webhook'ов для получения результатов

6.  **Оптимизация и масштабирование:**
    - Настройка кэширования результатов
    - Оптимизация промптов для снижения затрат
    - Балансировка нагрузки между различными API
    - Мониторинг и логирование

## 6. Безопасность и best practices

- Хранение API ключей в переменных окружения
- Использование rate limiting для предотвращения превышения квот
- Валидация входящих данных
- Модерация сгенерированного контента
- Логирование всех запросов для аудита
- Резервное копирование workflows

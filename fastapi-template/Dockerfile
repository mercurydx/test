# ==========================================
# ЭТАП 1: BUILDER (Сборка зависимостей)
# ==========================================
FROM python:3.11-slim-bullseye AS builder

# Устанавливаем переменные окружения для Python
ENV PYTHONDONTWRITEBYTECODE=1 \
    PYTHONUNBUFFERED=1 \
    PIP_NO_CACHE_DIR=off \
    PIP_DISABLE_PIP_VERSION_CHECK=on

# Устанавливаем системные зависимости, нужные ТОЛЬКО для сборки
# libpq-dev нужен для компиляции драйвера psycopg2 (PostgreSQL)
RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    libpq-dev \
    gcc \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /usr/src/app

# 1. Сначала копируем файлы зависимостей (чтобы использовать кэш Docker)
# Если Pipfile не менялся, Docker пропустит установку при следующих сборках!
COPY Pipfile Pipfile.lock ./

# Устанавливаем pipenv и зависимости в системный Python
RUN pip install --upgrade pip && \
    pip install pipenv && \
    pipenv install --system --deploy --ignore-pipfile

# 2. Только теперь копируем исходный код приложения
COPY . .

# ==========================================
# ЭТАП 2: RUNTIME (Финальный легкий образ)
# ==========================================
FROM python:3.11-slim-bullseye AS runtime

ENV PYTHONDONTWRITEBYTECODE=1 \
    PYTHONUNBUFFERED=1

# Устанавливаем ТОЛЬКО runtime-зависимости (без gcc и build-essential!)
# libpq5 - это скомпилированная библиотека для работы с Postgres (весит копейки)
RUN apt-get update && apt-get install -y --no-install-recommends \
    libpq5 \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /usr/src/app

# Создаем непривилегированного пользователя для безопасности (Best Practice!)
RUN useradd -m -u 1000 appuser

# Копируем установленные пакеты и код из этапа builder
COPY --from=builder /usr/local/lib/python3.11/site-packages /usr/local/lib/python3.11/site-packages
COPY --from=builder /usr/local/bin /usr/local/bin
COPY --from=builder /usr/src/app /usr/src/app

# Меняем владельца файлов на нашего нового пользователя
RUN chown -R appuser:appuser /usr/src/app

# Переключаемся на непривилегированного пользователя
USER appuser

# Открываем порт, который слушает приложение (в docker-compose указан 11000)
EXPOSE 11000

# Команда запуска (в идеале здесь должен быть uvicorn, но оставим как в оригинале)
# Если main.py содержит uvicorn.run(), то это сработает.
CMD ["python", "main.py"]
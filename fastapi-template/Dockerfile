# ==========================================
# ЭТАП 1: BUILDER (Сборка зависимостей)
# ==========================================
FROM python:3.11-slim-bullseye AS builder

ENV PYTHONDONTWRITEBYTECODE=1 \
    PYTHONUNBUFFERED=1 \
    PIP_NO_CACHE_DIR=off \
    PIP_DISABLE_PIP_VERSION_CHECK=on

# Устанавливаем системные зависимости для сборки
RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    libpq-dev \
    gcc \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /usr/src/app

# Копируем файлы зависимостей
COPY Pipfile Pipfile.lock ./

# Экспортируем зависимости из Pipfile в requirements.txt
RUN pip install --upgrade pip && \
    pip install pipenv && \
    pipenv requirements > requirements.txt

# Устанавливаем зависимости через pip (более надежно)
RUN pip install --no-cache-dir -r requirements.txt

# Копируем исходный код приложения
COPY . .

# ==========================================
# ЭТАП 2: RUNTIME (Финальный легкий образ)
# ==========================================
FROM python:3.11-slim-bullseye AS runtime

ENV PYTHONDONTWRITEBYTECODE=1 \
    PYTHONUNBUFFERED=1

# Устанавливаем ТОЛЬКО runtime-зависимости
RUN apt-get update && apt-get install -y --no-install-recommends \
    libpq5 \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /usr/src/app

# Создаем непривилегированного пользователя
RUN useradd -m -u 1000 appuser

# Копируем установленные пакеты из builder
COPY --from=builder /usr/local/lib/python3.11/site-packages /usr/local/lib/python3.11/site-packages
COPY --from=builder /usr/local/bin /usr/local/bin
COPY --from=builder /usr/src/app /usr/src/app

# Меняем владельца файлов
RUN chown -R appuser:appuser /usr/src/app

# Переключаемся на непривилегированного пользователя
USER appuser

EXPOSE 11000

CMD ["python", "main.py"]
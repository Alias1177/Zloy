-- Migration: 001_create_users_table.sql
-- Description: Создание таблицы пользователей

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    login VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    balance INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Создаем индекс для быстрого поиска по логину
CREATE INDEX IF NOT EXISTS idx_users_login ON users(login);

-- Создаем индекс для поиска по дате создания
CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at); 
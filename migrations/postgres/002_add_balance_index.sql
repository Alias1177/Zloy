-- Migration: 002_add_balance_index.sql
-- Description: Добавление индекса для баланса пользователей

CREATE INDEX IF NOT EXISTS idx_users_balance ON users(balance); 
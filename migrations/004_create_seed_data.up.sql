
-- db/migrations/YYYYMMDDHHMMSS_create_seed_data.up.sql

-- Вставка организаций
INSERT INTO organization (id, name, description, type)
VALUES ('org1_id', 'ООО Ромашка', 'Производство цветов', 'LLC')
ON CONFLICT (id) DO NOTHING; -- Идемпотентность: пропустить, если уже существует

INSERT INTO organization (id, name, description, type)
VALUES ('org2_id', 'ИП Петров', 'Ремонт квартир', 'IE')
ON CONFLICT (id) DO NOTHING;

-- Вставка пользователей
INSERT INTO employee (id, username, first_name, last_name)
VALUES ('user1_id', 'ivanov', 'Иван', 'Иванов')
ON CONFLICT (id) DO NOTHING;

INSERT INTO employee (id, username, first_name, last_name)
VALUES ('user2_id', 'petrov', 'Петр', 'Петров')
ON CONFLICT (id) DO NOTHING;

INSERT INTO employee (id, username, first_name, last_name)
VALUES ('user3_id', 'sidorov', 'Сидор', 'Сидоров')
ON CONFLICT (id) DO NOTHING;

-- Связь ответственных за организации
INSERT INTO organization_responsible (id, organization_id, user_id)
VALUES ('resp1_id', 'org1_id', 'user1_id')
ON CONFLICT (id) DO NOTHING;

INSERT INTO organization_responsible (id, organization_id, user_id)
VALUES ('resp2_id', 'org2_id', 'user2_id')
ON CONFLICT (id) DO NOTHING;
CREATE TABLE quotes (
    id SERIAL PRIMARY KEY,
    text TEXT NOT NULL
);

INSERT INTO quotes (text) VALUES 
('🚀 Полетели!'),
('🎉 Отличный день для кода!'),
('☕ Время для кофе и curl'),
('🔥 Сервер горяч, как мой процессор'),
('🐳 Docker и Kubernetes — наше всё');
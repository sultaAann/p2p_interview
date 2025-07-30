-- Создаем таблицу slots
CREATE TABLE slots (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    interviewer_id UUID NOT NULL,
    interviewee_id UUID,
    start_time TIMESTAMP WITH TIME ZONE NOT NULL,
    end_time TIMESTAMP WITH TIME ZONE NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'available',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- Ограничения
    CONSTRAINT fk_slots_interviewer FOREIGN KEY (interviewer_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_slots_interviewee FOREIGN KEY (interviewee_id) REFERENCES users(id) ON DELETE SET NULL,
    CONSTRAINT chk_slots_time_order CHECK (start_time < end_time),
    CONSTRAINT chk_slots_status CHECK (status IN ('available', 'booked', 'completed', 'cancelled'))
);

-- Создаем индексы для быстрого поиска
CREATE INDEX idx_slots_interviewer_id ON slots(interviewer_id);
CREATE INDEX idx_slots_interviewee_id ON slots(interviewee_id);
CREATE INDEX idx_slots_start_time ON slots(start_time);
CREATE INDEX idx_slots_end_time ON slots(end_time);
CREATE INDEX idx_slots_status ON slots(status);
CREATE INDEX idx_slots_created_at ON slots(created_at);

-- Составной индекс для поиска доступных слотов по времени
CREATE INDEX idx_slots_available_time ON slots(status, start_time, end_time) 
WHERE status = 'available';

-- Составной индекс для поиска слотов интервьюера
CREATE INDEX idx_slots_interviewer_status ON slots(interviewer_id, status);

-- Создаем триггер для автоматического обновления updated_at при изменении записи
CREATE TRIGGER update_slots_updated_at
    BEFORE UPDATE ON slots
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Добавляем комментарии к таблице и столбцам для документации
COMMENT ON TABLE slots IS 'Таблица временных слотов для интервью';
COMMENT ON COLUMN slots.id IS 'Уникальный идентификатор слота (UUID)';
COMMENT ON COLUMN slots.interviewer_id IS 'ID интервьюера (ссылка на users.id)';
COMMENT ON COLUMN slots.interviewee_id IS 'ID интервьюируемого (ссылка на users.id), может быть NULL для свободных слотов';
COMMENT ON COLUMN slots.start_time IS 'Время начала слота';
COMMENT ON COLUMN slots.end_time IS 'Время окончания слота';
COMMENT ON COLUMN slots.status IS 'Статус слота: available, booked, completed, cancelled';
COMMENT ON COLUMN slots.created_at IS 'Дата и время создания записи';
COMMENT ON COLUMN slots.updated_at IS 'Дата и время последнего обновления записи';
-- DOWN migration для таблицы slots
-- Удаляем все связанные объекты в обратном порядке создания

-- Удаляем триггер
DROP TRIGGER IF EXISTS update_slots_updated_at ON slots;

-- Удаляем индексы (они удалятся автоматически с таблицей, но для явности)
DROP INDEX IF EXISTS idx_slots_interviewer_status;
DROP INDEX IF EXISTS idx_slots_available_time;
DROP INDEX IF EXISTS idx_slots_created_at;
DROP INDEX IF EXISTS idx_slots_status;
DROP INDEX IF EXISTS idx_slots_end_time;
DROP INDEX IF EXISTS idx_slots_start_time;
DROP INDEX IF EXISTS idx_slots_interviewee_id;
DROP INDEX IF EXISTS idx_slots_interviewer_id;

-- Удаляем таблицу slots
DROP TABLE IF EXISTS slots;
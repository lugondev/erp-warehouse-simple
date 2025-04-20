-- Drop indexes
DROP INDEX IF EXISTS idx_report_schedules_next_run_at;
DROP INDEX IF EXISTS idx_report_schedules_created_by;
DROP INDEX IF EXISTS idx_report_schedules_active;
DROP INDEX IF EXISTS idx_report_schedules_frequency;
DROP INDEX IF EXISTS idx_report_schedules_report_type;
DROP INDEX IF EXISTS idx_reports_created_at;
DROP INDEX IF EXISTS idx_reports_status;
DROP INDEX IF EXISTS idx_reports_created_by;
DROP INDEX IF EXISTS idx_reports_type;
-- Drop tables
DROP TABLE IF EXISTS report_schedules;
DROP TABLE IF EXISTS reports;
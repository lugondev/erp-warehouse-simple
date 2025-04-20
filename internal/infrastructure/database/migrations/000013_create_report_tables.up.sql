-- Create reports table
CREATE TABLE IF NOT EXISTS reports (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	name VARCHAR(255) NOT NULL,
	description TEXT,
	type VARCHAR(50) NOT NULL,
	parameters JSONB,
	start_date TIMESTAMP NOT NULL,
	end_date TIMESTAMP NOT NULL,
	created_by INTEGER NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
	file_url TEXT,
	format VARCHAR(20) NOT NULL,
	status VARCHAR(20) NOT NULL DEFAULT 'PENDING'
);
-- Create report_schedules table
CREATE TABLE IF NOT EXISTS report_schedules (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	name VARCHAR(255) NOT NULL,
	description TEXT,
	report_type VARCHAR(50) NOT NULL,
	parameters JSONB,
	frequency VARCHAR(20) NOT NULL,
	format VARCHAR(20) NOT NULL,
	active BOOLEAN NOT NULL DEFAULT TRUE,
	recipients TEXT [],
	created_by INTEGER NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
	last_run_at TIMESTAMP,
	next_run_at TIMESTAMP
);
-- Create indexes
CREATE INDEX IF NOT EXISTS idx_reports_type ON reports(type);
CREATE INDEX IF NOT EXISTS idx_reports_created_by ON reports(created_by);
CREATE INDEX IF NOT EXISTS idx_reports_status ON reports(status);
CREATE INDEX IF NOT EXISTS idx_reports_created_at ON reports(created_at);
CREATE INDEX IF NOT EXISTS idx_report_schedules_report_type ON report_schedules(report_type);
CREATE INDEX IF NOT EXISTS idx_report_schedules_frequency ON report_schedules(frequency);
CREATE INDEX IF NOT EXISTS idx_report_schedules_active ON report_schedules(active);
CREATE INDEX IF NOT EXISTS idx_report_schedules_created_by ON report_schedules(created_by);
CREATE INDEX IF NOT EXISTS idx_report_schedules_next_run_at ON report_schedules(next_run_at);
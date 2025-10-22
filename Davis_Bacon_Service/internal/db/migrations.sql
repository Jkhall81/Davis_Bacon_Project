-- ============================================================
-- Davis-Bacon Database Schema
-- ============================================================

-- 1️ Main table: wage_determinations
CREATE TABLE IF NOT EXISTS wage_determinations (
    id TEXT PRIMARY KEY,
    wd_number TEXT NOT NULL,
    state TEXT NOT NULL,
    revision_number INT NOT NULL,
    published_date TIMESTAMP WITH TIME ZONE,
    modified_date TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_wd_state ON wage_determinations (state);
CREATE INDEX IF NOT EXISTS idx_wd_modified ON wage_determinations (modified_date DESC);

-- 2️ Counties covered by each wage determination
CREATE TABLE IF NOT EXISTS wd_counties (
    id SERIAL PRIMARY KEY,
    wd_id TEXT NOT NULL REFERENCES wage_determinations(id) ON DELETE CASCADE,
    county TEXT NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_wd_county ON wd_counties (wd_id, county);

-- 3️ Construction types per wage determination
CREATE TABLE IF NOT EXISTS wd_construction_types (
    id SERIAL PRIMARY KEY,
    wd_id TEXT NOT NULL REFERENCES wage_determinations(id) ON DELETE CASCADE,
    construction_type TEXT NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_wd_construction_type ON wd_construction_types (wd_id, construction_type);

-- 4️ Cached JSON/text documents fetched from SAM.gov
CREATE TABLE IF NOT EXISTS wd_documents (
    id SERIAL PRIMARY KEY,
    wd_id TEXT NOT NULL REFERENCES wage_determinations(id) ON DELETE CASCADE,
    revision_number INT NOT NULL,
    publish_date TIMESTAMP WITH TIME ZONE,
    document TEXT NOT NULL,
    UNIQUE (wd_id, revision_number)
);

CREATE INDEX IF NOT EXISTS idx_wd_documents_wd_id ON wd_documents (wd_id);

-- 5️ Parsed detail rows extracted from document text
CREATE TABLE IF NOT EXISTS wd_details (
    id SERIAL PRIMARY KEY,
    wd_id TEXT NOT NULL REFERENCES wage_determinations(id) ON DELETE CASCADE,
    classification TEXT NOT NULL,
    group_number INT,
    base_rate NUMERIC(10,2) DEFAULT 0,
    fringe_rate NUMERIC(10,2) DEFAULT 0,
    effective_date TIMESTAMP WITH TIME ZONE,
    notes TEXT
);

CREATE INDEX IF NOT EXISTS idx_wd_details_wd_id ON wd_details (wd_id);
CREATE INDEX IF NOT EXISTS idx_wd_details_class ON wd_details (classification);

-- 6️ Audit timestamps
ALTER TABLE wage_determinations
    ADD COLUMN IF NOT EXISTS created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW();

ALTER TABLE wd_documents
    ADD COLUMN IF NOT EXISTS created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW();

ALTER TABLE wd_details
    ADD COLUMN IF NOT EXISTS created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW();

-- 7️ Trigger to auto-update timestamps
CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_wd_updated
BEFORE UPDATE ON wage_determinations
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER trg_doc_updated
BEFORE UPDATE ON wd_documents
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER trg_details_updated
BEFORE UPDATE ON wd_details
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

-- ============================================================
-- Done
-- ============================================================

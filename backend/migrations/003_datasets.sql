-- Add dataset fields to transactions
ALTER TABLE transactions
ADD COLUMN IF NOT EXISTS source_dataset text,
ADD COLUMN IF NOT EXISTS synthetic boolean DEFAULT false,
ADD COLUMN IF NOT EXISTS fraud_label boolean;

-- Add a datasets metadata table (optional, but good for tracking scan history)
CREATE TABLE IF NOT EXISTS datasets (
    id text PRIMARY KEY,
    name text NOT NULL,
    source text,
    description text,
    row_count bigint,
    column_count int,
    quality_score numeric,
    has_fraud_labels boolean,
    status text DEFAULT 'AVAILABLE',
    last_scanned_at timestamptz DEFAULT NOW()
);

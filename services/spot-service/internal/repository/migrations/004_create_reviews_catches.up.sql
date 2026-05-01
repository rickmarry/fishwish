CREATE TABLE IF NOT EXISTS reviews (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    spot_id UUID REFERENCES spots(id) ON DELETE CASCADE,
    user_id UUID NOT NULL,
    rating INTEGER NOT NULL CHECK (rating >= 1 AND rating <= 5),
    content TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_reviews_spot ON reviews(spot_id);
CREATE INDEX idx_reviews_user ON reviews(user_id);

CREATE TABLE IF NOT EXISTS catch_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    spot_id UUID REFERENCES spots(id) ON DELETE SET NULL,
    species VARCHAR(100),
    weight_lbs DECIMAL(6,2),
    length_in DECIMAL(5,2),
    bait_used VARCHAR(100),
    weather_conditions JSONB,
    photos TEXT[],
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_catch_logs_user ON catch_logs(user_id);
CREATE INDEX idx_catch_logs_spot ON catch_logs(spot_id);
CREATE INDEX idx_catch_logs_species ON catch_logs(species);

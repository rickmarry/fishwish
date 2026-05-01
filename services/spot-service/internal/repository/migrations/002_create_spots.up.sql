CREATE TABLE IF NOT EXISTS spots (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    location GEOMETRY(Point, 4326) NOT NULL,
    type VARCHAR(50) NOT NULL CHECK (type IN ('lake', 'river', 'ocean', 'pond', 'stream', 'reservoir', 'estuary')),
    description TEXT,
    access_notes TEXT,
    parking VARCHAR(100),
    difficulty VARCHAR(20) DEFAULT 'easy' CHECK (difficulty IN ('easy', 'moderate', 'hard')),
    depth_info TEXT,
    bottom_type VARCHAR(100),
    regulations TEXT,
    rating DECIMAL(3,2) DEFAULT 0,
    review_count INTEGER DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_spots_location ON spots USING GIST(location);
CREATE INDEX idx_spots_type ON spots(type);
CREATE INDEX idx_spots_difficulty ON spots(difficulty);

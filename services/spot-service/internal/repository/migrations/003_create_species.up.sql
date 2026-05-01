CREATE TABLE IF NOT EXISTS species (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) UNIQUE NOT NULL,
    common_name VARCHAR(100) NOT NULL,
    description TEXT,
    avg_weight_lbs DECIMAL(5,2),
    record_weight_lbs DECIMAL(6,2),
    best_seasons TEXT[],
    preferred_bait TEXT[]
);

CREATE TABLE IF NOT EXISTS spot_species (
    spot_id UUID REFERENCES spots(id) ON DELETE CASCADE,
    species_id UUID REFERENCES species(id) ON DELETE CASCADE,
    abundance VARCHAR(20) DEFAULT 'common' CHECK (abundance IN ('rare', 'uncommon', 'common', 'abundant')),
    PRIMARY KEY (spot_id, species_id)
);

CREATE INDEX idx_species_name ON species(name);
CREATE INDEX idx_spot_species_species ON spot_species(species_id);

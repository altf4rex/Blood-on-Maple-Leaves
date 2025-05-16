CREATE TABLE saves (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    player_id UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    scene_id TEXT NOT NULL,
    rage INT NOT NULL DEFAULT 0,
    honor INT NOT NULL DEFAULT 0,
    karma INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
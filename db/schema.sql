-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Enums
CREATE TYPE position_enum AS ENUM (
    'QB', 'RB', 'WR', 'TE', 'PK', -- Football
    'PG', 'SG', 'SF', 'PF', 'C',          -- Basketball
    'SP', 'RP', 'CATCHER', '1B', '2B', '3B', 'SS', 'LF', 'CF', 'RF', 'DH' -- Baseball
);
CREATE TYPE player_status_enum AS ENUM ('ACTIVE', 'INJURED', 'PUP', 'SUSPENDED', 'RETIRED');
CREATE TYPE sport_type_enum AS ENUM ('FOOTBALL', 'BASKETBALL', 'BASEBALL');
CREATE TYPE draft_room_status_enum AS ENUM ('WAITING', 'DRAFTING', 'PAUSED', 'COMPLETE');

-- 1. Conferences
CREATE TABLE conferences (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- 2. Divisions
CREATE TABLE divisions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    conference_id UUID NOT NULL REFERENCES conferences(id),
    created_at TIMESTAMP DEFAULT NOW()
);

-- 3. Pro Teams
CREATE TABLE pro_teams (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    city TEXT NOT NULL,
    state TEXT,  -- Added to match Go type
    name TEXT NOT NULL,
    abbreviation TEXT NOT NULL, -- e.g. "MIN"
    logo_url TEXT,
    year_established INT,
    division_id UUID NOT NULL REFERENCES divisions(id),
    created_at TIMESTAMP DEFAULT NOW()
);

-- 4. Players
CREATE TABLE players (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    position position_enum NOT NULL,
    team_id UUID NOT NULL REFERENCES pro_teams(id),
    
    -- Physical / Career attributes
    height INT, -- in inches
    weight INT, -- in lbs
    age INT,
    years_of_experience INT CHECK (years_of_experience >= 0),
    draft_year INT,  -- Added to match Go type
    jersey_number INT,
    
    -- Meta
    status player_status_enum NOT NULL DEFAULT 'ACTIVE',
    skill DECIMAL(5,4), -- renamed from position_skill_factor, now 0.0000 - 1.0000
    headshot_url TEXT,
    
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 5. Yearly Stats
CREATE TABLE yearly_stats (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    player_id UUID NOT NULL REFERENCES players(id),
    year INT NOT NULL,
    sport_type sport_type_enum NOT NULL,
    
    -- The JSON payload
    stats JSONB NOT NULL,
    
    -- Fantasy Meta
    fantasy_points DECIMAL(10,2) NOT NULL DEFAULT 0,
    projected_fantasy_points DECIMAL(10,2) NOT NULL DEFAULT 0,
    is_projected BOOLEAN NOT NULL DEFAULT FALSE,
    games_played INT CHECK (games_played >= 0),
    
    -- Computed Column
    fantasy_points_per_game DECIMAL(10,2) GENERATED ALWAYS AS (
        CASE WHEN games_played > 0 THEN fantasy_points / games_played ELSE 0 END
    ) STORED,

    created_at TIMESTAMP DEFAULT NOW()
);

-- 6. Ranking Lists
CREATE TABLE ranking_lists (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title TEXT NOT NULL,
    author TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 7. Rankings (Updated with Surrogate Keys)
CREATE TABLE rankings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ranking_list_id UUID NOT NULL REFERENCES ranking_lists(id),
    player_id UUID NOT NULL REFERENCES players(id),
    rank INT NOT NULL CHECK (rank > 0),
    
    UNIQUE (ranking_list_id, player_id),
    UNIQUE (ranking_list_id, rank) -- No two players can have the same rank in a list
);

-- 8. Draft Rooms
CREATE TABLE draft_rooms (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    status draft_room_status_enum NOT NULL DEFAULT 'WAITING',
    timer_duration INT NOT NULL DEFAULT 60,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 9. Team Depth Charts (Pro) (Updated with Surrogate Keys)
CREATE TABLE team_depth_charts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    team_id UUID NOT NULL REFERENCES pro_teams(id),
    player_id UUID NOT NULL REFERENCES players(id),
    rank INT NOT NULL CHECK (rank > 0),
    position position_enum NOT NULL,
    updated_at TIMESTAMP DEFAULT NOW(),
    
    UNIQUE (team_id, position, rank)
);

-- 10. Users (Minimal Placeholder)
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT UNIQUE NOT NULL,
    username TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- 11. Fantasy Teams
CREATE TABLE fantasy_teams (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    draft_room_id UUID NOT NULL REFERENCES draft_rooms(id),
    user_id UUID REFERENCES users(id), -- Nullable if bot (or separate is_bot flag)
    name TEXT NOT NULL,
    draft_order_number INT,
    is_bot BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW()
);

-- 12. Fantasy Rosters (Updated with Surrogate Keys)
CREATE TABLE fantasy_rosters (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    fantasy_team_id UUID NOT NULL REFERENCES fantasy_teams(id),
    player_id UUID NOT NULL REFERENCES players(id),
    roster_spot TEXT NOT NULL, -- 'QB', 'BN', 'IR'
    created_at TIMESTAMP DEFAULT NOW(),
    
    UNIQUE (fantasy_team_id, player_id)
);

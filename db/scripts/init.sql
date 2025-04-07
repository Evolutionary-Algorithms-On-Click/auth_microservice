-- Existing tables
DROP TABLE IF EXISTS access;
DROP TABLE IF EXISTS run;
DROP TABLE IF EXISTS registerOtp;
DROP TABLE IF EXISTS users;

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    userName STRING UNIQUE NOT NULL,
    fullName STRING,
    email STRING UNIQUE NOT NULL,
    role STRING DEFAULT 'user',
    password STRING NOT NULL,
    accountStatus STRING DEFAULT 'active',
    createdAt TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    updatedAt TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL
);

CREATE TABLE IF NOT EXISTS registerOtp (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email STRING UNIQUE NOT NULL,
    otp STRING NOT NULL,
    createdAt TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    updatedAt TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL
);

CREATE TABLE IF NOT EXISTS run (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name STRING NOT NULL,
    description STRING,
    status STRING DEFAULT 'scheduled', -- 'scheduled', 'running', 'completed', 'failed'
    type STRING NOT NULL, -- 'ea', 'gp', 'ml', 'pso'
    command STRING NOT NULL,
    createdBy UUID REFERENCES users(id),
    createdAt TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    updatedAt TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL
);

CREATE TABLE IF NOT EXISTS access (
    runID UUID REFERENCES run(id),
    userID UUID REFERENCES users(id),
    mode STRING DEFAULT 'read', -- 'read', 'write'
    PRIMARY KEY (runID, userID),
    createdAt TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    updatedAt TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL
);

-- New team management tables
CREATE TABLE IF NOT EXISTS teams (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name STRING NOT NULL,
    description STRING,
    createdBy UUID REFERENCES users(id) ON DELETE CASCADE,
    createdAt TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    updatedAt TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL
);

CREATE TABLE IF NOT EXISTS team_members (
    teamID UUID REFERENCES teams(id) ON DELETE CASCADE,
    userID UUID REFERENCES users(id) ON DELETE CASCADE,
    createdAt TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    updatedAt TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    PRIMARY KEY (teamID, userID)
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_teams_createdBy ON teams(createdBy);
CREATE INDEX IF NOT EXISTS idx_team_members_userID ON team_members(userID);

-- Updated team_run_access table without grantedBy
CREATE TABLE IF NOT EXISTS team_run_access (
    runID UUID REFERENCES run(id) ON DELETE CASCADE,
    teamID UUID REFERENCES teams(id) ON DELETE CASCADE,
    mode STRING DEFAULT 'read', -- 'read', 'write'
    createdAt TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    updatedAt TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    PRIMARY KEY (runID, teamID)
);
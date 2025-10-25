-- Drop tables in reverse dependency order
DROP TABLE IF EXISTS password_reset_otps;
DROP TABLE IF EXISTS teamMembers;
DROP TABLE IF EXISTS access;
DROP TABLE IF EXISTS run;
DROP TABLE IF EXISTS registerOtp;
DROP TABLE IF EXISTS team;
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
    status STRING DEFAULT 'scheduled',
    -- 'scheduled', 'running', 'completed', 'failed'
    type STRING NOT NULL,
    -- 'ea', 'gp', 'ml', 'pso'
    command STRING NOT NULL,
    createdBy UUID REFERENCES users(id),
    createdAt TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    updatedAt TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL
);
CREATE TABLE IF NOT EXISTS access (
    runID UUID REFERENCES run(id),
    userID UUID REFERENCES users(id),
    mode STRING DEFAULT 'read',
    -- 'read', 'write'
    PRIMARY KEY (runID, userID),
    createdAt TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    updatedAt TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL
);

--table to maintain team metadata
CREATE TABLE IF NOT EXISTS team (
    teamID UUID PRIMARY KEY,
    createdBy STRING NOT NULL,
    role STRING,
    createdAt TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    updatedAt TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL
);

--table to maintain team members associated to that team
CREATE TABLE IF NOT EXISTS teamMembers (
    memberId UUID REFERENCES users(id),
    teamID UUID REFERENCES team(teamID),
    PRIMARY KEY (memberId, teamID),
    createdAt TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    updatedAt TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL
);

-- table to maintain password reset otps
CREATE TABLE IF NOT EXISTS password_reset_otps (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    otp_code STRING NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    is_used BOOLEAN DEFAULT FALSE
);

-- indexing for faster lookups
CREATE INDEX IF NOT EXISTS idx_password_reset_user_id ON password_reset_otps(user_id);

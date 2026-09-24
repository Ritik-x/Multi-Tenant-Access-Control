CREATE TABLE invitations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    organization_id UUID NOT NULL
     REFERENCES organizations(id)
   ON DELETE CASCADE,
    email TEXT NOT NULL,

    role_id UUID NOT NULL
        REFERENCES roles(id)
        ON DELETE RESTRICT,

    token_hash TEXT NOT NULL UNIQUE,

    expires_at TIMESTAMPTZ NOT NULL,

    accepted_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
)
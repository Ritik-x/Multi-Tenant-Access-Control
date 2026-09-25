CREATE TABLE audit_logs(
     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

       organization_id UUID NOT NULL
        REFERENCES organizations(id)
        ON DELETE CASCADE,
 user_id UUID
        REFERENCES users(id)
        ON DELETE SET NULL,

    action TEXT NOT NULL,

    resource TEXT NOT NULL,

 resource_id UUID,

    metadata JSONB,

    ip_address INET,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()



)
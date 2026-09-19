DELETE FROM roles
WHERE organization_id = (
    SELECT id
    FROM organizations
    WHERE slug = 'development'
)
AND name IN ('owner', 'admin', 'member');
DELETE FROM role_permissions
WHERE role_id IN (
    SELECT r.id
    FROM roles r
    JOIN organizations o
    ON o.id = r.organization_id
    WHERE o.slug = 'development'
      AND r.name IN ('owner', 'admin', 'member')
)
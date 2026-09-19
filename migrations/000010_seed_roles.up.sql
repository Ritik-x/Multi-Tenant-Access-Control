INSERT INTO roles (organization_id, name)
SELECT id, role_name
FROM  organizations
CROSS JOIN (
    VALUES
        ('owner'),
        ('admin'),
        ('member')
) AS role_data(role_name)
WHERE organizations.slug = 'development';
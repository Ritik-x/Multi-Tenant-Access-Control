-- owner get all the permissions 

INSERT INTO role_permissions( role_id , permission_id)
SELECT r.id , p.id 
FROM roles r 
CROSS JOIN permissions p 
JOIN organizations o
ON o.id = r.organization_id
WHERE o.slug = 'development'
  AND r.name = 'owner';

--   Admin gets all current permisssions

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
JOIN organizations o
    ON o.id = r.organization_id
WHERE o.slug = 'development'
  AND r.name = 'admin';


-- Member gets only read permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
JOIN organizations o
    ON o.id = r.organization_id
WHERE o.slug = 'development'
  AND r.name = 'member'
  AND p.name IN (
      'users.read',
      'team.read'
  );
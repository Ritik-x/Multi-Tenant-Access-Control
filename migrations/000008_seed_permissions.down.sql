DELETE FROM permissions
WHERE name IN (
    'users.read',
    'users.create',
    'users.update',
    'users.delete',
    'team.read',
    'team.invite',
    'team.remove',
    'audit.read'
);
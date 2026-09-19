INSERT INTO permissions ( name , description)
VALUES
 ('users.read', 'Read users'),
    ('users.create', 'Create users'),
    ('users.update', 'Update users'),
    ('users.delete', 'Delete users'),

    ('team.read', 'View team members'),
    ('team.invite', 'Invite users to the organization'),
    ('team.remove', 'Remove users from the organization'),

    ('audit.read', 'View audit logs');

    
package utils

const (
	SearchUserQuery   = `SELECT username, email, id  FROM users WHERE username ILIKE ?`
	SearchGroupQuery  = `SELECT name, id AS group_id FROM groups WHERE name ILIKE ?`
	GetAllGroupsQuery = `SELECT DISTINCT ON (gs.id) gs.* , public.groups_members.* , CASE WHEN gs.created_by = ? THEN true ELSE false END AS is_delete_able
FROM public.groups_members
LEFT JOIN public.groups AS gs
ON gs.id = public.groups_members.group_id
WHERE public.groups_members.user_id = ?`
)

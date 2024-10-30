package utils

const (
	SearchUserQuery   = `SELECT username, email, id FROM users WHERE username ILIKE ?`
	SearchGroupQuery  = `SELECT * FROM groups WHERE name ILIKE ?`
	GetAllGroupsQuery = `SELECT DISTINCT ON (gs.id) gs.* , public.groups_members.*
FROM public.groups_members
LEFT JOIN public.groups AS gs
ON gs.id = public.groups_members.group_id
WHERE public.groups_members.user_id = ?`
)

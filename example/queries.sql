--------------------------------------------------------------------------------------------------------------------
-- Users
--------------------------------------------------------------------------------------------------------------------

-- name: GetUsers :many
select *
from users
order by created_at
LIMIT $1 OFFSET $2;

-- name: GetUsersById :many
select *
from users
where id = ANY ($1::uuid[])
order by display_name;

-- name: GetUser :one
select *
from users
where id = $1;

-- name: FindByUsername :one
select *
from users
where username = $1;

-- name: CreateUser :one
insert into users (username, name, display_name, emails, active, locale, external_id, created_at, updated_at)
values ($1, $2, $3, $4, $5, $6, $7, now(), now())
returning *;

-- name: UpdateUser :exec
update users
set username     =$2,
    name         = $3,
    display_name = $4,
    emails       = $5,
    active       = $6,
    external_id  = $7,
    locale       = $8,
    updated_at   = now()
where id = $1;

-- name: PatchUser :exec
update users
set active     = $2,
    updated_at = now()
where id = $1;

-- name: DeleteUser :exec
delete
from users
where id = $1;

-- name: GetUserCount :one
select count(*)
from users;

--------------------------------------------------------------------------------------------------------------------
-- Groups
--------------------------------------------------------------------------------------------------------------------

-- name: GetGroups :many
select *
from groups
order by id
LIMIT $1 OFFSET $2;

-- name: GetGroup :one
select *
from groups
where id = $1;

-- name: CreateGroup :one
insert into groups (display_name, created_at, updated_at)
values ($1, now(), now())
returning *;

-- name: DeleteGroup :exec
delete
from groups
where id = $1;

-- name: GetGroupCount :one
select count(*)
from groups;

-- name: PatchGroupDisplayName :exec
update groups
set display_name = $2,
    updated_at   = now()
where id = $1;

--------------------------------------------------------------------------------------------------------------------
-- Membership
--------------------------------------------------------------------------------------------------------------------

-- name: GetGroupMembership :many
select group_users.*, users.username as username
from group_users
         left join users on users.id = group_users.user_id
where group_users.group_id = $1;

-- name: GetGroupMembershipForUser :one
select group_users.*, users.username as username
from group_users
         left join users on users.id = group_users.user_id
where group_users.group_id = $1
  and group_users.user_id = $2;

-- name: DropMembershipForGroup :exec
delete
from group_users
where group_id = $1;

-- name: DropMembershipForUserAndGroup :exec
delete
from group_users
where user_id = $1
  and group_id = $2;

-- name: CreateMembershipForUserAndGroup :exec
insert into group_users (user_id, group_id)
values ($1, $2)
on conflict (user_id, group_id) do nothing;

--------------------------------------------------------------------------------------------------------------------
-- SCIM API Key
--------------------------------------------------------------------------------------------------------------------

-- name: InsertAPIKey :one
insert into api_keys (encodedhash, system, owner, description, created_at, updated_at)
values ($1, $2, $3, $4, now(), now())
returning *;

-- name: InsertScimAPIKey :one
insert into scim_api_keys (api_key_id, domain, created_at, updated_at)
values ($1, 'default', now(), now())
returning *;

-- name: DeleteAPIKey :exec
delete
from api_keys
where id = $1;

-- name: DeleteScimAPIKey :exec
delete
from scim_api_keys
where domain = 'default';

-- name: FindAPIKey :one
select *
from api_keys
where id = $1
  and system = false;

-- name: FindAPIKeysById :many
select *
from api_keys
where id = ANY ($1::uuid[]);

-- name: FindScimAPIKey :one
select api_keys.*
from api_keys
         left join scim_api_keys on scim_api_keys.api_key_id = api_keys.id
where scim_api_keys.domain = 'default'
  and api_keys.system = true;

-- name: GetAPIKeys :many
select *
from api_keys
where system = false;

create table if not exists `todo_dependencies`
(
    `id`             bigint unsigned not null,
    `source_todo_id` bigint unsigned not null,
    `dest_todo_id`   bigint unsigned not null,
    `created_at`     datetime(3)     not null,
    `updated_at`     datetime(3)     not null,
    primary key (`id`),
    unique key `unique_dependency` (`source_todo_id`, `dest_todo_id`),
    index `reverse_dependency` (`dest_todo_id`, `source_todo_id`)
)
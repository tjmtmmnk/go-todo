create table if not exists `users`
(
    `id`         bigint unsigned not null,
    `name`       varchar(100)    not null,
    `nickname`   varchar(100) default null,
    -- bcryptは72byteなので余裕を持たせている
    `password`   varchar(255)    not null,
    `created_at` datetime(3)     not null,
    `updated_at` datetime(3)     not null,
    primary key (`id`),
    unique key `unique_name` (`name`)
);

create table if not exists `todos`
(
    `id`         bigint unsigned not null,
    `user_id`    bigint unsigned not null,
    `item_name`  varchar(1000)   not null,
    `done`       bool            not null default false,
    `created_at` datetime(3)     not null,
    `updated_at` datetime(3)     not null,
    primary key (`id`),
    index `user_id` (`user_id`),
    foreign key `fk_user_id` (`user_id`)
        references `users` (`id`)
        on delete cascade
);

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
);
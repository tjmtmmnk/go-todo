alter table `todos`
    add column `start_at` datetime(3) default null
        after `done`,
    add column `end_at`   datetime(3) default null
        after `start_at`
;
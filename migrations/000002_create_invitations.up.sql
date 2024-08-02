create table invitations
(
    event_id    int UNSIGNED not null,
    user_id     varchar(50)            not null,
    accepted_at datetime null,
    created_at  datetime default NOW() not null,
    updated_at  datetime default NOW() not null invisible,
    constraint invitations_pk
        primary key (event_id, user_id),
    CONSTRAINT `fk_invitations_events`
        FOREIGN KEY (event_id) REFERENCES events (id)
            ON DELETE RESTRICT
            ON UPDATE RESTRICT
) ENGINE = InnoDB
    DEFAULT CHARACTER SET = utf8
    COLLATE = utf8_unicode_ci;
CREATE TABLE IF NOT EXISTS `locations`
(
    `id`         INT UNSIGNED  NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `name`       VARCHAR(255)  NULL,
    `latitude`   DECIMAL(10, 6) NULL,
    `longitude`  DECIMAL(10, 6) NULL,
    `created_at` DATETIME      NOT NULL DEFAULT NOW(),
    `updated_at` DATETIME      NOT NULL DEFAULT NOW()
)
    ENGINE = InnoDB
    DEFAULT CHARACTER SET = utf8
    COLLATE = utf8_unicode_ci;

CREATE TABLE IF NOT EXISTS `events`
(
    `id`                      INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `external_id`             VARCHAR(50)  NOT NULL,
    `user`                    VARCHAR(50)  NOT NULL,
    `name`                    VARCHAR(255) NOT NULL,
    `location_id`             INT UNSIGNED NULL,
    `description`             TEXT         NULL,
    `start_date`              DATETIME     NULL,
    `end_date`                DATETIME     NULL,
    `capacity`                INT UNSIGNED NULL,
    `registration_start_date` DATETIME     NULL,
    `registration_end_date`   DATETIME     NULL,
    `public`                  BOOL         NOT NULL default false,
    `created_at`              DATETIME     NOT NULL DEFAULT NOW(),
    `updated_at`              DATETIME     NOT NULL DEFAULT NOW(),
    CONSTRAINT `fk_events_locations`
        FOREIGN KEY (location_id) REFERENCES locations (id)
            ON DELETE RESTRICT
            ON UPDATE RESTRICT
)
    ENGINE = InnoDB
    DEFAULT CHARACTER SET = utf8
    COLLATE = utf8_unicode_ci;



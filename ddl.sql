CREATE DATABASE `howareyou` /*!40100 DEFAULT CHARACTER SET utf8mb4 */;

CREATE TABLE `howareyou`.`user_group`
(
    `user_group_id`       BIGINT(20)   NOT NULL AUTO_INCREMENT,
    `user_group_slack_id` VARCHAR(100) NOT NULL,
    `channel_slack_id`    VARCHAR(100) NOT NULL,
    `created_at`          TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`          TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`user_group_id`),
    UNIQUE INDEX `idx_howareyou_u1` (`user_group_slack_id` ASC)
);

USE falcon_portal;


DROP TABLE IF EXISTS alarm_verification;
CREATE TABLE alarm_verification
(
    id                       INT UNSIGNED        NOT NULL    AUTO_INCREMENT,
    strategy_id              INT UNSIGNED       NOT NULL,
    strategy_copyid          INT UNSIGNED       NOT NULL,
    createdtime              INT(11) UNSIGNED  NOT NULL,
    verification_status      INT UNSIGNED       NOT NULL DEFAULT 0,
    PRIMARY KEY (id),
    KEY idx_id (id)
)
   ENGINE =InnoDB
   DEFAULT CHARSET =utf8
   COLLATE =utf8_unicode_ci;

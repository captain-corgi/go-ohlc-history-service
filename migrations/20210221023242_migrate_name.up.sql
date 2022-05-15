CREATE TABLE IF NOT EXISTS `ohlc`
(
    `id`     INT AUTO_INCREMENT NOT NULL,
    `unix`   INTEGER            NOT NULL,
    `symbol` VARCHAR(255)       NOT NULL,
    `open`   FLOAT              NOT NULL,
    `high`   FLOAT              NOT NULL,
    `low`    FLOAT              NOT NULL,
    `close`  FLOAT              NOT NULL,
    PRIMARY KEY (`id`)
);
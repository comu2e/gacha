-- DROP TABLE IF EXISTS 'maindb';

-- maindbという名前のデータベースを作成
CREATE DATABASE maindb;
-- maindb
use maindb;

CREATE TABLE IF not EXISTS characters
(
    name VARCHAR(100) NULL,
    id  INT NOT NULL PRIMARY KEY,
    CONSTRAINT name
    UNIQUE (name)
);

CREATE TABLE IF not EXISTS users
(
    id         INT NOT NULL PRIMARY KEY,
    name       VARCHAR(64)  NOT NULL,
    FirstName  VARCHAR(100) NULL,
    LastName   VARCHAR(100) NULL,
    Email      VARCHAR(100) NOT NULL,
    Password   VARCHAR(100) NOT NULL,
    Phone      VARCHAR(11)  NULL,
    UserStatus tinyint(1)   NOT NULL,
    xToken     VARCHAR(100)  NOT NULL,
    CONSTRAINT name
        UNIQUE (name)
);



CREATE TABLE IF not EXISTS user_character (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY ,
    user_id INT NOT NULL,
    FOREIGN KEY (user_id)
    REFERENCES users(id) ,
    character_id INT NOT NULL, FOREIGN KEY (character_id) REFERENCES characters(id)
);

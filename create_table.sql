-- testdbという名前のデータベースを作成
CREATE DATABASE testdb;
-- testdbをアクティブ
use testdb;

CREATE TABLE  characters
(
    name VARCHAR(100) NULL UNIQUE ,
    id  INT NOT NULL
        PRIMARY KEY
);

CREATE TABLE  users
(
    id  INT NOT NULL
        PRIMARY KEY,
    name   VARCHAR(64) UNIQUE NOT NULL ,
    FirstName  VARCHAR(100) NULL,
    LastName   VARCHAR(100) NULL,
    Email      VARCHAR(100) NOT NULL,
    Password   VARCHAR(100) NOT NULL,
    Phone      VARCHAR(11)  NULL,
    UserStatus TINYINT(1) NOT NULL,
    xToken     VARCHAR(20) NOT NULL
);

CREATE TABLE user_character (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY ,
    user_id INT NOT NULL,
    FOREIGN KEY (user_id)
    REFERENCES users(id) ,
    character_id
    INT NOT NULL,
    FOREIGN KEY (character_id)
    REFERENCES characters(id)
);

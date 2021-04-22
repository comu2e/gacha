
CREATE TABLE  characters
(
    name VARCHAR(100) NULL UNIQUE ,
    id  INT
                      NOT NULL
        PRIMARY KEY
);

CREATE TABLE  users
(
    id  INT NOT NULL
        PRIMARY KEY,
    name   TEXT NOT NULL ,
    FirstName  VARCHAR(100) NULL,
    LastName   VARCHAR(100) NULL,
    Email      VARCHAR(100) NOT NULL,
    Password   VARCHAR(100) NOT NULL,
    Phone      VARCHAR(11)  NULL,
    UserStatus TINYINT(1)
                            NOT NULL,
    xToken     VARCHAR(20)
                            NOT NULL
);

CREATE TABLE user_character (
                                id INT  PRIMARY KEY NOT NULL,
                                user_id
                                   INT NOT NULL,
                                FOREIGN KEY (id)
                                    REFERENCES users(id) ,
                                character_id
                                   INT NOT NULL,
                                FOREIGN KEY (id)
                                    REFERENCES characters(id)
);

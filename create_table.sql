create table user_character (
    id int primary key not null,
    user_id int not null ,
    character_id int not null
);

create table characters
(
    CharacterName varchar(100) null,
    id            int          not null
        primary key
);

create table users
(
    id         int          not null primary key,
    Username   text         not null,
    FirstName  varchar(100) null,
    LastName   varchar(100) null,
    Email      varchar(100) not null,
    Password   varchar(100) not null,
    Phone      varchar(11)  null,
    UserStatus tinyint(1)   not null
);


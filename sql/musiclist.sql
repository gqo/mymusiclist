-- Created by: gqo
-- 10/28/18
-- Intiial musiclistdb
create table users
    (u_id int not null AUTO_INCREMENT,
    username varchar(30) not null,
    pass varchar(64) not null,
    primary key (username)
    );

create table album
    (title varchar(255) not null,
    artist varchar(255) not null,
    cover varchar(255),
    primary key (title, artist)
    );

create table reviewed
    (username varchar(30),
    title varchar(255),
    artist varchar(255),
    rating int check (rating > 0 and rating < 11),
    primary key (username, title, artist),
    foreign key (username) references users(username)
        on delete cascade,
    foreign key (title, artist) references album(title, artist)
        on delete cascade
    );
alter user root identified with mysql_native_password by 'sqlpassword';
flush privileges;

create table if not exists Users (
    UserID INT NOT NULL auto_increment PRIMARY KEY,
    Email VARCHAR(254) NOT NULL UNIQUE,
    PassHash BINARY(128) NOT NULL,
    UserName VARCHAR(255) NOT NULL UNIQUE,
    FirstName VARCHAR(50) NOT NULL,
    LastName VARCHAR(50) NOT NULL
);

CREATE TABLE if not exists SignIns (
    SignInID INT NOT NULL auto_increment PRIMARY KEY,
    UserID INT NOT NULL,
    SignInDate DATETIME NOT NULL,
    IPAddress VARCHAR(100) NOT NULL UNIQUE
);

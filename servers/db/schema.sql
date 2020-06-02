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

CREATE TABLE IF NOT EXISTS Meeting (
    MeetingID INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    CreatorID INT NOT NULL,
    MeetingName VARCHAR(100) NOT NULL,
    MeetingDesc VARCHAR(300),
    FOREIGN KEY (CreatorID) REFERENCES Users(UserID)
);

CREATE TABLE IF NOT EXISTS MeetingMembers (
    MeetingMembersID INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    UserID INT NOT NULL,
    MeetingID INT NOT NULL,
    FOREIGN KEY (UserID) REFERENCES Users(UserID),
    FOREIGN KEY (MeetingID) REFERENCES Meeting(MeetingID)
);

CREATE TABLE IF NOT EXISTS [Time] (
    TimeID INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    TimeRange VARCHAR(30)
);

CREATE TABLE IF NOT EXISTS [Day] (
    DayID INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    DayName VARCHAR(20)
);

CREATE TABLE IF NOT EXISTS UserTimes (
    UserTimeID INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    UserID INT NOT NULL,
    TimeID INT NOT NULL,
    DayID INT NOT NULL,
    FOREIGN KEY (UserID) REFERENCES Users(UserID),
    FOREIGN KEY (TimeID) REFERENCES [Time](TimeID),
    FOREIGN KEY (DayID) REFERENCES [Day](DayID)
);
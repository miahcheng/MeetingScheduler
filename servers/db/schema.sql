create table if not exists Users (
    UserID INT NOT NULL auto_increment PRIMARY KEY,
    Email VARCHAR(254) NOT NULL UNIQUE,
    PassHash BINARY(128) NOT NULL,
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

CREATE TABLE IF NOT EXISTS `Time` (
    TimeID INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    TimeStart VARCHAR(4)
);

CREATE TABLE IF NOT EXISTS `Day` (
    DayID INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    DayName VARCHAR(20)
);

CREATE TABLE IF NOT EXISTS UserTimes (
    UserTimeID INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    UserID INT NOT NULL,
    TimeID INT NOT NULL,
    DayID INT NOT NULL,
    FOREIGN KEY (UserID) REFERENCES Users(UserID),
    FOREIGN KEY (TimeID) REFERENCES `Time`(TimeID),
    FOREIGN KEY (DayID) REFERENCES `Day`(DayID)
);

INSERT INTO `Day`(DayName)
VALUES("Sunday"),("Monday"),("Tuesday"),("Wednesday"),("Thursday"),("Friday"),("Saturday");

INSERT INTO `Time`(TimeStart)
VALUES ("0000"), ("0030"), ("0100"), ("0130"), ("0200"),
("0230"), ("0300"), ("0330"), ("0400"), ("0430"), ("0500"),
("0530"), ("0600"), ("0630"), ("0700"), ("0730"), ("0800"),
("0830"), ("0900"), ("0930"), ("1000"), ("1030"), ("1100"),
("1130"), ("1200"), ("1230"), ("1300"), ("1330"), ("1400"), 
("1430"), ("1500"), ("1530"), ("1600"), ("1630"), ("1700"),
("1730"), ("1800"), ("1830"), ("1900"), ("1930"), ("2000"),
("2030"), ("2100"), ("2130"), ("2200"), ("2230"), ("2300"),
("2330");
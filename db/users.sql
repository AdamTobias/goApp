USE goPractice;

DROP TABLE users;
CREATE TABLE users (
  id int(3) AUTO_INCREMENT,
  username varchar(75) UNIQUE,
  password varchar(75),
  PRIMARY KEY(id)
);
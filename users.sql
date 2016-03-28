USE goPractice;

/* Create other tables and define schemas for them here! */
DROP TABLE users;
CREATE TABLE users (
  id int(3) AUTO_INCREMENT,
  username varchar(75) UNIQUE,
  password varchar(75),
  PRIMARY KEY(id)
);

-- CREATE TABLE rooms (
--   id int(3) AUTO_INCREMENT,
--   roomname varchar(75) UNIQUE,
--   PRIMARY KEY(id) 
-- );

-- CREATE TABLE messages (
--   id int(3) AUTO_INCREMENT,
--   userid int(3),
--   text varchar(140),
--   roomid int(3),
--   PRIMARY KEY(id),
--   FOREIGN KEY(userId) REFERENCES users(id),
--   FOREIGN KEY(roomId) REFERENCES rooms(id)
-- );
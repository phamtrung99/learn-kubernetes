CREATE DATABASE demo;
CREATE TABLE `users` (
  `id` INT PRIMARY KEY AUTO_INCREMENT,
  `name` VARCHAR(255)
);

INSERT INTO `users` VALUES (1, "nobita"), (2, "shizuka");

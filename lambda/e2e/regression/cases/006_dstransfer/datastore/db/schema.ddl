DROP TABLE IF EXISTS expenditure;

CREATE TABLE expenditure
(
  id           INT AUTO_INCREMENT PRIMARY KEY,
  country      VARCHAR(255),
  year         int,
  category     VARCHAR(255),
  sub_category VARCHAR(255),
  expenditure  DECIMAL(7,2)
);


CREATE USER 'demonforgodojogo'@'localhost' 
IDENTIFIED BY '#replace password here before run it#';

GRANT SELECT, INSERT, UPDATE, DELETE
ON godojogo.*
TO 'demonforgodojogo'@'localhost';

FLUSH PRIVILEGES;
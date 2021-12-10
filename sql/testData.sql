USE OSUClassData;

REPLACE INTO Classes(ClassIdentifier, TermID, Students, Credits)
VALUES
    ROW('CS160', "202001", 10, 4),
    ROW('CS160', "202101", 11, 4),
    ROW('CS160', "202201", 12, 4),
    ROW('CS161', "202001", 20, 4),
    ROW('CS161', "202101", 21, 4),
    ROW('CS161', "202201", 22, 4);

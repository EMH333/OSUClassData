USE OSUClassData;

REPLACE INTO Classes(ClassIdentifier, TermID, Students, Credits, W, ClassGPA)
VALUES
    ROW('CS160', "202001", 10, 4, 1, 1),
    ROW('CS160', "202101", 11, 4, 2, 2),
    ROW('CS160', "202201", 12, 4, 3, 3),
    ROW('CS160', "202000", 3, 4, 1, 3),
    ROW('CS160', "202100", 3, 4, 2, 3),
    ROW('CS160', "202200", 4, 4, 3, 3),
    ROW('CS161', "202001", 20, 4, 4, 4),
    ROW('CS161', "202101", 21, 4, 5, 5),
    ROW('CS161', "202201", 22, 4, 6, 6);

REPLACE INTO ClassInfo(ClassIdentifier, Credits, ClassName)
VALUES
    ROW("CS160", 4, "Basics of Computer Science"),
    ROW("CS161", 4, "Intro to Computer Science I");

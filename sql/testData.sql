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
    ROW('CS161', "202201", 22, 4, 6, 3.90980980880608909890),
    ROW('CS999', "202001", 10, 4, 1, 1),
    ROW('HC407', "202001", 10, 4, 1, 1),
    ROW('CS162', "202001", 30, 4, 7, 7);

-- so we can test the bar chart
REPLACE INTO Classes (classIdentifier, termID, students, credits, a, aMinus, b, bPlus, bMinus, c, cPlus, cMinus, d, dPlus, dMinus, f, s, u, w, i, classGPA) VALUES
    ROW('CS160', '202201', 150, 3, 51, 14, 32, 15, 4, 12, 8, 4, 1, 0, 0, 4, 1, 2, 2, 0, 3.19);

REPLACE INTO ClassInfo(ClassIdentifier, Credits, ClassName, ClassDescription)
VALUES
    ROW('CS999', 4, 'Does not exist `anymore`', 'This class is not offered anymore'),
    ROW('HC407', 4, 'Does not have a name', 'This class does not have a name'),
    ROW("CS160", 4, "WRONG - Basics of Computer Science", 'Computer Science intro is a class designed to do a whole bunch of stuff blah blah blah this goes on for a while and needs to be formmatted correctly so there is that too'),
    ROW("CS161", 4, "WRONG - Intro to Computer Science I", 'Computer Science Intro (but this time it is actually programming) is a class designed to do a whole bunch of stuff blah blah blah this goes on for a while and needs to be formmatted correctly so there is that too'),
    ROW("CS162", 4, NULL, NULL);

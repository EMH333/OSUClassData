CREATE DATABASE IF NOT EXISTS OSUClassData;

USE OSUClassData;

GO
;

-- create Classes table
-- TODO assert students adds up to all the other values
CREATE TABLE IF NOT EXISTS Classes (
    ClassIdentifier VARCHAR(10) NOT NULL,
    TermID VARCHAR(10) NOT NULL,
    Students INTEGER NOT NULL,
    Credits INTEGER NOT NULL,
    ClassGPA FLOAT,
    A INTEGER,
    AMinus INTEGER,
    B INTEGER,
    BMinus INTEGER,
    BPlus INTEGER,
    C INTEGER,
    CMinus INTEGER,
    CPlus INTEGER,
    D INTEGER,
    DMinus INTEGER,
    DPlus INTEGER,
    F INTEGER,
    S INTEGER,
    U INTEGER,
    W INTEGER,
    I INTEGER,
    Visible BOOLEAN NOT NULL DEFAULT 1,
    PRIMARY KEY (ClassIdentifier, TermID)
);

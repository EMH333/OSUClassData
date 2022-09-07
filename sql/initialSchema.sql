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

CREATE TABLE IF NOT EXISTS ClassInfo (
    ClassIdentifier VARCHAR(10) PRIMARY KEY,
    Credits INTEGER NOT NULL,
    ClassName VARCHAR(255),
    -- whether or not the class name has been retrieved from the OSU API
    RetrievedClassName BOOLEAN NOT NULL DEFAULT 0,
    -- whether or not the class name has been normalized (can set whole column to 0 if there is an update to the algorithm)
    NormalizedClassName BOOLEAN NOT NULL DEFAULT 0,
    -- class description
    ClassDescription VARCHAR(4096)
);

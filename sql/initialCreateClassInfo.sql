--Create/update class info table, keeping class name the same if it has already been filled
REPLACE INTO ClassInfo(ClassIdentifier, Credits, ClassName)
SELECT
    InfoFromClasses.ClassIdentifier,
    Credits,
    ClassName
FROM
    (
        SELECT
            DISTINCT ClassIdentifier,
            Credits,
            TermID
        FROM
            Classes
        ORDER BY
            TermID DESC
    ) as InfoFromClasses LEFT JOIN
    (
        SELECT
            ClassIdentifier,
            ClassName
        FROM
            ClassInfo
    ) as InfoFromClassInfo On InfoFromClasses.ClassIdentifier = InfoFromClassInfo.ClassIdentifier;

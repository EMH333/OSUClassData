import wretch from "wretch";
import { dedupe, retry, throttlingCache } from "wretch/middlewares/index";
import queryString from "wretch/addons/queryString";

export const wretchInstance = wretch()
    .addon(queryString)
    .url("/api/v0/")
    .middlewares([dedupe(), retry(), throttlingCache()]);

export const chartOptions = {
    displayModeBar: false,
    responsive: true,
    scrollZoom: false,//!isTouchEnabled(),
};

export const datasetOptions = {
    spanGaps: true,
    normalized: true,
    cubicInterpolationMode: 'monotone',
    tension: 0.1,
};

export const chartColor = "rgb(215, 63, 9)";

function isTouchEnabled() {
    return ('ontouchstart' in window) ||
        (navigator.maxTouchPoints > 0) ||
        window.matchMedia("(pointer: coarse)").matches;
}

export function termIDtoString(termID: number): string {
    if (termID == undefined) return "";
    if (termID.toString().length != 6) return "";

    let year = parseInt(termID.toString().substring(0, 4));
    let part = parseInt(termID.toString().substring(4, 6));
    switch (part) {
        case 0:
            return "Summer " + (year - 1);
        case 1:
            return "Fall " + (year - 1);
        case 2:
            return "Winter " + year;
        case 3:
            return "Spring " + year;
    }
}

export function termIDtoPlotID(termID: number): number {
    if (termID == undefined) return 0;
    if (termID.toString().length != 6) return 0;

    let year = parseInt(termID.toString().substring(0, 4));
    let part = parseInt(termID.toString().substring(4, 6));
    // scales the termid based on a 100 year rotation, assumming 4 terms per year
    return (year % 100) * 4 + part;
}

export function convertRawDataToPlotData(terms: Array<number>, data: Array<any>): Array<{ x: string, y: string }> {
    const rawMap = new Map();
    const termConversionMap = new Map();
    for (let i = 0; i < terms.length; i++) {
        rawMap.set(termIDtoPlotID(terms[i]), data[i]);
        termConversionMap.set(termIDtoPlotID(terms[i]), termIDtoString(terms[i]));
    }

    const mediumData = Array();
    for (let index = termIDtoPlotID(terms[0]); index <= termIDtoPlotID(terms[terms.length - 1]); index++) {
        if (rawMap.has(index)) {
            mediumData.push({
                x: termConversionMap.get(index),
                y: rawMap.get(index),
            });
        } else {
            const numFake = termIDtoPlotID(terms[terms.length - 1]) - termIDtoPlotID(terms[0]) + 1; // the number of fake data points to add
            let fake = " ";
            for (let i = 0; i < index % numFake; i++) {
                fake += "​";
            }
            mediumData.push({ x: fake, y: null });
        }
    }
    return mediumData;
}

export function GPAToLetterGrade(averageGPA: number): string {
    if (averageGPA == undefined) return "";
    //converge gpa to letter grade
    if (averageGPA == 4) return "A";
    if (averageGPA >= 3.7) return "A-";
    if (averageGPA >= 3.3) return "B+";
    if (averageGPA >= 3.0) return "B";
    if (averageGPA >= 2.7) return "B-";
    if (averageGPA >= 2.3) return "C+";
    if (averageGPA >= 2.0) return "C";
    if (averageGPA >= 1.7) return "C-";
    if (averageGPA >= 1.3) return "D+";
    if (averageGPA >= 1.0) return "D";
    if (averageGPA >= 0.7) return "D-";
    return "F";
}

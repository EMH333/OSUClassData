import wretch from "wretch";
import { dedupe, retry, throttlingCache } from "wretch-middlewares";

export const wretchInstance = wretch()
    .url("api/v0/")
    .middlewares([dedupe(), retry(), throttlingCache()]);

export const chartOptions = {
    displayModeBar: false,
    responsive: true,
    scrollZoom: !isTouchEnabled(),
};

function isTouchEnabled() {
    return ( 'ontouchstart' in window ) ||
           ( navigator.maxTouchPoints > 0 ) ||
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

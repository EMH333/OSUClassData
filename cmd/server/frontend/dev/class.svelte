<script lang="ts">
  import BasicClassInfo from "./components/BasicClassInfo.svelte";
  import { onMount } from "svelte";
  import {
    wretchInstance,
    termIDtoString,
    chartColor,
    convertRawDataToPlotData,
  } from "./util";
  import {
    Chart,
    LineElement,
    LineController,
    BarElement,
    BarController,
    Legend,
    Title,
    Tooltip,
    CategoryScale,
    LinearScale,
    PointElement,
ChartData,
ScatterDataPoint,
  } from "chart.js";

  Chart.register(
    LineElement,
    LineController,
    BarElement,
    BarController,
    Legend,
    Title,
    Tooltip,
    CategoryScale,
    LinearScale,
    PointElement
  ); //make sure we register all the plugins we need

  let selectedClass: string;

  // TODO handle non-existent class
  //get class from query string on mount
  onMount(() => {
    const query = new URLSearchParams(window.location.search);
    selectedClass = query.get("class");

    if (selectedClass == null) {
      return; //don't bother rendering anything if we don't specify a class
    }

    // set page title to selectedClass
    document.title = selectedClass + " - OSU Class Data";

    createStudentsPerTermChart();
    createAvgGPAPerTermChart();
    createWithdrawalRatePerTermChart();
    createLastTermGradeDistributionChart();
  });

  function createStudentsPerTermChart() {
    wretchInstance
      .url("chart/studentsPerTerm")
      .query({ class: selectedClass })
      .get()
      .json((data) => {
        const students = data.SpecificData;
        const terms = data.Terms.map((term: string) => Number(term));

        const chartData = {
          datasets: [
            {
              label: "Students",
              data: convertRawDataToPlotData(terms, students),
              backgroundColor: chartColor,
              borderColor: chartColor,
              spanGaps: true,
              normalized: true,
            },
          ],
        };
        const chart = new Chart(
          document.getElementById("studentsPerTermChart") as HTMLCanvasElement,
          {
            type: "line",
            data: chartData,
            options: {
              maintainAspectRatio: false,
              plugins: {
                title: {
                  display: true,
                  text: "Students Per Term",
                  font: {
                    size: 20,
                  },
                },
                legend: {
                  display: false,
                },
              },
            },
          }
        );
      })
      .catch((err) => {
        console.error(err);
      });
  }
  function createAvgGPAPerTermChart() {
    wretchInstance
      .url("chart/avgGPAPerTerm")
      .query({ class: selectedClass })
      .get()
      .json((data) => {
        const avgGPA = data.SpecificData;
        const terms = data.Terms.map((term: string) => Number(term));

        const chartData = {
          datasets: [
            {
              label: "GPA",
              data: convertRawDataToPlotData(terms, avgGPA),
              backgroundColor: chartColor,
              borderColor: chartColor,
              spanGaps: true,
              normalized: true,
            },
          ],
        };
        const chart = new Chart(
          document.getElementById("avgGPAPerTermChart") as HTMLCanvasElement,
          {
            type: "line",
            data: chartData,
            options: {
              maintainAspectRatio: false,
              plugins: {
                title: {
                  display: true,
                  text: "Average GPA per Term",
                  font: {
                    size: 20,
                  },
                },
                legend: {
                  display: false,
                },
              },
            },
          }
        );
      })
      .catch((err) => {
        console.error(err);
      });
  }

  function createWithdrawalRatePerTermChart() {
    wretchInstance
      .url("chart/withdrawalRatePerTerm")
      .query({ class: selectedClass })
      .get()
      .json((data) => {
        const terms: Array<number> = data.Terms.map((term: string) =>
          Number(term)
        );
        const withdrawalRate = data.SpecificData.map((rate: number) =>
          (Number(rate) * 100).toFixed(2)
        );

        const chartData = {
          datasets: [
            {
              label: "Withdrawal Rate",
              data: convertRawDataToPlotData(terms, withdrawalRate),
              backgroundColor: chartColor,
              borderColor: chartColor,
              spanGaps: true,
              normalized: true,
            },
          ],
        };
        const chart = new Chart(
          document.getElementById(
            "withdrawalRatePerTermChart"
          ) as HTMLCanvasElement,
          {
            type: "line",
            data: chartData,
            options: {
              maintainAspectRatio: false,
              plugins: {
                title: {
                  display: true,
                  text: "Withdrawal Rate per Term",
                  font: {
                    size: 20,
                  },
                },
                legend: {
                  display: false,
                },
                tooltip: {
                  callbacks: {
                    label: function (context) {
                      if (context.parsed.y !== null) {
                        return context.parsed.y + "%";
                      }
                    },
                  },
                },
              },
              scales: {
                y: {
                  ticks: {
                    callback: function (value, index, ticks) {
                      return value + "%";
                    },
                  },
                },
              },
            },
          }
        );
      })
      .catch((err) => {
        console.error(err);
      });
  }

  function createLastTermGradeDistributionChart() {
    wretchInstance
      .url("chart/lastTermGradeDistribution")
      .query({ class: selectedClass })
      .get()
      .json((data) => {
        const chartData = {
          labels: [
            "A",
            "A-",
            "B+",
            "B",
            "B-",
            "C+",
            "C",
            "C-",
            "D+",
            "D",
            "D-",
            "F",
            "S (Pass)",
            "U (Fail)",
            "I (Incomplete)",
            "W (Withdrawal)",
          ],
          datasets: [
            {
              data: [
                data.A,
                data.AMinus,
                data.BPlus,
                data.B,
                data.BMinus,
                data.CPlus,
                data.C,
                data.CMinus,
                data.DPlus,
                data.D,
                data.DMinus,
                data.F,
                data.S,
                data.U,
                data.I,
                data.W,
              ],
              backgroundColor: chartColor,
              borderColor: chartColor,
            },
          ],
        };
        const chart = new Chart(
          document.getElementById(
            "lastTermGradeDistributionChart"
          ) as HTMLCanvasElement,
          {
            type: "bar",
            data: chartData,
            options: {
              maintainAspectRatio: false,
              plugins: {
                title: {
                  display: true,
                  text:
                    "Grade Distribution (" + termIDtoString(data.TermID) + ")",
                  font: {
                    size: 20,
                  },
                },
                legend: {
                  display: false,
                },
              },
              scales: {
                x: {
                  ticks: {
                    autoSkip: false,
                  },
                },
              },
            },
          }
        );
      })
      .catch((err) => {
        console.error(err);
      });
  }
</script>

<p style="font-weight: bold;"><a href="/" class="button-link">Go Back</a></p>
<div class="center buffer"><BasicClassInfo {selectedClass} /></div>

<div class="chart-container">
  <canvas id="lastTermGradeDistributionChart" />
</div>
<div class="chart-container">
  <canvas id="studentsPerTermChart" />
</div>
<div class="chart-container">
  <canvas id="avgGPAPerTermChart" />
</div>
<div class="chart-container">
  <canvas id="withdrawalRatePerTermChart" />
</div>
<br />
<!--Possible Graphics:
- Grade Distribution Pie Chart (from all data)-->
{#if selectedClass == null}
  <h2 class="center buffer">Please go back and select a class</h2>
{/if}
<p class="center">Copyright © 2021 Ethan Hampton</p>

<style>
  @import "./css/buttonLink.css";

  .center {
    margin: auto;
    width: max-content;
    max-width: 100%;
  }
  .buffer {
    margin-bottom: 3em;
  }
  .chart-container {
    height: 45vh;
    min-height: 200px;
    width: 100%;
    margin-bottom: 3em;
  }
</style>

<script lang="ts">
  import BasicClassInfo from "./components/BasicClassInfo.svelte";
  import { onMount } from "svelte";
  import {
    wretchInstance,
    termIDtoString,
    chartColor,
    convertRawDataToPlotData,
    datasetOptions,
  } from "./util";
  import {
    Chart,
    LineElement,
    LineController,
    BarElement,
    BarController,
    Title,
    Tooltip,
    CategoryScale,
    LinearScale,
    PointElement,
  } from "chart.js";
  import Footer from "./components/Footer.svelte";

  import type { ChartConfiguration } from "chart.js";

  Chart.register(
    LineElement,
    LineController,
    BarElement,
    BarController,
    Title,
    Tooltip,
    CategoryScale,
    LinearScale,
    PointElement
  ); //make sure we register all the plugins we need

  interface CombinedResponse {
    Terms: number[];
    SpecificData: Map<string, Float64Array>;
  }

  let selectedClass: string;

  // TODO handle non-existent class
  //get class from query string (or URL) on mount
  onMount(() => {
    const query = new URLSearchParams(window.location.search);
    selectedClass = query.get("class");

    if (selectedClass == null) {
      // then try to get it from the url
      const url = new URL(window.location.href);
      selectedClass = url.pathname.split("/").pop();
      //make sure and uppercase the class
      selectedClass = selectedClass.toUpperCase();
      //confirm it is a valid class with regex
      if (!selectedClass.match(/^[A-Z]{1,4}\d{3}/)) {
        // if it's not a valid class, then just return
        console.error("Invalid class: " + selectedClass);
        return;
      }
      console.log("Got class from url");
    }

    // set page title to selectedClass
    document.title = selectedClass + " - OSU Class Data";

    getCombinedData(selectedClass)
      .then((data) => {
        //console.log(data);
        createStudentsPerTermChart(data);
        createAvgGPAPerTermChart(data);
        createWithdrawalRatePerTermChart(data);
      })
      .catch((err) => {
        console.error(err);
      });

    createLastTermGradeDistributionChart();
  });

  function getCombinedData(selected: string) {
    return wretchInstance
      .url("chart/combinedData/" + selected)
      .get()
      .json((data) => {
        data.Terms = data.Terms.map((term: string) => Number(term)); // do some formatting
        return data;
      });
  }

  function createStudentsPerTermChart(data: CombinedResponse) {
    const students = data.SpecificData["S"];
    const terms = data.Terms;

    const chartData = {
      datasets: [
        {
          label: "Students",
          data: convertRawDataToPlotData(terms, students),
          backgroundColor: chartColor,
          borderColor: chartColor,
          ...datasetOptions,
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
      } as ChartConfiguration<"line">
    );
  }
  function createAvgGPAPerTermChart(data: CombinedResponse) {
    const avgGPA = data.SpecificData["GPA"];
    const terms = data.Terms;

    const chartData = {
      datasets: [
        {
          label: "GPA",
          data: convertRawDataToPlotData(terms, avgGPA),
          backgroundColor: chartColor,
          borderColor: chartColor,
          ...datasetOptions,
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
      } as ChartConfiguration<"line">
    );
  }

  function createWithdrawalRatePerTermChart(data: CombinedResponse) {
    const terms = data.Terms;
    const withdrawalRate = data.SpecificData["WR"].map((wr) => wr * 100);

    const chartData = {
      datasets: [
        {
          label: "Withdrawal Rate",
          data: convertRawDataToPlotData(terms, withdrawalRate),
          backgroundColor: chartColor,
          borderColor: chartColor,
          ...datasetOptions,
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
      } as ChartConfiguration<"line">
    );
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
          }  as ChartConfiguration<"bar">
        );
      })
      .catch((err) => {
        console.error(err);
      });
  }
</script>

<p><a href="/" class="button-link">Go Back</a></p>
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
<Footer />

<style>
  @import "./css/globalCSS.css";

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

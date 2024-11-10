<script lang="ts">
  import { run } from 'svelte/legacy';

  import { onMount } from "svelte";
  import { wretchInstance, chartColor, convertRawDataToPlotData, datasetOptions } from "./util";
  import {
    Chart,
    LineElement,
    LineController,
    Title,
    Tooltip,
    CategoryScale,
    LinearScale,
    PointElement,
  } from "chart.js";
  import AutoComplete from "svelecte";
  import Footer from "./components/Footer.svelte";

  import type { ChartConfiguration } from "chart.js";

  Chart.register(
    LineElement,
    LineController,
    Title,
    Tooltip,
    CategoryScale,
    LinearScale,
    PointElement
  ); //make sure we register all the plugins we need

  onMount(() => {
    loadSubjects();
  });

  let selectedSubject: string = $state("University-wide");

  //for search
  let classesToPick: string[] = $state(["University-wide"]);


  function loadSubjects() {
    console.log("Loading subjects");
    wretchInstance
      .url("subjects")
      .get()
      .json((json) => {
        classesToPick.push(...(json as string[]));
        classesToPick = classesToPick; //have to trigger svelte refresh

        const query = new URLSearchParams(window.location.search);
        let requestedSubject = query.get("subject");
        if (requestedSubject) {
          console.log("Selecting subject", requestedSubject);
          selectedSubject = requestedSubject;
        }
      })
      .catch((error) => {
        console.log(error);
      });
  }

  function createAvgGPAPerTermChart() {
    wretchInstance
      .url("subject/chart/avgGPAPerTerm")
      .query({ subject: selectedSubject })
      .get()
      .json((data) => {
        const avgGPA = data.GPA;
        const terms = data.Terms.map((term: string) => Number(term));

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
        let chartStatus = Chart.getChart("avgGPAPerTermChart");
        if (chartStatus != undefined) {
          chartStatus.destroy();
        }
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
      })
      .catch((err) => {
        console.error(err);
      });
  }

  function createWithdrawalRatePerTermChart() {
    wretchInstance
      .url("subject/chart/withdrawalRatePerTerm")
      .query({ subject: selectedSubject })
      .get()
      .json((data) => {
        const withdrawalRate = data.WithdrawalRate.map((rate: number) =>
          (Number(rate) * 100).toFixed(2)
        );
        const terms = data.Terms.map((term: string) => Number(term));

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
        let chartStatus = Chart.getChart("withdrawalRatePerTermChart");
        if (chartStatus != undefined) {
          chartStatus.destroy();
        }
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
      })
      .catch((err) => {
        console.error(err);
      });
  }
  run(() => {
    // reload charts when subject is changed
    if (selectedSubject) {
      createAvgGPAPerTermChart();
      createWithdrawalRatePerTermChart();
    }
  });
</script>

<p><a href="/" class="button-link">Go Back</a></p>
<div class="selector">
  <AutoComplete
    options={classesToPick}
    bind:value={selectedSubject}
    disableHighlight={true}
  />
</div>
<div class="chart-container">
  <canvas id="avgGPAPerTermChart"></canvas>
</div>
<div class="chart-container">
  <canvas id="withdrawalRatePerTermChart"></canvas>
</div>
<br />
<Footer />

<style>
  @import "./css/globalCSS.css";
  @import "./css/classSelector.css";

  .selector {
    margin: auto;
    width: 12em;
    max-width: 100%;
    margin-bottom: 3em;
  }
  .chart-container {
    height: 45vh;
    min-height: 200px;
    width: 100%;
    margin-bottom: 3em;
  }
</style>

<script lang="ts">
  import { onMount } from "svelte";
  import {
    wretchInstance,
    termIDtoString,
    termIDtoPlotID,
    chartColor,
  } from "./util";
  import Chart from "chart.js/auto"; //TODO change this to just import what we need
  import AutoComplete from "simple-svelte-autocomplete/src/SimpleAutocomplete.svelte";

  onMount(() => {
    loadSubjects();
    selectedClassAny = "University-wide" as unknown as any[];
  });

  let selectedSubject: string = "University-wide";

  //for search
  let classesToPick: string[] = ["University-wide"];
  let selectedClassAny: any[];

  $: selectedSubject = selectedClassAny as unknown as string;

  $: {
    // reload charts when subject is changed
    if (selectedSubject) {
      createAvgGPAPerTermChart();
      createWithdrawalRatePerTermChart();
    }
  }

  function loadSubjects() {
    console.log("Loading subjects");
    wretchInstance
      .url("subjects")
      .get()
      .json((json) => {
        classesToPick.push(...(json as string[]));
        classesToPick = classesToPick; //have to trigger svelte refresh
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
          labels: terms.map((term: number) => termIDtoString(term)),
          datasets: [
            {
              label: "GPA",
              data: avgGPA,
              backgroundColor: chartColor,
              borderColor: chartColor,
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
          }
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
        const withdrawalRate = data.WithdrawalRate;
        const terms = data.Terms.map((term: string) => Number(term));

        const chartData = {
          labels: terms.map((term: number) => termIDtoString(term)),
          datasets: [
            {
              label: "Withdrawal Rate",
              data: withdrawalRate,
              backgroundColor: chartColor,
              borderColor: chartColor,
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
          }
        );
      })
      .catch((err) => {
        console.error(err);
      });
  }
</script>

<p style="font-weight: bold;"><a href="/">Go Back</a></p>
<div class="center buffer">
  <AutoComplete
    text={undefined}
    items={classesToPick}
    bind:selectedItem={selectedClassAny}
    showClear={true}
    onFocus={() => {
      selectedClassAny = undefined;
    }}
  />
</div>
<div class="chart-container">
  <canvas id="avgGPAPerTermChart" />
</div>
<div class="chart-container">
  <canvas id="withdrawalRatePerTermChart" />
</div>
<br />
<p class="center">Copyright © 2021 Ethan Hampton</p>

<style>
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

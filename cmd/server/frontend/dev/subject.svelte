<script lang="ts">
  import { onMount } from "svelte";
  import {
    wretchInstance,
    termIDtoString,
    termIDtoPlotID,
    chartOptions,
  } from "./util";
  import Plotly from "plotly.js-basic-dist-min";
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
        classesToPick = classesToPick;//have to trigger svelte refresh
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
          x: terms.map((term: number) => termIDtoPlotID(term)),
          y: avgGPA,
          mode: "lines+markers",
          name: "GPA",
        };
        const chartLayout = {
          title: "Average GPA per Term",
          xaxis: {
            tickmode: "array",
            tickvals: terms.map((term: number) => termIDtoPlotID(term)),
            ticktext: terms.map((term: number) => termIDtoString(term)),
          },
        };
        Plotly.purge("avgGPAPerTermChart"); //clear previous chart
        Plotly.newPlot(
          "avgGPAPerTermChart",
          [chartData],
          chartLayout,
          chartOptions
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
          x: terms.map((term: number) => termIDtoPlotID(term)),
          y: withdrawalRate,
          mode: "lines+markers",
          name: "Withdrawal Rate",
        };
        const chartLayout = {
          title: "Withdrawal Rate per Term",
          xaxis: {
            tickmode: "array",
            tickvals: terms.map((term: number) => termIDtoPlotID(term)),
            ticktext: terms.map((term: number) => termIDtoString(term)),
          },
          yaxis: {
            tickformat: ".0%",
          },
        };
        Plotly.purge("withdrawalRatePerTermChart"); // clear previous chart
        Plotly.newPlot(
          "withdrawalRatePerTermChart",
          [chartData],
          chartLayout,
          chartOptions
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
      console.log("focus");
      selectedClassAny = undefined;
    }}
  />
</div>
<div id="avgGPAPerTermChart" />
<div id="withdrawalRatePerTermChart" />
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
</style>

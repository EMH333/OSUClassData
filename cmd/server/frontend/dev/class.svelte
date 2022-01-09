<script lang="ts">
  import BasicClassInfo from "./components/BasicClassInfo.svelte";
  import { onMount } from "svelte";
  import {
    wretchInstance,
    termIDtoString,
    termIDtoPlotID,
    chartOptions,
  } from "./util";
  import Plotly from "plotly.js-basic-dist";

  let selectedClass: string;

  // TODO handle non-existent class and/or query string
  //get class from query string on mount
  onMount(() => {
    const query = new URLSearchParams(window.location.search);
    selectedClass = query.get("class");

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
          x: terms.map((term: number) => termIDtoPlotID(term)),
          y: students,
          mode: "lines+markers",
          name: "Students",
        };
        const chartLayout = {
          title: "Students per Term",
          xaxis: {
            tickmode: "array",
            tickvals: terms.map((term: number) => termIDtoPlotID(term)),
            ticktext: terms.map((term) => termIDtoString(term)),
          },
        };
        Plotly.newPlot(
          "studentsPerTermChart",
          [chartData],
          chartLayout,
          chartOptions
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
      .url("chart/withdrawalRatePerTerm")
      .query({ class: selectedClass })
      .get()
      .json((data) => {
        const withdrawalRate = data.SpecificData;
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

  function createLastTermGradeDistributionChart() {
    wretchInstance
      .url("chart/lastTermGradeDistribution")
      .query({ class: selectedClass })
      .get()
      .json((data) => {
        const chartData = {
          x: [
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
            "S",
            "U",
            "I",
            "W",
          ],
          y: [
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
          type: "bar",
        };
        const chartLayout = {
          title: "Grade Distribution (last term)",
        };
        Plotly.newPlot(
          "lastTermGradeDistributionChart",
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
<div class="center buffer"><BasicClassInfo {selectedClass} /></div>

<div id="lastTermGradeDistributionChart" />
<div id="studentsPerTermChart" />
<div id="avgGPAPerTermChart" />
<div id="withdrawalRatePerTermChart" />
<br />
<!--Possible Graphics:<br />
- Grade Distribution Pie Chart (from all data)<br />
- Grade Distribution Pie Chart (from last term)<br />
- Num As/Bs/etc over time<br />
- Students per term over time<br />
- Withdrawal Rate over time<br />-->
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

<script lang="ts">
  import BasicClassInfo from "./components/BasicClassInfo.svelte";
  import { onMount } from "svelte";
  import { wretchInstance, termIDtoString, termIDtoPlotID, chartOptions } from "./util";
  import Plotly from "plotly.js-basic-dist";

  let selectedClass: string;

  // TODO handle non-existent class and/or query string
  //get class from query string on mount
  onMount(() => {
    const query = new URLSearchParams(window.location.search);
    selectedClass = query.get("class");

    createStudentsPerTermChart();
    createAvgGPAPerTermChart();
  });

  function createStudentsPerTermChart() {
    wretchInstance
      .url("chart/studentsPerTerm")
      .query({ class: selectedClass })
      .get()
      .json((data) => {
        const students = data.Students;
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
        Plotly.newPlot("studentsPerTermChart", [chartData], chartLayout, chartOptions);
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
        Plotly.newPlot("avgGPAPerTermChart", [chartData], chartLayout, chartOptions);
      })
      .catch((err) => {
        console.error(err);
      });
  }
</script>
<p><a href="/">Go Back</a></p>
<BasicClassInfo {selectedClass} />
<div id="studentsPerTermChart" />
<div id="avgGPAPerTermChart" />
<br />
Possible Graphics:<br />
- Grade Distribution Pie Chart (from all data)<br />
- Grade Distribution Pie Chart (from last term)<br />
- Num As/Bs/etc over time<br />
- Students per term over time<br />
- Withdrawl Rate over time<br />
<p>Copyright © 2021 Ethan Hampton</p>
<style>
</style>

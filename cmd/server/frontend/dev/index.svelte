<script lang="ts">
  import AutoComplete from "simple-svelte-autocomplete/src/SimpleAutocomplete.svelte";
  import { onMount } from "svelte";
  import wretch from "wretch";
  import { dedupe, retry, throttlingCache } from "wretch-middlewares";

  interface BasicClass {
    displayName: string;
    id: string;
  }

  const wretchInstance = wretch()
    .url("api/v0/")
    .middlewares([dedupe(), retry(), throttlingCache()]);

  //for search
  let classesToPick: BasicClass[];
  let selectedClass: any[];

  //for display
  let classInfo: any;

  onMount(() => loadClasses());

  $: {
    if (selectedClass) {
      loadClass((selectedClass as unknown as BasicClass).id);
    }
  }

  function loadClasses() {
    console.log("Loading classes");
    wretchInstance
      .url("classes")
      .get()
      .json((json) => {
        classesToPick = (json as string[]).flatMap((className: string) => ({
          displayName: className,
          id: className,
        }));
      })
      .catch((error) => {
        console.log(error);
      });
  }

  function loadClass(id: string) {
    wretchInstance
      .url("classInfo")
      .query({ class: id })
      .get()
      .json((json) => {
        classInfo = json;
      })
      .catch((error) => {
        console.log(error);
      });
  }

  function GPAToLetterGrade(averageGPA: number): string {
    if (averageGPA == undefined) return "";
    return averageGPA.toFixed(2); //TODO: convert to letter grade
  }
</script>

<AutoComplete
  text="Search for a class"
  items={classesToPick}
  bind:selectedItem={selectedClass}
  labelFieldName="displayName"
/>

{#if selectedClass}
  <!-- <div>Class info: {JSON.stringify(classInfo)}</div> -->
  {#if classInfo && classInfo.ClassName != null}
    <div>Class ID: {classInfo.ClassIdentifier}</div>
    <div>Class Name: {classInfo.ClassName}</div>
    <div>Credits: {classInfo.Credits}</div>
    <div>Last Term With Data: {classInfo.LastTerm}</div>
    <div>
      Average Grade (from all data): {GPAToLetterGrade(classInfo.AverageGPA)}
    </div>
    <div>
      Average Grade (from last term): {GPAToLetterGrade(
        classInfo.AverageGPALastTerm
      )}
    </div>
    <div>Average Students Per Term: {classInfo.AverageStudents}</div>
    <div>Students last term: {classInfo.StudentsLastTerm}</div>
    <div>Withdrawl Rate: {(classInfo.WithdrawlRate * 100).toFixed(2)}%</div>
  {/if}

  <br />
  Possible Graphics:<br />
  - Grade Distribution Pie Chart (from all data)<br />
  - Grade Distribution Pie Chart (from last term)<br />
  - Num As/Bs/etc over time<br />
  - Students per term over time<br />
  - Withdrawl Rate over time<br />
{:else}
  <div>Please pick a class!</div>
{/if}

<style>
  /* your styles go here */
</style>

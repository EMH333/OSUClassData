<script lang="ts">
  import AutoComplete from "svelecte";
  import { onMount } from "svelte";
  import { wretchInstance } from "../util";
  import type { BasicClass } from "../types";

  //for subjects
  let subjectsToPick: BasicClass[];
  let selectedSubject: BasicClass;

  //for classes
  let classesToPick: BasicClass[];
  export let selectedClass: BasicClass;

  onMount(() => loadSubjects());

  function loadSubjects() {
    console.log("Loading subjects");
    wretchInstance
      .url("subjects")
      .get()
      .json((json) => {
        subjectsToPick = (json as string[]).flatMap((className: string) => ({
          label: className,
          id: className,
        }));
      })
      .catch((error) => {
        console.error(error);
      });
  }

  function loadClasses(subject) {
    console.log("Loading classes");
    wretchInstance
      .url(`classes/${subject}`)
      .get()
      .json((json) => {
        classesToPick = (json as string[]).flatMap((className: string) => ({
          label: className,
          id: className,
        }));
      })
      .catch((error) => {
        console.error(error);
      });
  }
</script>

<div class="selector">
  <div class="subject">
    <div class="label">Subject:</div>
    <AutoComplete
      placeholder="Select a subject"
      options={subjectsToPick}
      bind:value={selectedSubject}
      on:change={()=>{
      classesToPick = []; // clear the class list
      selectedClass = null;
      loadClasses(selectedSubject.id); // load the classes
    }}
      valueAsObject={true}
      virtualList={true}
      disableHighlight={true}
    />
  </div>
  <div class="class">
    <div class="label">Class:</div>
    <AutoComplete
      placeholder="Search for a class"
      options={classesToPick}
      bind:value={selectedClass}
      disabled={!(classesToPick?.length > 0 && selectedSubject)}
      valueAsObject={true}
      virtualList={true}
      disableHighlight={true}
    />
  </div>
</div>

<style>
  @import "../css/classSelector.css";

  .selector {
    max-width: 22.5em;
    margin: 0 auto;
  }

  .label {
      padding: 0.4em 0;
  }
</style>

<script lang="ts">
  import AutoComplete from "simple-svelte-autocomplete/src/SimpleAutocomplete.svelte";
  import { onMount } from "svelte";
  import { wretchInstance } from "./util";
  import BasicClassInfo from "./components/BasicClassInfo.svelte";

  interface BasicClass {
    displayName: string;
    id: string;
  }

  //for search
  let classesToPick: BasicClass[];
  let selectedClassAny: any[];
  let selectedClass: BasicClass;

  onMount(() => loadClasses());

  $: selectedClass = selectedClassAny as unknown as BasicClass;

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
</script>

<div class="container">
  <h1>OSU Class Data Explorer</h1>

  <AutoComplete
    text="Search for a class"
    items={classesToPick}
    bind:selectedItem={selectedClassAny}
    labelFieldName="displayName"
  />

  {#if selectedClass}
    <!-- <div>Class info: {JSON.stringify(classInfo)}</div> -->
    <BasicClassInfo selectedClass={selectedClass.id} />

    <!--Link to more info-->
    <a href={`class.html?class=${selectedClass.id}`} style="font-weight: bold;">More info about {selectedClass.id}</a>
  {:else}
    <div>Please pick a class!</div>
  {/if}
  <p><a href="subject.html">Stats by Subject</a></p>
  <p><a href="about.html">About This Website</a></p>
  <p>Copyright © 2021 Ethan Hampton</p>
</div>

<style>
  .container {
    margin: auto;
    width: max-content;
    max-width: 100%;
    text-align: center;
  }
</style>

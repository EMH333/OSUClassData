<script lang="ts">
  import AutoComplete from "simple-svelte-autocomplete/src/SimpleAutocomplete.svelte";
  import { onMount } from "svelte";
  import wretch from "wretch";
  import { dedupe, retry, throttlingCache } from "wretch-middlewares";

interface BasicClass { displayName: string; id: string }

  const wretchInstance = wretch()
    .url("api/v1/")
    .middlewares([dedupe(), retry(), throttlingCache()]);

  //for search
  let classesToPick: BasicClass[];
  let selectedClass: any[];

  //for display
  let classInfo: any;

  onMount(() => loadClasses());

  $: {
    if(selectedClass){
      loadClass((selectedClass as unknown as BasicClass).id, "202001")
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

  function loadClass(id: string, term: string) {
    wretchInstance
      .url("class")
      .query({ class: id, term: term })
      .get()
      .json((json) => {
        classInfo = json;
      })
      .catch((error) => {
        console.log(error);
      });
  }
</script>

<AutoComplete
  text="Search for a class"
  items={classesToPick}
  bind:selectedItem={selectedClass}
  labelFieldName="displayName"
/>

{#if selectedClass}
  <div>Class info: {JSON.stringify(classInfo)}</div>
{:else}
  <div>Please pick a class!</div>
{/if}

<style>
  /* your styles go here */
</style>

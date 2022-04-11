<script lang="ts">
  import AutoComplete from "simple-svelte-autocomplete";
  import { onMount } from "svelte";
  import { wretchInstance } from "../util";
  import type { BasicClass } from "../types";

  //for search
  let classesToPick: BasicClass[];
  let selectedClassAny: any[];
  export let selectedClass: BasicClass;

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
        console.error(error);
      });
  }
</script>

<AutoComplete
    text="Search for a class"
    items={classesToPick}
    bind:selectedItem={selectedClassAny}
    labelFieldName="displayName"
  />

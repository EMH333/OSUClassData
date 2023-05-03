<script lang="ts">
  import AutoComplete from "svelecte";
  import { onMount } from "svelte";
  import { wretchInstance } from "../util";
  import type { BasicClass } from "../types";

  //for search
  let classesToPick: BasicClass[];
  export let selectedClass: BasicClass;

  onMount(() => loadClasses());

  function loadClasses() {
    console.log("Loading classes");
    wretchInstance
      .url("classes")
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
  <AutoComplete
    placeholder="Search for a class"
    options={classesToPick}
    bind:value={selectedClass}
    valueAsObject={true}
    virtualList={true}
    disableHighlight={true}
  />
</div>

<style>
  @import "../css/classSelector.css";

  .selector {
    max-width: 22.5em;
    margin: 0 auto;
  }
</style>

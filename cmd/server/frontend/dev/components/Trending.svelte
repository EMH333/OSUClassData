<script lang="ts">
  import { onMount } from "svelte";
  import { wretchInstance } from "../util";

  export let college: string = null;
  let classes: Set<string> = new Set(); // this is a set because backend might send duplicates

  onMount(() => {
    //allow for a college to be passed in
    let url = "trendingClasses";
    if (college) {
      url += `?college=${college}`;
    }

    wretchInstance
      .url(url)
      .get()
      .json((json) => {
        //for each json entry, add to classes set
        for (const c of (json as string[])) {
          classes.add(c);
        }
        classes = classes;//force update
      })
      .catch((error) => {
        console.error(error);
      });
  });
</script>

<!--only display if able to load properly-->
{#if classes != null && classes.size > 0}
  <div id="trendingClasses">
    <div style="font-weight: bold;">Trending Classes: </div>
    {#each [...classes] as c}
      <span class="trendingClass"><a href="class.html?class={c}">{c}</a></span>
    {/each}
  </div>
{/if}

<style>
  .trendingClass {
    margin: 0.5rem;
  }
  #trendingClasses {
    justify-content: center;
    margin: 1rem;
  }
</style>

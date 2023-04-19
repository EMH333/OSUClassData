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
        if (json) {
          //for each json entry, add to classes set
          for (const c of json as string[]) {
            classes.add(c);
          }
          classes = classes; //force update
        }
      })
      .catch((error) => {
        console.error(error);
      });
  });
</script>

<!--only display if able to load properly-->
{#if classes != null && classes.size > 0}
  <div id="trendingClasses">
    <h3 style="margin: 0;">Trending Classes:</h3>
    {#each [...classes] as c}
      <div class="trendingClass">
        <a href="/class/{c}" class="button-link">{c}</a>
      </div>
    {/each}
  </div>
{/if}

<style>
  .trendingClass {
    margin: 0.5rem;
    display: inline-block;
  }
  #trendingClasses {
    justify-content: center;
    margin: 1rem;
  }
</style>

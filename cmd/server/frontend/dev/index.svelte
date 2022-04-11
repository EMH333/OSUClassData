<script lang="ts">
  import ClassSelector from "./components/ClassSelector.svelte";
  import BasicClassInfo from "./components/BasicClassInfo.svelte";
  import Trending from "./components/Trending.svelte";
  import type { BasicClass } from "./types";

  let selectedClass: BasicClass;

  let Beta;

  async function loadBeta(eventDetails) {
    if (eventDetails.target.open && Beta === undefined) {
      Beta = (await import("./components/BetaSwitch.svelte")).default;
    }
  }
</script>

<div class="container">
  <h1>OSU Class Data Explorer</h1>

  <ClassSelector bind:selectedClass />

  {#if selectedClass}
    <!-- <div>Class info: {JSON.stringify(classInfo)}</div> -->
    <BasicClassInfo selectedClass={selectedClass.id} />

    <!--Link to more info-->
    <a href={`class.html?class=${selectedClass.id}`} style="font-weight: bold;">
      More info about {selectedClass.id}
    </a>
  {:else}
    <div>Please pick a class!</div>
  {/if}
  <Trending />
  <p><a href="subject.html">Stats by Subject</a></p>
  <p><a href="about.html">About This Website</a></p>

  <details on:toggle={loadBeta}>
    <summary>Advanced</summary>
    <svelte:component this={Beta} />
  </details>

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

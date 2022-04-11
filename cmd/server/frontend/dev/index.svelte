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
    <p>
      <a href={`class.html?class=${selectedClass.id}`} style="font-weight: bold;" class="button-link">
        More info about {selectedClass.id}
      </a>
    </p>
  {:else}
    <div>Please pick a class!</div>
  {/if}
  <div class="spacer" />
  <Trending />
  <div class="spacer" />
  <p><a href="subject.html" class="button-link">Stats by Subject</a></p>
  <p><a href="about.html" class="button-link">About This Website</a></p>

  <!--<details on:toggle={loadBeta}>
    <summary>Advanced</summary>
    <svelte:component this={Beta} />
  </details>-->

  <p>Copyright © 2021 Ethan Hampton</p>
</div>

<style global>
  @import "./globalCSS.css";
</style>

<script lang="ts">
  import ClassSelector from "./components/ClassSelector.svelte";
  import BasicClassInfo from "./components/BasicClassInfo.svelte";
  import Trending from "./components/Trending.svelte";
  import type { BasicClass } from "./types";
  import { wretchInstance } from "./util";

  let selectedClass: BasicClass;

  let Beta;

  async function loadBeta(eventDetails) {
    if (eventDetails.target.open && Beta === undefined) {
      Beta = (await import("./components/BetaSwitch.svelte")).default;
    }
  }

  function submitEmail(form: SubmitEvent) {
    form.preventDefault();
    let email = (
      (form.target as HTMLFormElement).elements[0] as HTMLInputElement
    ).value;
    wretchInstance
      .url("subscribe")
      .post({ email })
      .json()
      .then(() => {
        alert("Subscribed!");
      });
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
      <a
        href={`class.html?class=${selectedClass.id}`}
        style="font-weight: bold;"
        class="button-link"
      >
        More info about {selectedClass.id}
      </a>
    </p>
  {:else}
    <div>Please pick a class!</div>
  {/if}
  <noscript>
    <div>
      <p>
        <strong> You must have JavaScript enabled to use this site. </strong>
      </p>
      <p>
        This web app allows the exploration of data about courses, subjects, and
        classes at OSU (Oregon State University) which was obtained through a
        FOIA request to the university. This website is in no way affiliated
        with Oregon State University.
      </p>
    </div>
  </noscript>
  <div class="spacer" />
  <Trending />
  <div class="spacer" />
  <div>
    <p style="width: 50%; min-width: 15em; margin:auto;">
      Want to get notified when new terms or features are added? Submit your
      email and we'll let you know when new data is available:
    </p>
    <form on:submit={submitEmail}>
      <input type="email" name="email" />
      <input type="submit" value="Subscribe" />
    </form>
  </div>
  <div class="spacer" />
  <p><a href="subject.html" class="button-link">Stats by Subject</a></p>
  <p><a href="leaderboards" class="button-link">Class Leaderboards</a></p>
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

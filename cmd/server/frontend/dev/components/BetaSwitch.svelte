<script>
  import { run } from 'svelte/legacy';

  import { onMount } from "svelte";
  import Toggle from "svelte-toggle";
  import { BetaEnabledKey, isBetaEnabled } from "../utils/beta";

  let betaEnabled = $state(false);

  // make sure it reflects the correct state across reloads
  onMount(() => {
    betaEnabled = isBetaEnabled();
  });

  run(() => {
    if (betaEnabled) {
      //put in local storage
      localStorage.setItem(BetaEnabledKey, "true");
    } else {
      //remove from local storage
      localStorage.removeItem(BetaEnabledKey);
    }
  });
</script>

<div class="container">
  <Toggle
    bind:toggled={betaEnabled}
    label="Beta Mode Enabled"
    on="Yes"
    off="No"
    toggledColor="#D73F09"
  />
</div>

<style>
  .container {
    margin: auto;
    width: max-content;
    max-width: 100%;
    text-align: center;
  }
</style>

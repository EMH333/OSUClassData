<script lang="ts">
  import BasicClassInfo from "./components/BasicClassInfo.svelte";
  import { onMount } from "svelte";

  let selectedClass: string;

  // TODO handle non-existent class and/or query string
  //get class from query string on mount
  onMount(() => {
    const query = new URLSearchParams(window.location.search);
    selectedClass = query.get("class");
  });
</script>

<BasicClassInfo {selectedClass} />
<br />
  Possible Graphics:<br />
  - Grade Distribution Pie Chart (from all data)<br />
  - Grade Distribution Pie Chart (from last term)<br />
  - Num As/Bs/etc over time<br />
  - Students per term over time<br />
  - Withdrawl Rate over time<br />
<style>
</style>

<script lang="ts">
  export let classID: string;
  export let difficulty: number;

  let subject: string;
  let difficultyText: string;

  $: {
    switch (difficulty) {
      case 0:
        difficultyText = "much easier than";
        break;
      case 1:
        difficultyText = "easier than";
        break;
      case 2:
        difficultyText = "about the same as";
        break;
      case 3:
        difficultyText = "harder than";
        break;
      case 4:
        difficultyText = "much harder than";
        break;
      default:
        difficultyText = "Unknown";
        break;
    }
  }

  $: {
    subject = getSubjectFromClassID(classID);
  }

  function getSubjectFromClassID(classID: string): string {
    return classID.substring(0, classID.match(/\d+/)![0].length);
  }
</script>

<div>
  <span
    >This class is {{ difficultyText }} other classes in {{
      subject,
    }}</span
  >
  <span
    ><form action="/TODO" method="post">
      <input type="hidden" name="classToVote" value={{ classID }} />
      <!--Add some sort of CORS value-->
      <button type="submit" name="vote" value="Agree">Agree</button>
      <button type="submit" name="vote" value="Disagree">Disagree</button>
      <!--If disagree, then what direction should it go?-->
    </form></span
  >
</div>

<style>
</style>

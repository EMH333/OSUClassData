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
        difficultyText = "slightly easier than";
        break;
      case 3:
        difficultyText = "about the same as";
        break;
      case 4:
        difficultyText = "slightly harder than";
        break;
      case 5:
        difficultyText = "harder than";
        break;
      case 6:
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
  <span>This class is {{ difficultyText }} other {{ subject }} classes</span>
  <!--<span
    ><form action="/TODO" method="post">
      <input type="hidden" name="classToVote" value={{ classID }} />-->
  <!--Add some sort of CORS value-->
  <!--Make sure to include the rating this vote is agreeing or disagreeing with-->
  <!--<button type="submit" name="vote" value="Agree">Agree</button>
      <button type="submit" name="vote" value="DisagreeTooHigh">Disagree ({{ classID }} is easier)</button>
      <button type="submit" name="vote" value="DisagreeTooLow">Disagree ({{ classID }} is harder)</button>-->
  <!--If disagree, then what direction should it go?-->
  <!--</form></span>-->
</div>

<style>
</style>

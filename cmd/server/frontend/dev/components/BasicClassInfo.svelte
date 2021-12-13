<script lang="ts">
  import { wretchInstance, termIDtoString } from "../util";

  export let selectedClass: string;
  let classInfo: any;

  $: {
    if (selectedClass) {
      loadClass(selectedClass);
    }
  }

  function loadClass(id: string) {
    wretchInstance
      .url("classInfo")
      .query({ class: id })
      .get()
      .json((json) => {
        classInfo = json;
      })
      .catch((error) => {
        console.log(error);
      });
  }

  function GPAToLetterGrade(averageGPA: number): string {
    if (averageGPA == undefined) return "";
    return averageGPA.toFixed(2); //TODO: convert to letter grade
  }
</script>

{#if classInfo && classInfo.ClassName != null}
  <div>Class ID: {classInfo.ClassIdentifier}</div>
  <div>Class Name: {classInfo.ClassName}</div>
  <div>Credits: {classInfo.Credits}</div>
  <div>Last Term With Data: {termIDtoString(classInfo.LastTerm)}</div>
  <div>
    Average Grade (from all data): {GPAToLetterGrade(classInfo.AverageGPA)}
  </div>
  <div>
    Average Grade (from last term): {GPAToLetterGrade(
      classInfo.AverageGPALastTerm
    )}
  </div>
  <div>Average Students Per Term: {classInfo.AverageStudents}</div>
  <div>Students last term: {classInfo.StudentsLastTerm}</div>
  <div>Withdrawl Rate: {(classInfo.WithdrawlRate * 100).toFixed(2)}%</div>
{/if}

<style>
</style>

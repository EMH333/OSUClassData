<script lang="ts">
  import { wretchInstance } from "../util";

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

  function termIDtoString(termID: number): string {
    if (termID == undefined) return "";
    if (termID.toString().length != 6) return "";

    let year = parseInt(termID.toString().substring(0, 4));
    let part = parseInt(termID.toString().substring(4, 6));
    switch (part) {
      case 0:
        return "Summer " + (year - 1);
      case 1:
        return "Fall " + (year - 1);
      case 2:
        return "Winter " + year;
      case 3:
        return "Spring " + year;
    }
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

<script lang="ts">
  import { wretchInstance, termIDtoString, GPAToLetterGrade } from "../util";

  export let selectedClass: string;
  let classInfo: any;

  $: {
    if (selectedClass) {
      loadClass(selectedClass);
    }
  }

  function loadClass(id: string) {
    wretchInstance
      .url("classInfo/" + id)
      .get()
      .json((json) => {
        classInfo = json;
      })
      .catch((error) => {
        console.error(error);
      });

    // track what clsses are being looked at
    try {
      //get around typescript error here
      (window as Window & typeof globalThis & { umami: any }).umami.track(
        "classInfo",
        { class: id }
      );
    } catch (e) {
      // umami not loaded so don't worry about it
    }
  }
</script>

<div class="display">
  {#if classInfo && classInfo.ClassName != null}
    <h2>
      {classInfo.ClassIdentifier}
      {#if classInfo.ClassName != null && classInfo.ClassName != ""}
        - {classInfo.ClassName}
      {/if}
    </h2>

    <div>
      Credits: <span class="data">
        {classInfo.Credits > 0 ? classInfo.Credits : "Variable"}
      </span>
    </div>
    <div>
      Last Term With Data: <span class="data">
        {termIDtoString(classInfo.LastTerm)}
      </span>
    </div>
    <div>
      Average Grade (from all data): <span class="data">
        {GPAToLetterGrade(classInfo.AverageGPA)}
      </span>
    </div>
    <div>
      Average Grade (from last term): <span class="data">
        {GPAToLetterGrade(classInfo.AverageGPALastTerm)}
      </span>
    </div>
    <div>
      Average Students Per Term: <span class="data"
        >{parseInt(classInfo.AverageStudents).toFixed(2)}</span
      >
    </div>
    <div>
      Students last term: <span class="data">{classInfo.StudentsLastTerm}</span>
    </div>
    <div>
      Withdrawal Rate: <span class="data">
        {(classInfo.WithdrawalRate * 100).toFixed(2)}%
      </span>
    </div>
    <div>
      Pass Rate (C or better): <span class="data">
        {(classInfo.PassRate * 100).toFixed(2)}%
      </span>
    </div>
    {#if classInfo.ClassDescription}
      <div class="description">
        Description: <br />
        <span>
          {classInfo.ClassDescription}
        </span>
      </div>
    {/if}
  {/if}
</div>

<style>
  .data {
    font-weight: bold;
  }
  .display {
    text-align: center;
  }
  .description {
    min-width: 100%;
    width: 0;
    margin-top: 1em;
  }
</style>

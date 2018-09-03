<template>
  <div>
    <table>
      <thead><th>Name</th><th>Branch</th></thead>
      <tbody>
        <tr v-for="(project, idx) in data.projects" :key="idx" >
          <td>{{ project.name }}</td>
          <td>{{ project.currentBranch }}</td>
        </tr>
      </tbody>
      </table>
    <p>{{ JSON.stringify(data) }}</p>
  </div>
</template>

<script>
export default {
  props: {
    data: Object
  }
};
</script>

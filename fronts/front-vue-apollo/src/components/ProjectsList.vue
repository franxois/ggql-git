<template>
  <div class="apollo-example">
    <!-- Apollo watched Graphql query -->
    <ApolloQuery
      :query="require('../graphql/ProjectList.gql')"
    >
      <template slot-scope="{ result: { error, data }, isLoading }">

        <!-- isLoading -->
        <div v-if="isLoading" class="loading apollo">Loading...</div>

        <!-- Error -->
        <div v-else-if="error" class="error apollo">An error occured</div>

        <!-- Result -->
        <div v-else-if="data" class="result apollo"> <projecs-display :data="data"></projecs-display> </div>

        <!-- No result -->
        <div v-else class="no-result apollo">No result :(</div>
      </template>
    </ApolloQuery>
  </div>
</template>

<script>
import ProjecsDisplay from "./ProjectsList.display";

export default {
  components: {
    "projecs-display": ProjecsDisplay
  }
};
</script>

<style scoped>
.form,
.input,
.apollo,
.message {
  padding: 12px;
}

.input {
  font-family: inherit;
  font-size: inherit;
  border: solid 2px #ccc;
  border-radius: 3px;
}

.error {
  color: red;
}
</style>

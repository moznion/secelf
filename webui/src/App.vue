<template>
  <div id="app">
    <h1>SECELF</h1>

    <div role="tablist">
      <b-card no-body class="mb-1">
        <b-card-header header-tag="header" class="p-1" role="tab">
          <a block href="#" v-b-toggle.upload-pane>▼ Upload</a>
        </b-card-header>
        <b-collapse id="upload-pane" accordion="upload-pane" role="tabpanel">
          <b-card-body>
            <form action="/api/files" method="POST" enctype="multipart/form-data">
              <div class="input-group">
                <div class="custom-file">
                  <input type="file" class="custom-file-input" id="file" name="file">
                  <label class="custom-file-label" for="file">Choose file</label>
                </div>
                <div class="input-group-append">
                  <button class="btn btn-primary" type="submit">Upload</button>
                </div>
              </div>
            </form>
          </b-card-body>
        </b-collapse>
      </b-card>
    </div>

    <form v-on:submit.prevent="search">
      <div class="searchbox input-group mb-3">
        <input type="text" class="form-control" v-model="searchQuery" placeholder="Query">
        <div class="input-group-append">
          <button class="btn btn-outline-secondary" type="submit">Search</button>
        </div>
      </div>
    </form>

    <table class="table">
      <thead>
        <tr>
          <th>ID</th>
          <th>File Name</th>
          <th></th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="file in files">
          <td>{{ file.id }}</td>
          <td>{{ file.file_name }}</td>
          <td><a v-bind:href="'/files/' + file.id">Download</a></td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script>
import { APIClient } from './api_client.ts';

export default {
  name: 'app',
  data() {
    return {
      searchQuery: '',
      files: [],
    };
  },
  methods: {
    search: function() {
      const self = this;
      APIClient.search(this.searchQuery, (result, err) => {
        if (err !== null) {
          // TODO error handling
          return;
        }
        self.files = result;
      });
    },
  },
};
</script>

<style lang="scss">
#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  margin-top: 10px;
}

.searchbox {
  margin-top: 20px;
}

h1,
h2 {
  font-weight: normal;
}

ul {
  list-style-type: none;
  padding: 0;
}

li {
  display: inline-block;
  margin: 0 10px;
}

a {
  color: #42b983;
}
</style>

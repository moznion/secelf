<template>
  <div id="app">
    <h1>SECELF</h1>

    <form action="/api/files" method="POST" enctype="multipart/form-data">
      <div class="input-group">
        <div class="custom-file">
          <input type="file" class="custom-file-input" id="file">
          <label class="custom-file-label" for="file">Choose file</label>
        </div>
        <div class="input-group-append">
          <button class="btn btn-primary" type="submit">Upload</button>
        </div>
      </div>
    </form>
  </div>
</template>

<script>
export default {
  name: 'app',
  data() {
    return {
      msg: 'Welcome to Your Vue.js App',
    };
  },
  methods: {
    foo: function() {
      this.msg = 'hih';
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

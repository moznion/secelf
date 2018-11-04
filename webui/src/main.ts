import BootstrapVue from 'bootstrap-vue';
import 'bootstrap-vue/dist/bootstrap-vue.css';
import 'bootstrap/dist/css/bootstrap.min.css';
import Vue from 'vue';
import App from './App.vue';

Vue.use(BootstrapVue);

/* tslint:disable:no-unused-expression */
new Vue({
  el: '#app',
  render: h => h(App),
});

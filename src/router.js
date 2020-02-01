import VueRouter from 'vue-router'
import Vue from 'vue'
import Recover from './views/Recover.vue'

Vue.use(VueRouter)

const routes = [
    { path: '/', component: Recover },
    { path: '/recover', component: null },
    { path: '/new', component: null },
    { path: '/about', component: null },
    { path: '/send', component: null },
    { path: '/receive', component: null },
    { path: '/history', component: null },
  ];

export default new VueRouter({
    mode: 'history',
    routes
});
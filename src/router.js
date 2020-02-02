import VueRouter from 'vue-router'
import Vue from 'vue'
import Recover from './views/Recover.vue'
import New from './views/New.vue'
import About from './views/About.vue'
import Send from './views/Send.vue'
import Receive from './views/Receive.vue'

Vue.use(VueRouter)

const routes = [
    { path: '/', component: Recover },
    { path: '/new', component: New },
    { path: '/about', component: About },
    { path: '/send', component: Send },
    { path: '/receive', component: Receive },
    { path: '/history', component: null },
  ];

export default new VueRouter({
    mode: 'history',
    routes
});
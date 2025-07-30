import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import SetsView from '../views/SetsView.vue'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: HomeView,
  },
  {
    path: '/sets',
    name: 'Sets',
    component: SetsView,
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router

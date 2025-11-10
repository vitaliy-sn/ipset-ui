import { createRouter, createWebHistory } from 'vue-router'
import EntriesView from '../views/EntriesView.vue'
import SetsView from '../views/SetsView.vue'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: EntriesView,
  },
  {
    path: '/sets',
    name: 'Sets',
    component: SetsView,
  },
  {
    path: '/entries',
    name: 'Entries',
    component: EntriesView,
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router

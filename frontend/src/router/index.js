import { createRouter, createWebHistory } from 'vue-router'
import WordsView from '../views/WordsView.vue'

const routes = [
  {
    path: '/',
    name: 'Words',
    component: WordsView
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
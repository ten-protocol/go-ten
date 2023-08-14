import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView
    },
    {
      path: '/about',
      name: 'about',
      // route level code-splitting
      // this generates a separate chunk (About.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      component: () => import('../views/AboutView.vue')
    },
    {
      path: '/personal',
      name: 'personal',
      component: () => import('../views/PersonalView.vue')
    },
    {
      path: '/transactions',
      name: 'transactions',
      component: () => import('../views/TransactionsView.vue')
    },
    {
      path: '/batches',
      name: 'batches',
      component: () => import('../views/BatchesView.vue')
    },
    {
      path: '/rollups',
      name: 'rollups',
      component: () => import('../views/RollupsView.vue')
    }
  ]
})

export default router

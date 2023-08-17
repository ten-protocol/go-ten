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
      path: '/blocks',
      name: 'blocks',
      component: () => import('../views/BlocksView.vue')
    },
    {
      path: '/decrypt',
      name: 'decrypt',
      component: () => import('../views/DecryptView.vue')
    },
    {
      path: '/verified',
      name: 'verifiedContracts',
      component: () => import('../views/VerifiedContracts.vue')
    }
  ]
})

export default router

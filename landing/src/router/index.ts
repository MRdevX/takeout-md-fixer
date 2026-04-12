import { createRouter, createWebHistory } from 'vue-router'

import HomeView from '@/views/HomeView.vue'
import PrivacyView from '@/views/PrivacyView.vue'

export const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    { path: '/', name: 'home', component: HomeView },
    { path: '/privacy', name: 'privacy', component: PrivacyView },
  ],
  scrollBehavior() {
    return { top: 0 }
  },
})

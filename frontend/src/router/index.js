import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/Login.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('../views/Register.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: () => import('../views/DashboardView.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/',
    redirect: to => {
      const isAuthenticated = !!localStorage.getItem('access_token')
      return isAuthenticated ? '/dashboard' : '/login'
    }
  },
  {
    path: '/rankings',
    name: 'Rankings',
    component: () => import('../views/RankingsView.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/statistics',
    name: 'Statistics',
    component: () => import('../views/StatisticsView.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/profile',
    name: 'Profile',
    component: () => import('../views/Profile.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/profile/security',
    name: 'SecuritySettings',
    component: () => import('../views/SecuritySettings.vue'),
    meta: { requiresAuth: true }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// Navigation guard
router.beforeEach((to, from, next) => {
  const isAuthenticated = !!localStorage.getItem('access_token')
  
  if (to.meta.requiresAuth) {
    if (!isAuthenticated) {
      // Save the attempted URL for redirection after login
      next({ 
        name: 'Login',
        query: { redirect: to.fullPath }
      })
    } else {
      // Check if token is valid by checking its expiration
      const token = localStorage.getItem('access_token')
      try {
        const base64Url = token.split('.')[1]
        const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/')
        const payload = JSON.parse(window.atob(base64))
        
        // Check if token is expired
        if (payload.exp && payload.exp * 1000 < Date.now()) {
          // Token is expired, remove it and redirect to login
          localStorage.removeItem('access_token')
          localStorage.removeItem('refresh_token')
          next({ 
            name: 'Login',
            query: { redirect: to.fullPath }
          })
        } else {
          next()
        }
      } catch (e) {
        // If token is invalid, remove it and redirect to login
        localStorage.removeItem('access_token')
        localStorage.removeItem('refresh_token')
        next({ 
          name: 'Login',
          query: { redirect: to.fullPath }
        })
      }
    }
  } else if (isAuthenticated && ['/login', '/register'].includes(to.path)) {
    next({ name: 'Dashboard' })
  } else {
    next()
  }
})

export default router

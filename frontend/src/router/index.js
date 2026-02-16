import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: () => import('../views/HomeView.vue')
  },
  {
    path: '/lessons',
    name: 'Lessons',
    component: () => import('../views/LessonsView.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/lesson/:id',
    name: 'Lesson',
    component: () => import('../views/LessonView.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/words',
    name: 'Words',
    component: () => import('../views/WordsView.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/LoginView.vue')
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('../views/RegisterView.vue')
  },
  {
    path: '/teacher/login',
    name: 'TeacherLogin',
    component: () => import('../views/TeacherLoginView.vue')
  },
  {
    path: '/teacher/dashboard',
    name: 'TeacherDashboard',
    component: () => import('../views/TeacherDashboardView.vue'),
    meta: { requiresTeacherAuth: true }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// Navigation guard для проверки аутентификации
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  const teacherToken = localStorage.getItem('teacher_token')
  
  // Если пользователь авторизован и пытается попасть на главную, перенаправляем на уроки
  if (to.path === '/' && (token || teacherToken)) {
    next('/lessons')
    return
  }
  
  // Проверка обычной аутентификации
  if (to.meta.requiresAuth && !token) {
    next('/login')
    return
  }
  
  // Проверка аутентификации учителя
  if (to.meta.requiresTeacherAuth && !teacherToken) {
    next('/teacher/login')
    return
  }
  
  // Если учитель пытается зайти в обычный раздел
  if (teacherToken && !to.path.startsWith('/teacher') && to.path !== '/login' && to.path !== '/register') {
    if (to.path !== '/lessons' && to.path !== '/lesson/:id' && to.path !== '/words') {
      next('/teacher/dashboard')
      return
    }
  }
  
  // Если обычный пользователь пытается зайти в раздел учителя
  if (token && to.path.startsWith('/teacher')) {
    next('/lessons')
    return
  }
  
  next()
})

export default router
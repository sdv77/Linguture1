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
    path: '/profile',
    name: 'Profile',
    component: () => import('../views/ProfileView.vue'),
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
  // 🔹 НОВЫЙ: Страница настройки профиля
  {
    path: '/setup-profile',
    name: 'ProfileSetup',
    component: () => import('../views/ProfileSetupView.vue'),
    meta: { requiresAuth: true }
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
router.beforeEach(async (to, from, next) => {
  const token = localStorage.getItem('token')
  const teacherToken = localStorage.getItem('teacher_token')
  
  // Если пользователь авторизован и пытается попасть на главную, перенаправляем на уроки
  if (to.path === '/' && (token || teacherToken)) {
    next('/lessons')
    return
  }
  
  // 🔹 ПРОВЕРКА: если токен есть, но профиль не настроен
  if (token && to.path !== '/setup-profile' && to.path !== '/login' && to.path !== '/register') {
    const isSetup = localStorage.getItem('profile_is_setup')
    
    // Если не знаем статус или профиль не настроен → проверяем через API
    if (isSetup === 'false' || isSetup === null) {
      try {
        // Делаем запрос к API для получения актуальных данных
        const response = await fetch('http://localhost:8080/api/user/me', {
          headers: {
            'Authorization': `Bearer ${token}`
          }
        })
        
        if (response.ok) {
          const userData = await response.json()
          localStorage.setItem('profile_is_setup', userData.is_setup)
          
          // Если профиль действительно не настроен → перенаправляем
          if (!userData.is_setup && to.path !== '/setup-profile') {
            next('/setup-profile')
            return
          }
        }
      } catch (error) {
        console.error('Ошибка проверки профиля:', error)
        // Не прерываем навигацию при ошибке
      }
    }
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
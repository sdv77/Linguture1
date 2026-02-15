<template>
  <div id="app">
    <nav v-if="showNavbar" class="navbar">
      <div class="nav-container">
        <router-link to="/" class="nav-logo">Words App</router-link>
        
        <div v-if="isLoggedIn" class="nav-links">
          <router-link to="/lessons" class="nav-link">📚 Уроки</router-link>
          <router-link to="/" class="nav-link">📖 Мои слова</router-link>
          <button @click="logout" class="nav-link nav-logout">🚪 Выйти</button>
        </div>
        
        <div v-if="isTeacher" class="nav-links">
          <router-link to="/teacher/dashboard" class="nav-link">👨‍🏫 Панель учителя</router-link>
          <button @click="logoutTeacher" class="nav-link nav-logout">🚪 Выйти</button>
        </div>
      </div>
    </nav>
    
    <router-view />
  </div>
</template>

<script>
import { computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'

export default {
  setup() {
    const router = useRouter()
    const isLoggedIn = computed(() => !!localStorage.getItem('token'))
    const isTeacher = computed(() => !!localStorage.getItem('teacher_token'))
    
    // Показывать навигацию только на определённых страницах
    const showNavbar = computed(() => {
      const currentPath = router.currentRoute.value.path
      return !currentPath.startsWith('/login') && 
             !currentPath.startsWith('/register') && 
             !currentPath.startsWith('/teacher/login')
    })
    
    const logout = () => {
      localStorage.removeItem('token')
      router.push('/login')
    }
    
    const logoutTeacher = () => {
      localStorage.removeItem('teacher_token')
      router.push('/teacher/login')
    }
    
    return {
      isLoggedIn,
      isTeacher,
      showNavbar,
      logout,
      logoutTeacher
    }
  }
}
</script>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
  background-color: #f5f7fa;
  color: #333;
  line-height: 1.6;
}

#app {
  min-height: 100vh;
}

.navbar {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  position: sticky;
  top: 0;
  z-index: 1000;
}

.nav-container {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 60px;
}

.nav-logo {
  color: white;
  font-size: 24px;
  font-weight: bold;
  text-decoration: none;
}

.nav-links {
  display: flex;
  gap: 30px;
  align-items: center;
}

.nav-link {
  color: white;
  text-decoration: none;
  font-size: 16px;
  font-weight: 500;
  padding: 8px 16px;
  border-radius: 6px;
  transition: all 0.2s;
}

.nav-link:hover {
  background: rgba(255, 255, 255, 0.2);
}

.nav-link.router-link-active {
  background: rgba(255, 255, 255, 0.3);
}

.nav-logout {
  background: none;
  border: none;
  cursor: pointer;
  color: white;
  font-size: 16px;
  font-weight: 500;
  padding: 8px 16px;
  border-radius: 6px;
  transition: all 0.2s;
}

.nav-logout:hover {
  background: rgba(255, 255, 255, 0.2);
}
</style>
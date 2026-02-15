<template>
  <div class="auth-container">
    <div class="auth-card">
      <div class="auth-header">
        <h2>👨‍🏫 Вход для учителя</h2>
        <p>Панель управления уроками</p>
      </div>
      
      <form @submit.prevent="login">
        <div class="form-group">
          <label for="email">Email:</label>
          <input 
            type="email" 
            id="email" 
            v-model="formData.email" 
            required 
            placeholder="Введите ваш email"
          />
        </div>
        
        <div class="form-group">
          <label for="password">Пароль:</label>
          <input 
            type="password" 
            id="password" 
            v-model="formData.password" 
            required 
            placeholder="Введите пароль"
          />
        </div>
        
        <div v-if="errorMessage" class="error-message">
          {{ errorMessage }}
        </div>
        
        <button type="submit" class="btn btn-primary" :disabled="loading">
          {{ loading ? 'Вход...' : 'Войти' }}
        </button>
      </form>
      
      <div class="auth-footer">
        <p>Тестовый учитель: teacher@wordsapp.com / teacher123</p>
        <router-link to="/login" class="back-link">← Вернуться на обычный вход</router-link>
      </div>
    </div>
  </div>
</template>

<script>
import { ref } from 'vue'
import { useRouter } from 'vue-router'

export default {
  name: 'TeacherLoginView',
  
  setup() {
    const router = useRouter()
    const formData = ref({
      email: '',
      password: ''
    })
    const loading = ref(false)
    const errorMessage = ref('')
    
    const API_URL = 'http://localhost:8080/api/teacher/auth/login'
    
    const login = async () => {
      loading.value = true
      errorMessage.value = ''
      
      try {
        const response = await fetch(API_URL, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({
            email: formData.value.email,
            password: formData.value.password
          })
        })
        
        if (!response.ok) {
          const errorData = await response.json().catch(() => ({}))
          throw new Error(errorData.message || 'Ошибка входа')
        }
        
        const data = await response.json()
        
        localStorage.setItem('teacher_token', data.token)
        
        router.push('/teacher/dashboard')
      } catch (error) {
        console.error('Ошибка входа:', error)
        errorMessage.value = error.message || 'Ошибка при входе'
      } finally {
        loading.value = false
      }
    }
    
    return {
      formData,
      loading,
      errorMessage,
      login
    }
  }
}
</script>

<style scoped>
.auth-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
}

.auth-card {
  background: white;
  padding: 40px;
  border-radius: 12px;
  box-shadow: 0 15px 30px rgba(0, 0, 0, 0.2);
  width: 100%;
  max-width: 450px;
}

.auth-header {
  text-align: center;
  margin-bottom: 30px;
}

.auth-header h2 {
  margin: 0 0 10px 0;
  color: #333;
  font-size: 28px;
}

.auth-header p {
  margin: 0;
  color: #666;
  font-size: 16px;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  font-weight: 600;
  color: #555;
  font-size: 15px;
}

.form-group input {
  width: 100%;
  padding: 14px;
  border: 2px solid #e0e0e0;
  border-radius: 8px;
  font-size: 16px;
  transition: border-color 0.3s;
}

.form-group input:focus {
  outline: none;
  border-color: #667eea;
}

.error-message {
  background-color: #ffebee;
  color: #c62828;
  padding: 12px;
  border-radius: 8px;
  margin-bottom: 15px;
  border-left: 4px solid #c62828;
  font-size: 14px;
}

.btn {
  width: 100%;
  padding: 14px;
  border: none;
  border-radius: 8px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s;
}

.btn-primary {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.btn-primary:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 8px 16px rgba(102, 126, 234, 0.4);
}

.btn-primary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.auth-footer {
  text-align: center;
  margin-top: 25px;
  color: #666;
}

.auth-footer p {
  margin: 0 0 10px 0;
  font-size: 14px;
  color: #999;
}

.back-link {
  color: #667eea;
  text-decoration: none;
  font-weight: 600;
  font-size: 15px;
}

.back-link:hover {
  text-decoration: underline;
}
</style>
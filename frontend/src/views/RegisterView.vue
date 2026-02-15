<template>
  <div class="auth-container">
    <div class="auth-card">
      <h2>Регистрация</h2>
      <form @submit.prevent="register">
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
            placeholder="Введите пароль (минимум 6 символов)"
            minlength="6"
          />
        </div>
        
        <div class="form-group">
          <label for="confirmPassword">Подтвердите пароль:</label>
          <input 
            type="password" 
            id="confirmPassword" 
            v-model="formData.confirmPassword" 
            required 
            placeholder="Повторите пароль"
          />
        </div>
        
        <div v-if="errorMessage" class="error-message">
          {{ errorMessage }}
        </div>
        
        <button type="submit" class="btn btn-primary" :disabled="loading">
          {{ loading ? 'Регистрация...' : 'Зарегистрироваться' }}
        </button>
      </form>
      
      <div class="auth-footer">
        <p>Уже есть аккаунт? <router-link to="/login">Войти</router-link></p>
      </div>
    </div>
  </div>
</template>

<script>
import { ref } from 'vue';
import { useRouter } from 'vue-router';

export default {
  name: 'RegisterView',
  
  setup() {
    const router = useRouter();
    const formData = ref({
      email: '',
      password: '',
      confirmPassword: ''
    });
    const loading = ref(false);
    const errorMessage = ref('');
    
    const API_URL = 'http://localhost:8080/api/auth/register';
    
    const register = async () => {
      // Валидация
      if (formData.value.password !== formData.value.confirmPassword) {
        errorMessage.value = 'Пароли не совпадают';
        return;
      }
      
      if (formData.value.password.length < 6) {
        errorMessage.value = 'Пароль должен содержать минимум 6 символов';
        return;
      }
      
      loading.value = true;
      errorMessage.value = '';
      
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
        });
        
        if (!response.ok) {
          const errorData = await response.json();
          throw new Error(errorData.message || 'Ошибка регистрации');
        }
        
        // Регистрация успешна, переходим на страницу входа
        alert('Регистрация успешна! Проверьте почту для подтверждения.');
        router.push('/login');
      } catch (error) {
        console.error('Ошибка регистрации:', error);
        errorMessage.value = error.message || 'Ошибка при регистрации';
      } finally {
        loading.value = false;
      }
    };
    
    return {
      formData,
      loading,
      errorMessage,
      register
    };
  }
};
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
  border-radius: 10px;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.2);
  width: 100%;
  max-width: 400px;
}

.auth-card h2 {
  text-align: center;
  margin-bottom: 30px;
  color: #333;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  font-weight: 600;
  color: #555;
}

.form-group input {
  width: 100%;
  padding: 12px;
  border: 1px solid #ddd;
  border-radius: 6px;
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
  border-radius: 6px;
  margin-bottom: 15px;
  border-left: 4px solid #c62828;
}

.btn {
  width: 100%;
  padding: 12px;
  border: none;
  border-radius: 6px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: background-color 0.3s;
}

.btn-primary {
  background-color: #667eea;
  color: white;
}

.btn-primary:hover:not(:disabled) {
  background-color: #5568d3;
}

.btn-primary:disabled {
  background-color: #a0a8e0;
  cursor: not-allowed;
}

.auth-footer {
  text-align: center;
  margin-top: 20px;
  color: #666;
}

.auth-footer a {
  color: #667eea;
  text-decoration: none;
  font-weight: 600;
}

.auth-footer a:hover {
  text-decoration: underline;
}
</style>
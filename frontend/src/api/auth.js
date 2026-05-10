// src/api/auth.js

const API_BASE_URL = `${import.meta.env.VITE_API_URL || 'http://localhost:8080'}/api`

// Вспомогательная функция для выполнения запросов
async function request(url, options = {}) {
  const token = localStorage.getItem('token')
  
  const config = {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...(token && { 'Authorization': `Bearer ${token}` }),
      ...options.headers,
    },
  }

  const response = await fetch(`${API_BASE_URL}${url}`, config)
  
  if (!response.ok) {
    const error = await response.text()
    throw new Error(error || `HTTP error! status: ${response.status}`)
  }
  
  return response.json()
}

// 🔹 Логин пользователя
export async function login(email, password) {
  return request('/auth/login', {
    method: 'POST',
    body: JSON.stringify({ email, password }),
  })
}

// 🔹 Регистрация пользователя
export async function register(email, password) {
  return request('/auth/register', {
    method: 'POST',
    body: JSON.stringify({ email, password }),
  })
}

// 🔹 Первичная настройка профиля
export async function setupProfile(nickname, nativeLanguage, learningLanguage) {
  return request('/user/setup', {
    method: 'POST',
    body: JSON.stringify({
      nickname,
      native_language: nativeLanguage,
      learning_language: learningLanguage,
    }),
  })
}

// 🔹 Получение данных текущего пользователя
export async function getCurrentUser() {
  return request('/user/me')
}

<template>
  <div class="profile-container">
    <div class="profile-card">
      <!-- 🔹 Аватарка -->
      <div class="avatar-section">
        <div class="avatar">👤</div>
        <h1 class="nickname">{{ user.nickname || 'Пользователь' }}</h1>
      </div>

      <!-- 🔹 Информация о языках -->
      <div class="languages-section">
        <div class="language-item native">
          <div class="language-icon">🏠</div>
          <div class="language-info">
            <span class="language-label">Родной язык</span>
            <span class="language-value">
              {{ getLanguageEmoji(user.native_language) }} {{ getLanguageName(user.native_language) }}
            </span>
          </div>
        </div>

        <div class="language-arrow">↓</div>

        <div class="language-item learning">
          <div class="language-icon">🎯</div>
          <div class="language-info">
            <span class="language-label">Изучаю</span>
            <span class="language-value">
              {{ getLanguageEmoji(user.learning_language) }} {{ getLanguageName(user.learning_language) }}
            </span>
          </div>
        </div>
      </div>

      <!-- 🔹 Email -->
      <div class="email-section">
        <span class="email-label">📧 Email:</span>
        <span class="email-value">{{ user.email }}</span>
      </div>

      <!-- 🔹 Кнопка выхода (еле заметная) -->
      <button @click="logout" class="logout-btn">
        <span class="logout-icon">🚪</span>
        <span>Выйти</span>
      </button>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getCurrentUser } from '../api/auth'

export default {
  name: 'ProfileView',
  setup() {
    const router = useRouter()
    const user = ref({
      nickname: '',
      email: '',
      native_language: '',
      learning_language: '',
    })

    // 🔹 Список языков для отображения
    const languages = {
      ru: { name: 'Русский', emoji: '🇷🇺' },
      en: { name: 'English', emoji: '🇺🇸' },
      es: { name: 'Español', emoji: '🇪' },
      de: { name: 'Deutsch', emoji: '🇩🇪' },
      fr: { name: 'Français', emoji: '🇫' },
      it: { name: 'Italiano', emoji: '🇮' },
      pt: { name: 'Português', emoji: '🇵🇹' },
      zh: { name: '中文', emoji: '🇨' },
      ja: { name: '日本語', emoji: '🇯🇵' },
      ko: { name: '한국어', emoji: '🇰🇷' },
      ar: { name: 'العربية', emoji: '🇸🇦' },
      tr: { name: 'Türkçe', emoji: '🇹🇷' },
    }

    const getLanguageName = (code) => {
      return languages[code]?.name || code
    }

    const getLanguageEmoji = (code) => {
      return languages[code]?.emoji || ''
    }

    const logout = () => {
      localStorage.removeItem('token')
      localStorage.removeItem('profile_is_setup')
      localStorage.removeItem('user_id')
      router.push('/')
    }

    onMounted(async () => {
      try {
        const userData = await getCurrentUser()
        user.value = userData
      } catch (error) {
        console.error('Ошибка загрузки профиля:', error)
        // Если не удалось загрузить, используем данные из localStorage
        user.value = {
          nickname: localStorage.getItem('nickname') || '',
          email: localStorage.getItem('email') || '',
          native_language: localStorage.getItem('native_language') || '',
          learning_language: localStorage.getItem('learning_language') || '',
        }
      }
    })

    return {
      user,
      getLanguageName,
      getLanguageEmoji,
      logout,
    }
  },
}
</script>

<style scoped>
.profile-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.profile-card {
  background: white;
  border-radius: 24px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  padding: 50px 40px;
  max-width: 500px;
  width: 100%;
  text-align: center;
  position: relative;
}

/* 🔹 Аватарка и никнейм */
.avatar-section {
  margin-bottom: 40px;
}

.avatar {
  width: 120px;
  height: 120px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 64px;
  margin: 0 auto 20px;
  box-shadow: 0 8px 20px rgba(102, 126, 234, 0.3);
}

.nickname {
  font-size: 28px;
  color: #333;
  font-weight: bold;
  margin: 0;
}

/* 🔹 Секция языков */
.languages-section {
  background: #f8f9fa;
  border-radius: 16px;
  padding: 30px;
  margin-bottom: 30px;
}

.language-item {
  display: flex;
  align-items: center;
  gap: 15px;
  padding: 15px;
  border-radius: 12px;
  margin-bottom: 15px;
  transition: all 0.3s;
}

.language-item:last-child {
  margin-bottom: 0;
}

.language-item.native {
  background: white;
  border: 2px solid #e0e0e0;
}

.language-item.learning {
  background: linear-gradient(135deg, #4CAF50 0%, #45a049 100%);
  color: white;
  border: 2px solid #4CAF50;
}

.language-icon {
  font-size: 32px;
}

.language-info {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  flex: 1;
}

.language-label {
  font-size: 12px;
  opacity: 0.7;
  margin-bottom: 4px;
}

.language-value {
  font-size: 18px;
  font-weight: 600;
}

.language-arrow {
  font-size: 24px;
  color: #667eea;
  margin: 10px 0;
  animation: bounce 1s ease infinite;
}

@keyframes bounce {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-5px); }
}

/* 🔹 Email */
.email-section {
  margin-bottom: 30px;
  padding: 15px;
  background: #f8f9fa;
  border-radius: 12px;
}

.email-label {
  display: block;
  font-size: 12px;
  color: #999;
  margin-bottom: 5px;
}

.email-value {
  font-size: 16px;
  color: #666;
  word-break: break-all;
}

/* 🔹 Кнопка выхода (еле заметная) */
.logout-btn {
  background: transparent;
  border: 1px solid rgba(102, 126, 234, 0.2);
  color: rgba(102, 126, 234, 0.5);
  padding: 10px 20px;
  border-radius: 8px;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.3s;
  display: inline-flex;
  align-items: center;
  gap: 8px;
  opacity: 0.4;
}

.logout-btn:hover {
  background: rgba(231, 76, 60, 0.1);
  border-color: #e74c3c;
  color: #e74c3c;
  opacity: 1;
  transform: translateY(-2px);
}

.logout-icon {
  font-size: 16px;
}

/* 🔹 Адаптивность */
@media (max-width: 640px) {
  .profile-card {
    padding: 30px 20px;
  }

  .avatar {
    width: 100px;
    height: 100px;
    font-size: 48px;
  }

  .nickname {
    font-size: 22px;
  }

  .languages-section {
    padding: 20px;
  }

  .language-value {
    font-size: 16px;
  }
}
</style>
<template>
  <div class="setup-container">
    <div class="setup-card">
      <div class="setup-header">
        <div class="setup-icon">🎉</div>
        <h1 class="setup-title">Добро пожаловать в Linguture!</h1>
        <p class="setup-subtitle">Давайте настроим ваш профиль для начала обучения</p>
      </div>

      <form @submit.prevent="handleSubmit" class="setup-form">
        <!-- 🔹 Поле никнейма -->
        <div class="form-group">
          <label for="nickname" class="form-label">
            <span class="label-icon">👤</span>
            Придумайте никнейм
          </label>
          <input
            id="nickname"
            v-model="form.nickname"
            type="text"
            class="form-input"
            placeholder="Например: ivan_learner"
            :class="{ 'input-error': errors.nickname }"
            required
            minlength="3"
            maxlength="30"
          />
          <p v-if="errors.nickname" class="error-message">{{ errors.nickname }}</p>
          <p class="form-hint">3-30 символов, только латинские буквы, цифры и подчёркивание</p>
        </div>

        <!-- 🔹 Выбор родного языка -->
        <div class="form-group">
          <label class="form-label">
            <span class="label-icon">🏠</span>
            Ваш родной язык
          </label>
          <div class="language-grid">
            <div
              v-for="lang in languages"
              :key="lang.code"
              class="language-option"
              :class="{ selected: form.nativeLanguage === lang.code }"
              @click="form.nativeLanguage = lang.code"
            >
              <input
                type="radio"
                :id="`native-${lang.code}`"
                :value="lang.code"
                v-model="form.nativeLanguage"
                class="language-radio"
              />
              <label :for="`native-${lang.code}`" class="language-label">
                <span class="flag">{{ lang.emoji }}</span>
                <span class="lang-name">{{ lang.name }}</span>
              </label>
            </div>
          </div>
          <p v-if="errors.nativeLanguage" class="error-message">{{ errors.nativeLanguage }}</p>
        </div>

        <!-- 🔹 Выбор изучаемого языка -->
        <div class="form-group">
          <label class="form-label">
            <span class="label-icon">🎯</span>
            Какой язык хотите изучать?
          </label>
          <div class="language-grid">
            <div
              v-for="lang in languages"
              :key="`learn-${lang.code}`"
              class="language-option"
              :class="{ 
                selected: form.learningLanguage === lang.code,
                disabled: form.nativeLanguage === lang.code
              }"
              @click="selectLearningLanguage(lang.code)"
            >
              <input
                type="radio"
                :id="`learn-${lang.code}`"
                :value="lang.code"
                v-model="form.learningLanguage"
                class="language-radio"
                :disabled="form.nativeLanguage === lang.code"
              />
              <label :for="`learn-${lang.code}`" class="language-label">
                <span class="flag">{{ lang.emoji }}</span>
                <span class="lang-name">{{ lang.name }}</span>
                <span v-if="form.nativeLanguage === lang.code" class="lang-badge">Родной</span>
              </label>
            </div>
          </div>
          <p v-if="errors.learningLanguage" class="error-message">{{ errors.learningLanguage }}</p>
        </div>

        <!-- 🔹 Кнопка отправки -->
        <button type="submit" class="submit-btn" :disabled="loading">
          <span v-if="loading" class="loading-spinner">⏳</span>
          <span v-else>🚀</span>
          {{ loading ? 'Сохранение...' : 'Начать обучение' }}
        </button>
      </form>
    </div>
  </div>
</template>

<script>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { setupProfile } from '../api/auth'

export default {
  name: 'ProfileSetupView',
  setup() {
    const router = useRouter()
    const loading = ref(false)
    const errors = reactive({})
    
    // 🔹 Список поддерживаемых языков с эмодзи
    const languages = [
      { code: 'ru', name: 'Русский', emoji: '🇷🇺' },
      { code: 'en', name: 'English', emoji: '🇺🇸' },
      { code: 'es', name: 'Español', emoji: '🇪🇸' },
      { code: 'de', name: 'Deutsch', emoji: '🇩🇪' },
      { code: 'fr', name: 'Français', emoji: '🇫🇷' },
      { code: 'it', name: 'Italiano', emoji: '🇮🇹' },
      { code: 'pt', name: 'Português', emoji: '🇵🇹' },
      { code: 'zh', name: '中文', emoji: '🇨🇳' },
      { code: 'ja', name: '日本語', emoji: '🇯🇵' },
      { code: 'ko', name: '한국어', emoji: '🇰🇷' },
      { code: 'ar', name: 'العربية', emoji: '🇸🇦' },
      { code: 'tr', name: 'Türkçe', emoji: '🇹🇷' },
    ]

    // 🔹 Данные формы
    const form = reactive({
      nickname: '',
      nativeLanguage: '',
      learningLanguage: '',
    })

    // 🔹 Выбор изучаемого языка (исключаем родной)
    const selectLearningLanguage = (code) => {
      if (form.nativeLanguage !== code) {
        form.learningLanguage = code
      }
    }

    // 🔹 Валидация никнейма
    const validateNickname = (nickname) => {
      const regex = /^[a-zA-Z0-9_]{3,30}$/
      return regex.test(nickname)
    }

    // 🔹 Отправка формы
    const handleSubmit = async () => {
      // Сбрасываем ошибки
      Object.keys(errors).forEach(key => delete errors[key])
      
      // Валидация
      if (!form.nickname) {
        errors.nickname = 'Никнейм обязателен'
        return
      }
      if (!validateNickname(form.nickname)) {
        errors.nickname = 'Никнейм может содержать только латинские буквы, цифры и подчёркивание (3-30 символов)'
        return
      }
      if (!form.nativeLanguage) {
        errors.nativeLanguage = 'Выберите родной язык'
        return
      }
      if (!form.learningLanguage) {
        errors.learningLanguage = 'Выберите язык для изучения'
        return
      }
      if (form.nativeLanguage === form.learningLanguage) {
        errors.learningLanguage = 'Язык изучения не может совпадать с родным'
        return
      }

      loading.value = true

      try {
        await setupProfile(
          form.nickname,
          form.nativeLanguage,
          form.learningLanguage
        )
        
        // ✅ Успешно! Перенаправляем на уроки
        router.push('/lessons')
      } catch (error) {
        console.error('Ошибка настройки профиля:', error)
        
        // Обработка ошибок
        if (error.message.includes('никнейм уже занят')) {
          errors.nickname = 'Этот никнейм уже занят, выберите другой'
        } else if (error.message.includes('неподдерживаемый')) {
          errors.nativeLanguage = 'Выбранный язык не поддерживается'
        } else {
          alert('Произошла ошибка при сохранении профиля. Попробуйте ещё раз.')
        }
      } finally {
        loading.value = false
      }
    }

    return {
      form,
      languages,
      loading,
      errors,
      selectLearningLanguage,
      handleSubmit,
    }
  },
}
</script>

<style scoped>
.setup-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.setup-card {
  background: white;
  border-radius: 24px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  padding: 50px 40px;
  max-width: 700px;
  width: 100%;
}

.setup-header {
  text-align: center;
  margin-bottom: 40px;
}

.setup-icon {
  font-size: 64px;
  margin-bottom: 20px;
  animation: bounce 1s ease;
}

@keyframes bounce {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-10px); }
}

.setup-title {
  font-size: 32px;
  color: #333;
  margin-bottom: 10px;
  font-weight: bold;
}

.setup-subtitle {
  font-size: 18px;
  color: #666;
}

.setup-form {
  display: flex;
  flex-direction: column;
  gap: 30px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.form-label {
  font-size: 16px;
  font-weight: 600;
  color: #333;
  display: flex;
  align-items: center;
  gap: 8px;
}

.label-icon {
  font-size: 20px;
}

.form-input {
  padding: 14px 18px;
  border: 2px solid #e0e0e0;
  border-radius: 12px;
  font-size: 16px;
  transition: all 0.3s;
  outline: none;
}

.form-input:focus {
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.form-input.input-error {
  border-color: #e74c3c;
}

.form-hint {
  font-size: 13px;
  color: #999;
  margin-top: 5px;
}

.error-message {
  color: #e74c3c;
  font-size: 14px;
  margin-top: 5px;
}

/* 🔹 Сетка языков */
.language-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(140px, 1fr));
  gap: 12px;
  margin-top: 10px;
}

.language-option {
  position: relative;
  border: 2px solid #e0e0e0;
  border-radius: 12px;
  overflow: hidden;
  cursor: pointer;
  transition: all 0.3s;
}

.language-option:hover:not(.disabled) {
  border-color: #667eea;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.2);
}

.language-option.selected {
  border-color: #667eea;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.language-option.disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.language-radio {
  position: absolute;
  opacity: 0;
  pointer-events: none;
}

.language-label {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 16px 12px;
  cursor: pointer;
  gap: 8px;
}

.flag {
  font-size: 32px;
}

.lang-name {
  font-size: 14px;
  font-weight: 500;
  text-align: center;
}

.lang-badge {
  font-size: 10px;
  background: rgba(0, 0, 0, 0.1);
  padding: 2px 8px;
  border-radius: 10px;
  margin-top: 4px;
}

/* 🔹 Кнопка отправки */
.submit-btn {
  background: linear-gradient(135deg, #4CAF50 0%, #45a049 100%);
  color: white;
  border: none;
  padding: 18px 40px;
  border-radius: 12px;
  font-size: 18px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  margin-top: 20px;
}

.submit-btn:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(76, 175, 80, 0.4);
}

.submit-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.loading-spinner {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

/* 🔹 Адаптивность */
@media (max-width: 640px) {
  .setup-card {
    padding: 30px 20px;
  }
  
  .setup-title {
    font-size: 24px;
  }
  
  .language-grid {
    grid-template-columns: repeat(auto-fill, minmax(100px, 1fr));
  }
  
  .flag {
    font-size: 24px;
  }
  
  .lang-name {
    font-size: 12px;
  }
}
</style>
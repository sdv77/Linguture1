<template>
  <div class="lesson-container">
    <div class="lesson-header">
      <button @click="goBack" class="btn btn-back">← Назад</button>
      <h1>{{ lesson?.lesson?.title }}</h1>
      <div class="lesson-info">
        <span class="lesson-level">Уровень {{ lesson?.lesson?.level }}</span>
        <span class="lesson-progress">{{ currentExercise + 1 }}/{{ lesson?.words?.length || 0 }}</span>
      </div>
    </div>
    
    <div v-if="loading" class="loading">Загрузка урока...</div>
    
    <div v-else-if="lesson" class="lesson-content">
      <div v-if="currentExercise < lesson.words.length" class="exercise-container">
        <div class="exercise-card">
          <div class="word-display">
            <h2>{{ currentWord.word }}</h2>
            <p v-if="currentWord.transcription" class="transcription">{{ currentWord.transcription }}</p>
          </div>
          
          <div class="exercise-type">
            <h3>Выберите правильный перевод:</h3>
          </div>
          
          <div class="options-grid">
            <button 
              v-for="(option, index) in shuffledOptions" 
              :key="index"
              @click="checkAnswer(option)"
              class="option-btn"
              :class="{
                'correct': selectedOption === option && option.correct,
                'incorrect': selectedOption === option && !option.correct,
                'disabled': selectedOption !== null
              }"
            >
              {{ option.text }}
            </button>
          </div>
          
          <div v-if="feedback" class="feedback" :class="feedback.type">
            {{ feedback.message }}
          </div>
          
          <button 
            v-if="showNextButton" 
            @click="nextExercise" 
            class="btn btn-primary btn-next"
          >
            {{ currentExercise + 1 === lesson.words.length ? 'Завершить урок' : 'Следующее слово' }}
          </button>
        </div>
      </div>
      
      <div v-else class="lesson-complete">
        <div class="complete-icon">🎉</div>
        <h2>Урок завершен!</h2>
        <p class="score">Вы набрали {{ score }}/{{ lesson.words.length * 10 }} очков</p>
        
        <div class="complete-actions">
          <button @click="completeLesson" class="btn btn-primary btn-complete">
            Добавить слова в профиль
          </button>
          <button @click="goBack" class="btn btn-secondary">
            Вернуться к урокам
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'

export default {
  name: 'LessonView',
  
  setup() {
    const route = useRoute()
    const router = useRouter()
    const lesson = ref(null)
    const loading = ref(false)
    const currentExercise = ref(0)
    const selectedOption = ref(null)
    const feedback = ref(null)
    const showNextButton = ref(false)
    const score = ref(0)
    
    const API_URL = 'http://localhost:8080/api/lessons'
    
    const getAuthHeaders = () => {
      const token = localStorage.getItem('token')
      return {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
      }
    }
    
    const currentWord = computed(() => {
      if (!lesson.value || currentExercise.value >= lesson.value.words.length) {
        return null
      }
      return lesson.value.words[currentExercise.value]
    })
    
    const shuffledOptions = computed(() => {
      if (!currentWord.value) return []
      
      // Создаем правильный вариант
      const correctOption = {
        text: currentWord.value.meaning,
        correct: true
      }
      
      // Создаем неправильные варианты из других слов урока
      const otherWords = lesson.value.words
        .filter((_, index) => index !== currentExercise.value)
        .slice(0, 3)
      
      const wrongOptions = otherWords.map(word => ({
        text: word.meaning,
        correct: false
      }))
      
      // Добавляем правильный вариант и перемешиваем
      const options = [correctOption, ...wrongOptions]
      return options.sort(() => Math.random() - 0.5)
    })
    
    const loadLesson = async () => {
      loading.value = true
      
      try {
        const lessonId = route.params.id
        const response = await fetch(`${API_URL}/${lessonId}/words`)
        
        if (!response.ok) {
          throw new Error('Ошибка загрузки урока')
        }
        
        lesson.value = await response.json()
      } catch (error) {
        console.error('Ошибка:', error)
        alert('Ошибка загрузки урока')
      } finally {
        loading.value = false
      }
    }
    
    const checkAnswer = (option) => {
      if (selectedOption.value !== null) return
      
      selectedOption.value = option
      
      if (option.correct) {
        feedback.value = {
          type: 'correct',
          message: '✅ Правильно!'
        }
        score.value += 10
      } else {
        feedback.value = {
          type: 'incorrect',
          message: `❌ Неправильно. Правильный ответ: ${currentWord.value.meaning}`
        }
      }
      
      showNextButton.value = true
    }
    
    const nextExercise = () => {
      selectedOption.value = null
      feedback.value = null
      showNextButton.value = false
      currentExercise.value++
    }
    
    const completeLesson = async () => {
      try {
        const lessonId = route.params.id
        const response = await fetch(`${API_URL}/${lessonId}/complete`, {
          method: 'POST',
          headers: getAuthHeaders(),
          body: JSON.stringify({ score: score.value })
        })
        
        if (!response.ok) {
          throw new Error('Ошибка завершения урока')
        }
        
        alert('Урок успешно завершен! Слова добавлены в ваш профиль.')
        router.push('/lessons')
      } catch (error) {
        console.error('Ошибка:', error)
        alert('Ошибка завершения урока: ' + error.message)
      }
    }
    
    const goBack = () => {
      router.push('/lessons')
    }
    
    onMounted(() => {
      loadLesson()
    })
    
    return {
      lesson,
      loading,
      currentExercise,
      currentWord,
      shuffledOptions,
      selectedOption,
      feedback,
      showNextButton,
      score,
      checkAnswer,
      nextExercise,
      completeLesson,
      goBack
    }
  }
}
</script>

<style scoped>
.lesson-container {
  max-width: 800px;
  margin: 0 auto;
  padding: 30px 20px;
}

.lesson-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
  padding-bottom: 20px;
  border-bottom: 2px solid #e0e0e0;
}

.lesson-header h1 {
  margin: 0;
  color: #333;
  font-size: 28px;
}

.lesson-info {
  display: flex;
  align-items: center;
  gap: 15px;
}

.lesson-level, .lesson-progress {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 6px 16px;
  border-radius: 20px;
  font-size: 14px;
  font-weight: 600;
}

.btn-back {
  background: #f5f5f5;
  color: #333;
  padding: 8px 16px;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 15px;
  font-weight: 600;
}

.btn-back:hover {
  background: #e0e0e0;
}

.loading {
  text-align: center;
  padding: 60px;
  font-size: 18px;
  color: #666;
}

.exercise-container {
  display: flex;
  justify-content: center;
}

.exercise-card {
  background: white;
  border-radius: 16px;
  padding: 40px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1);
  width: 100%;
  max-width: 600px;
}

.word-display {
  text-align: center;
  margin-bottom: 30px;
}

.word-display h2 {
  font-size: 48px;
  font-weight: bold;
  color: #333;
  margin: 0 0 10px 0;
}

.transcription {
  font-size: 18px;
  color: #999;
  font-family: monospace;
  margin: 0;
}

.exercise-type h3 {
  text-align: center;
  color: #666;
  margin-bottom: 25px;
  font-size: 20px;
}

.options-grid {
  display: grid;
  grid-template-columns: 1fr;
  gap: 12px;
  margin-bottom: 25px;
}

.option-btn {
  padding: 16px 20px;
  border: 2px solid #e0e0e0;
  border-radius: 12px;
  background: white;
  color: #333;
  font-size: 18px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
  text-align: left;
}

.option-btn:hover:not(.disabled) {
  border-color: #667eea;
  background: #f8f9ff;
}

.option-btn.correct {
  background: #d4edda;
  border-color: #c3e6cb;
  color: #155724;
}

.option-btn.incorrect {
  background: #f8d7da;
  border-color: #f5c6cb;
  color: #721c24;
}

.option-btn.disabled {
  cursor: not-allowed;
  opacity: 0.7;
}

.feedback {
  padding: 15px;
  border-radius: 10px;
  margin-bottom: 25px;
  text-align: center;
  font-size: 18px;
  font-weight: 600;
}

.feedback.correct {
  background: #d4edda;
  color: #155724;
  border: 1px solid #c3e6cb;
}

.feedback.incorrect {
  background: #f8d7da;
  color: #721c24;
  border: 1px solid #f5c6cb;
}

.btn-next {
  width: 100%;
  padding: 16px;
  font-size: 18px;
}

.lesson-complete {
  text-align: center;
  padding: 60px 20px;
}

.complete-icon {
  font-size: 80px;
  margin-bottom: 20px;
}

.lesson-complete h2 {
  font-size: 36px;
  color: #4CAF50;
  margin-bottom: 15px;
}

.score {
  font-size: 24px;
  color: #666;
  margin-bottom: 30px;
  font-weight: 600;
}

.complete-actions {
  display: flex;
  flex-direction: column;
  gap: 15px;
  max-width: 400px;
  margin: 0 auto;
}

.btn-complete {
  background: linear-gradient(135deg, #4CAF50 0%, #45a049 100%);
}

.btn-secondary {
  background: #f5f5f5;
  color: #333;
}
</style>
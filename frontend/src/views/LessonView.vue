<template>
  <div class="lesson-container">
    <div class="lesson-header">
      <button @click="goBack" class="btn btn-back">← Назад</button>
      <h1>{{ lesson?.lesson?.title }}</h1>
      <div class="lesson-info">
        <span class="lesson-level">Уровень {{ lesson?.lesson?.level }}</span>
      </div>
    </div>
    
    <div v-if="loading" class="loading">Загрузка урока...</div>
    
    <div v-else-if="lesson" class="lesson-content">
      <p class="lesson-description">{{ lesson.lesson.description }}</p>
      
      <div class="words-list">
        <h2>Слова урока ({{ lesson.words.length }})</h2>
        
        <div class="words-grid">
          <div v-for="word in lesson.words" :key="word.id" class="word-item">
            <div class="word-main">
              <span class="word-text">{{ word.word }}</span>
              <span v-if="word.transcription" class="word-transcription">{{ word.transcription }}</span>
            </div>
            <div class="word-meaning">{{ word.meaning }}</div>
            <div v-if="word.example" class="word-example">
              <span class="example-label">Пример:</span>
              {{ word.example }}
            </div>
          </div>
        </div>
      </div>
      
      <div class="lesson-actions">
        <button @click="completeLesson" class="btn btn-primary">Завершить урок</button>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';

export default {
  name: 'LessonView',
  
  setup() {
    const route = useRoute();
    const router = useRouter();
    const lesson = ref(null);
    const loading = ref(false);
    
    const API_URL = 'http://localhost:8080/api/lessons';
    
    const getAuthHeaders = () => {
      const token = localStorage.getItem('token');
      return {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
      };
    };
    
    const loadLesson = async () => {
      loading.value = true;
      
      try {
        const lessonId = route.params.id;
        const response = await fetch(`${API_URL}/${lessonId}/words`, {
          headers: getAuthHeaders()
        });
        
        if (!response.ok) {
          throw new Error('Ошибка загрузки урока');
        }
        
        lesson.value = await response.json();
      } catch (error) {
        console.error('Ошибка:', error);
        alert('Ошибка загрузки урока');
      } finally {
        loading.value = false;
      }
    };
    
    const completeLesson = async () => {
      if (!confirm('Вы уверены, что хотите завершить этот урок?')) {
        return;
      }
      
      try {
        const lessonId = route.params.id;
        const response = await fetch(`${API_URL}/${lessonId}/complete`, {
          method: 'POST',
          headers: getAuthHeaders(),
          body: JSON.stringify({ score: 100 })
        });
        
        if (!response.ok) {
          throw new Error('Ошибка завершения урока');
        }
        
        alert('Урок успешно завершён!');
        router.push('/lessons');
      } catch (error) {
        console.error('Ошибка:', error);
        alert('Ошибка завершения урока');
      }
    };
    
    const goBack = () => {
      router.push('/lessons');
    };
    
    onMounted(() => {
      loadLesson();
    });
    
    return {
      lesson,
      loading,
      completeLesson,
      goBack
    };
  }
};
</script>

<style scoped>
.lesson-container {
  max-width: 900px;
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
}

.lesson-info {
  display: flex;
  align-items: center;
  gap: 15px;
}

.lesson-level {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 6px 16px;
  border-radius: 20px;
  font-size: 14px;
}

.btn-back {
  background: #f5f5f5;
  color: #333;
  padding: 8px 16px;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 15px;
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

.lesson-description {
  background: #f9f9f9;
  padding: 20px;
  border-radius: 10px;
  margin-bottom: 30px;
  color: #555;
  line-height: 1.6;
}

.words-list h2 {
  margin-bottom: 20px;
  color: #333;
}

.words-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
  gap: 20px;
  margin-bottom: 30px;
}

.word-item {
  background: white;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
  border: 1px solid #e0e0e0;
}

.word-main {
  display: flex;
  flex-direction: column;
  gap: 5px;
  margin-bottom: 10px;
}

.word-text {
  font-size: 22px;
  font-weight: bold;
  color: #333;
}

.word-transcription {
  font-size: 14px;
  color: #999;
  font-family: monospace;
}

.word-meaning {
  color: #555;
  margin-bottom: 10px;
  font-style: italic;
}

.word-example {
  background: #f5f5f5;
  padding: 10px;
  border-radius: 6px;
  font-size: 14px;
  color: #666;
}

.example-label {
  font-weight: bold;
  color: #333;
}

.lesson-actions {
  display: flex;
  justify-content: center;
  gap: 15px;
}

.btn {
  padding: 14px 32px;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-size: 16px;
  font-weight: 600;
  transition: all 0.2s;
}

.btn-primary {
  background: linear-gradient(135deg, #4CAF50 0%, #45a049 100%);
  color: white;
}

.btn-primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 12px rgba(76, 175, 80, 0.4);
}
</style>
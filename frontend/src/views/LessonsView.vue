<template>
  <div class="lessons-container">
    <div class="header">
      <h1>📚 Уроки</h1>
      <div v-if="progressStats" class="stats">
        <div class="stat-card">
          <span class="stat-value">{{ progressStats.completed_lessons || 0 }}</span>
          <span class="stat-label">Завершено</span>
        </div>
        <div class="stat-card">
          <span class="stat-value">{{ progressStats.total_score || 0 }}</span>
          <span class="stat-label">Очков</span>
        </div>
      </div>
    </div>
    
    <div v-if="loading" class="loading">Загрузка уроков...</div>
    
    <div v-else-if="error" class="error-message">
      {{ error }}
    </div>
    
    <div v-else-if="lessons.length === 0" class="empty-state">
      <p>Нет доступных уроков</p>
    </div>
    
    <div v-else class="lessons-grid">
      <div 
        v-for="lesson in lessons" 
        :key="lesson.id" 
        class="lesson-card"
        :class="{ 
          'completed': userProgress[lesson.id]?.status === 'completed',
          'in-progress': userProgress[lesson.id]?.status === 'in_progress'
        }"
      >
        <div class="lesson-header">
          <div class="lesson-badge">
            <span v-if="userProgress[lesson.id]?.status === 'completed'">✓</span>
            <span v-else-if="userProgress[lesson.id]?.status === 'in_progress'">▶</span>
            <span v-else>{{ lesson.order_num }}</span>
          </div>
          <div class="lesson-level">
            <span>Уровень {{ lesson.level }}</span>
          </div>
        </div>
        
        <h3 class="lesson-title">{{ lesson.title }}</h3>
        <p class="lesson-description">{{ lesson.description }}</p>
        
        <div class="lesson-footer">
          <span class="lesson-type">{{ getLessonTypeLabel(lesson.lesson_type) }}</span>
          <button 
            @click="startLesson(lesson.id)" 
            class="btn btn-primary"
            :disabled="userProgress[lesson.id]?.status === 'completed'"
          >
            {{ getButtonText(lesson.id) }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';

export default {
  name: 'LessonsView',
  
  setup() {
    const router = useRouter();
    const lessons = ref([]);
    const userProgress = ref({});
    const progressStats = ref(null);
    const loading = ref(false);
    const error = ref(null);
    
    const API_URL = 'http://localhost:8080/api/lessons';
    
    const getAuthHeaders = () => {
      const token = localStorage.getItem('token');
      return {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
      };
    };
    
    const loadLessons = async () => {
      loading.value = true;
      error.value = null;
      
      try {
        // Загружаем уроки (публичный эндпоинт)
        const lessonsResponse = await fetch(`${API_URL}`);
        
        if (!lessonsResponse.ok) {
          throw new Error(`Ошибка загрузки уроков: ${lessonsResponse.status}`);
        }
        
        lessons.value = await lessonsResponse.json();
        
        // Загружаем прогресс пользователя
        await loadUserProgress();
        
        // Загружаем статистику
        await loadProgressStats();
      } catch (err) {
        console.error('Ошибка загрузки:', err);
        error.value = err.message || 'Ошибка загрузки уроков';
      } finally {
        loading.value = false;
      }
    };
    
    const loadUserProgress = async () => {
      try {
        const response = await fetch(`${API_URL}/my`, {
          headers: getAuthHeaders()
        });
        
        if (response.ok) {
          const progress = await response.json();
          
          const progressObj = {};
          progress.forEach(item => {
            progressObj[item.lesson_id] = item;
          });
          
          userProgress.value = progressObj;
        }
      } catch (err) {
        console.error('Ошибка загрузки прогресса:', err);
      }
    };
    
    const loadProgressStats = async () => {
      try {
        const response = await fetch(`${API_URL}/stats`, {
          headers: getAuthHeaders()
        });
        
        if (response.ok) {
          progressStats.value = await response.json();
        }
      } catch (err) {
        console.error('Ошибка загрузки статистики:', err);
      }
    };
    
    const startLesson = async (lessonId) => {
      const currentProgress = userProgress.value[lessonId];
      
      if (currentProgress?.status === 'completed') {
        router.push(`/lesson/${lessonId}`);
        return;
      }
      
      try {
        const response = await fetch(`${API_URL}/${lessonId}/start`, {
          method: 'POST',
          headers: getAuthHeaders()
        });
        
        if (!response.ok) {
          throw new Error(`Ошибка начала урока: ${response.status}`);
        }
        
        await loadUserProgress();
        router.push(`/lesson/${lessonId}`);
      } catch (err) {
        console.error('Ошибка начала урока:', err);
        alert(err.message || 'Ошибка начала урока');
      }
    };
    
    const getButtonText = (lessonId) => {
      const progress = userProgress.value[lessonId];
      
      if (progress?.status === 'completed') {
        return 'Повторить';
      }
      
      if (progress?.status === 'in_progress') {
        return 'Продолжить';
      }
      
      return 'Начать';
    };
    
    const getLessonTypeLabel = (type) => {
      const labels = {
        'vocabulary': 'Словарь',
        'grammar': 'Грамматика',
        'listening': 'Аудирование',
        'reading': 'Чтение',
        'writing': 'Письмо'
      };
      
      return labels[type] || type;
    };
    
    onMounted(() => {
      loadLessons();
    });
    
    return {
      lessons,
      userProgress,
      progressStats,
      loading,
      error,
      startLesson,
      getButtonText,
      getLessonTypeLabel
    };
  }
};
</script>


<style scoped>
.lessons-container {
  max-width: 1400px;
  margin: 0 auto;
  padding: 20px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
}

.header h1 {
  margin: 0;
  color: #333;
}

.stats {
  display: flex;
  gap: 20px;
}

.stat-card {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 15px 25px;
  border-radius: 10px;
  text-align: center;
}

.stat-value {
  font-size: 28px;
  font-weight: bold;
  display: block;
}

.stat-label {
  font-size: 14px;
  opacity: 0.9;
}

.loading {
  text-align: center;
  padding: 60px;
  font-size: 18px;
  color: #666;
}

.empty-state {
  text-align: center;
  padding: 60px;
  color: #999;
}

.lessons-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 25px;
}

.lesson-card {
  background: white;
  border-radius: 12px;
  padding: 25px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
  border: 2px solid #e0e0e0;
}

.lesson-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.15);
}

.lesson-card.completed {
  border-color: #4CAF50;
  background: linear-gradient(135deg, #f1f8e9 0%, #ffffff 100%);
}

.lesson-card.in-progress {
  border-color: #2196F3;
  background: linear-gradient(135deg, #e3f2fd 0%, #ffffff 100%);
}

.lesson-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
}

.lesson-badge {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
  font-weight: bold;
}

.lesson-card.completed .lesson-badge {
  background: linear-gradient(135deg, #4CAF50 0%, #45a049 100%);
}

.lesson-card.in-progress .lesson-badge {
  background: linear-gradient(135deg, #2196F3 0%, #1976D2 100%);
}

.lesson-level {
  background: #f5f5f5;
  padding: 5px 12px;
  border-radius: 20px;
  font-size: 13px;
  color: #666;
}

.lesson-title {
  margin: 0 0 10px 0;
  color: #333;
  font-size: 20px;
}

.lesson-description {
  margin: 0 0 20px 0;
  color: #666;
  line-height: 1.6;
}

.lesson-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.lesson-type {
  font-size: 13px;
  color: #999;
  font-style: italic;
}

.btn {
  padding: 10px 20px;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 15px;
  font-weight: 600;
  transition: all 0.2s;
}

.btn-primary {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.btn-primary:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(102, 126, 234, 0.4);
}

.btn-primary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}
</style>
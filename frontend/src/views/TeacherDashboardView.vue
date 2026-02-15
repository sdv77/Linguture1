<template>
  <div class="teacher-dashboard">
    <div class="dashboard-header">
      <h1>👨‍🏫 Панель учителя</h1>
      <div class="header-actions">
        <button @click="showCreateModal = true" class="btn btn-primary">
          ➕ Создать урок
        </button>
        <button @click="logout" class="btn btn-secondary">
          🚪 Выйти
        </button>
      </div>
    </div>
    
    <div v-if="loading" class="loading">Загрузка уроков...</div>
    
    <div v-else-if="error" class="error-message">
      {{ error }}
    </div>
    
    <div v-else-if="lessons.length === 0" class="empty-state">
      <p>У вас пока нет уроков</p>
      <button @click="showCreateModal = true" class="btn btn-primary">
        Создать первый урок
      </button>
    </div>
    
    <div v-else class="lessons-list">
      <div class="lessons-header">
        <h2>Ваши уроки ({{ lessons.length }})</h2>
      </div>
      
      <div class="lessons-grid">
        <div v-for="lesson in lessons" :key="lesson.id" class="lesson-card">
          <div class="lesson-card-header">
            <h3>{{ lesson.title }}</h3>
            <div class="lesson-card-actions">
              <button @click="editLesson(lesson)" class="btn btn-small btn-edit">✏️</button>
              <button @click="deleteLesson(lesson.id)" class="btn btn-small btn-delete">🗑️</button>
            </div>
          </div>
          
          <p class="lesson-description">{{ lesson.description }}</p>
          
          <div class="lesson-meta">
            <span class="lesson-level">Уровень {{ lesson.level }}</span>
            <span class="lesson-type">{{ getLessonTypeLabel(lesson.lesson_type) }}</span>
            <span class="lesson-order">Порядок: {{ lesson.order_num }}</span>
          </div>
        </div>
      </div>
    </div>
    
    <!-- Модальное окно -->
    <div v-if="showCreateModal" class="modal-overlay" @click="closeModal">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>{{ editingLesson ? 'Редактировать урок' : 'Создать урок' }}</h3>
          <button @click="closeModal" class="btn-close">×</button>
        </div>
        
        <form @submit.prevent="saveLesson" class="lesson-form">
          <div class="form-group">
            <label>Название урока *</label>
            <input v-model="lessonForm.title" required placeholder="Например: Основные приветствия" />
          </div>
          
          <div class="form-group">
            <label>Описание</label>
            <textarea v-model="lessonForm.description" rows="3" placeholder="Краткое описание урока"></textarea>
          </div>
          
          <div class="form-row">
            <div class="form-group">
              <label>Уровень сложности</label>
              <input type="number" v-model.number="lessonForm.level" min="1" max="10" />
            </div>
            
            <div class="form-group">
              <label>Порядок</label>
              <input type="number" v-model.number="lessonForm.order_num" min="0" />
            </div>
          </div>
          
          <div class="form-group">
            <label>Тип урока</label>
            <select v-model="lessonForm.lesson_type">
              <option value="vocabulary">Словарь</option>
              <option value="grammar">Грамматика</option>
              <option value="listening">Аудирование</option>
              <option value="reading">Чтение</option>
              <option value="writing">Письмо</option>
            </select>
          </div>
          
          <div class="modal-actions">
            <button type="submit" class="btn btn-primary">
              {{ editingLesson ? 'Сохранить' : 'Создать' }}
            </button>
            <button type="button" @click="closeModal" class="btn btn-secondary">
              Отмена
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';

export default {
  name: 'TeacherDashboardView',
  
  setup() {
    const router = useRouter();
    const lessons = ref([]);
    const loading = ref(false);
    const error = ref(null);
    const showCreateModal = ref(false);
    const editingLesson = ref(null);
    
    const lessonForm = ref({
      title: '',
      description: '',
      level: 1,
      lesson_type: 'vocabulary',
      order_num: 0
    });
    
    const API_URL = 'http://localhost:8080/api/teacher/lessons';
    
    const getAuthHeaders = () => {
      const token = localStorage.getItem('teacher_token');
      return {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
      };
    };
    
    const loadLessons = async () => {
      loading.value = true;
      error.value = null;
      
      try {
        const response = await fetch(API_URL, {
          headers: getAuthHeaders()
        });
        
        if (!response.ok) {
          throw new Error(`Ошибка загрузки уроков: ${response.status}`);
        }
        
        lessons.value = await response.json();
      } catch (err) {
        console.error('Ошибка:', err);
        error.value = err.message || 'Ошибка загрузки уроков';
      } finally {
        loading.value = false;
      }
    };
    
    const saveLesson = async () => {
      try {
        const url = editingLesson.value 
          ? `${API_URL}/${editingLesson.value.id}`
          : API_URL;
        
        const method = editingLesson.value ? 'PUT' : 'POST';
        
        const response = await fetch(url, {
          method: method,
          headers: getAuthHeaders(),
          body: JSON.stringify(lessonForm.value)
        });
        
        if (!response.ok) {
          const errorData = await response.json().catch(() => ({}));
          throw new Error(errorData.message || `Ошибка сохранения: ${response.status}`);
        }
        
        await loadLessons();
        closeModal();
        alert(editingLesson.value ? 'Урок успешно обновлен!' : 'Урок успешно создан!');
      } catch (err) {
        console.error('Ошибка:', err);
        alert('Ошибка сохранения урока: ' + err.message);
      }
    };
    
    const editLesson = (lesson) => {
      editingLesson.value = lesson;
      lessonForm.value = {
        title: lesson.title,
        description: lesson.description,
        level: lesson.level,
        lesson_type: lesson.lesson_type,
        order_num: lesson.order_num
      };
      showCreateModal.value = true;
    };
    
    const deleteLesson = async (id) => {
      if (!confirm('Вы уверены, что хотите удалить этот урок?')) {
        return;
      }
      
      try {
        const response = await fetch(`${API_URL}/${id}`, {
          method: 'DELETE',
          headers: getAuthHeaders()
        });
        
        if (!response.ok) {
          throw new Error(`Ошибка удаления: ${response.status}`);
        }
        
        await loadLessons();
        alert('Урок успешно удален!');
      } catch (err) {
        console.error('Ошибка:', err);
        alert('Ошибка удаления урока: ' + err.message);
      }
    };
    
    const closeModal = () => {
      showCreateModal.value = false;
      editingLesson.value = null;
      lessonForm.value = {
        title: '',
        description: '',
        level: 1,
        lesson_type: 'vocabulary',
        order_num: 0
      };
    };
    
    const logout = () => {
      localStorage.removeItem('teacher_token');
      router.push('/teacher/login');
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
      const token = localStorage.getItem('teacher_token');
      if (!token) {
        router.push('/teacher/login');
        return;
      }
      
      loadLessons();
    });
    
    return {
      lessons,
      loading,
      error,
      showCreateModal,
      editingLesson,
      lessonForm,
      saveLesson,
      editLesson,
      deleteLesson,
      closeModal,
      logout,
      getLessonTypeLabel
    };
  }
};
</script>


<style scoped>
.teacher-dashboard {
  max-width: 1400px;
  margin: 0 auto;
  padding: 30px 20px;
}

.dashboard-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
  padding-bottom: 20px;
  border-bottom: 2px solid #e0e0e0;
}

.dashboard-header h1 {
  margin: 0;
  color: #333;
  font-size: 32px;
}

.header-actions {
  display: flex;
  gap: 10px;
}

.loading {
  text-align: center;
  padding: 60px;
  font-size: 18px;
  color: #666;
}

.empty-state {
  text-align: center;
  padding: 80px 20px;
  color: #999;
}

.empty-state p {
  margin-bottom: 20px;
  font-size: 18px;
}

.lessons-list {
  background: white;
  border-radius: 12px;
  padding: 30px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

.lessons-header {
  margin-bottom: 25px;
}

.lessons-header h2 {
  margin: 0;
  color: #333;
}

.lessons-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 25px;
}

.lesson-card {
  background: #f9f9f9;
  border-radius: 10px;
  padding: 25px;
  border: 2px solid #e0e0e0;
  transition: all 0.3s;
}

.lesson-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.1);
  border-color: #667eea;
}

.lesson-card-header {
  display: flex;
  justify-content: space-between;
  align-items: start;
  margin-bottom: 15px;
}

.lesson-card-header h3 {
  margin: 0;
  color: #333;
  font-size: 20px;
}

.lesson-card-actions {
  display: flex;
  gap: 8px;
}

.lesson-description {
  margin: 0 0 20px 0;
  color: #666;
  line-height: 1.6;
}

.lesson-meta {
  display: flex;
  gap: 15px;
  padding-top: 15px;
  border-top: 1px solid #e0e0e0;
  color: #999;
  font-size: 14px;
}

.lesson-level, .lesson-type, .lesson-order {
  display: flex;
  align-items: center;
  gap: 5px;
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  background: white;
  border-radius: 12px;
  max-width: 600px;
  width: 90%;
  max-height: 90vh;
  overflow-y: auto;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  border-bottom: 2px solid #e0e0e0;
}

.modal-header h3 {
  margin: 0;
  color: #333;
}

.btn-close {
  background: none;
  border: none;
  font-size: 32px;
  cursor: pointer;
  color: #999;
  padding: 0;
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.btn-close:hover {
  color: #333;
  background: #f5f5f5;
  border-radius: 50%;
}

.lesson-form {
  padding: 20px;
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

.form-group input,
.form-group textarea,
.form-group select {
  width: 100%;
  padding: 12px;
  border: 2px solid #e0e0e0;
  border-radius: 8px;
  font-size: 16px;
  font-family: inherit;
}

.form-group input:focus,
.form-group textarea:focus,
.form-group select:focus {
  outline: none;
  border-color: #667eea;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 15px;
}

.modal-actions {
  display: flex;
  gap: 10px;
  margin-top: 20px;
  justify-content: flex-end;
}

.btn {
  padding: 12px 24px;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-size: 16px;
  font-weight: 600;
  transition: all 0.2s;
}

.btn-primary {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.btn-primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 12px rgba(102, 126, 234, 0.4);
}

.btn-secondary {
  background: #f5f5f5;
  color: #333;
}

.btn-secondary:hover {
  background: #e0e0e0;
}

.btn-small {
  padding: 6px 12px;
  font-size: 14px;
}

.btn-edit {
  background: #2196F3;
  color: white;
}

.btn-edit:hover {
  background: #0b7dda;
}

.btn-delete {
  background: #f44336;
  color: white;
}

.btn-delete:hover {
  background: #da190b;
}
</style>
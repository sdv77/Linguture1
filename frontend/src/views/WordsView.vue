<template>
  <div class="words-container">
    <h1>Словарь слов</h1>
    
    <!-- Форма добавления слова -->
    <div class="form-container">
      <h2>{{ editingWord ? 'Редактировать слово' : 'Добавить новое слово' }}</h2>
      <form @submit.prevent="saveWord">
        <div class="form-group">
          <label for="word">Слово:</label>
          <input 
            type="text" 
            id="word" 
            v-model="formData.word" 
            required 
            placeholder="Введите слово"
          />
        </div>
        
        <div class="form-group">
          <label for="meaning">Значение:</label>
          <textarea 
            id="meaning" 
            v-model="formData.meaning" 
            required 
            placeholder="Введите значение слова"
            rows="3"
          ></textarea>
        </div>
        
        <div class="form-actions">
          <button type="submit" class="btn btn-primary">
            {{ editingWord ? 'Сохранить' : 'Добавить' }}
          </button>
          <button 
            type="button" 
            v-if="editingWord" 
            @click="cancelEdit"
            class="btn btn-secondary"
          >
            Отмена
          </button>
        </div>
      </form>
    </div>
    
    <!-- Сообщение об ошибке -->
    <div v-if="errorMessage" class="error-message">
      {{ errorMessage }}
    </div>
    
    <!-- Список слов -->
    <div class="words-list">
      <h2>Все слова ({{ words.length }})</h2>
      
      <div v-if="loading" class="loading">Загрузка...</div>
      
      <div v-else-if="words.length === 0" class="empty-state">
        <p>Нет слов в словаре. Добавьте первое слово!</p>
      </div>
      
      <div v-else class="words-grid">
        <div 
          v-for="word in words" 
          :key="word.id" 
          class="word-card"
        >
          <div class="word-header">
            <h3>{{ word.word }}</h3>
            <div class="word-actions">
              <button 
                @click="editWord(word)" 
                class="btn btn-small btn-edit"
                title="Редактировать"
              >
                ✏️
              </button>
              <button 
                @click="deleteWord(word.id)" 
                class="btn btn-small btn-delete"
                title="Удалить"
              >
                🗑️
              </button>
            </div>
          </div>
          <p class="word-meaning">{{ word.meaning }}</p>
          <div class="word-meta">
            <small>Создано: {{ formatDate(word.created_at) }}</small>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue';

export default {
  name: 'WordsView',
  
  setup() {
    // Состояние компонента
    const words = ref([]);
    const loading = ref(false);
    const editingWord = ref(null);
    const formData = ref({
      word: '',
      meaning: ''
    });
    const errorMessage = ref('');
    
    // Базовый URL API
    const API_URL = 'http://localhost:8080/api/words';
    
    // Загрузка всех слов
    const loadWords = async () => {
      loading.value = true;
      errorMessage.value = '';
      
      try {
        const response = await fetch(API_URL);
        
        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }
        
        const data = await response.json();
        words.value = data;
      } catch (error) {
        console.error('Ошибка загрузки слов:', error);
        errorMessage.value = 'Ошибка при загрузке слов: ' + error.message;
      } finally {
        loading.value = false;
      }
    };
    
    // Сохранение слова (создание или обновление)
    const saveWord = async () => {
      try {
        const url = editingWord.value 
          ? `${API_URL}/${editingWord.value.id}` 
          : API_URL;
        
        const method = editingWord.value ? 'PUT' : 'POST';
        
        console.log('Отправка запроса:', { url, method, data: formData.value });
        
        const response = await fetch(url, {
          method: method,
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify(formData.value)
        });
        
        console.log('Ответ сервера:', {
          status: response.status,
          ok: response.ok,
          headers: Object.fromEntries(response.headers.entries())
        });
        
        if (!response.ok) {
          const errorText = await response.text();
          console.error('Тело ошибки:', errorText);
          throw new Error(`Ошибка ${method}: ${response.status} - ${errorText || response.statusText}`);
        }
        
        // Обновляем список
        await loadWords();
        resetForm();
        alert(editingWord.value ? 'Слово успешно обновлено!' : 'Слово успешно добавлено!');
      } catch (error) {
        console.error('Ошибка сохранения слова:', error);
        alert('Ошибка при сохранении слова: ' + error.message);
      }
    };
    
    // Редактирование слова
    const editWord = (word) => {
      editingWord.value = word;
      formData.value = {
        word: word.word,
        meaning: word.meaning
      };
    };
    
    // Отмена редактирования
    const cancelEdit = () => {
      editingWord.value = null;
      resetForm();
    };
    
    // Удаление слова
    const deleteWord = async (id) => {
      if (!confirm('Вы уверены, что хотите удалить это слово?')) {
        return;
      }
      
      try {
        const url = `${API_URL}/${id}`;
        console.log('Отправка запроса на удаление:', url);
        
        const response = await fetch(url, {
          method: 'DELETE'
        });
        
        console.log('Ответ сервера:', {
          status: response.status,
          ok: response.ok
        });
        
        if (!response.ok) {
          const errorText = await response.text();
          console.error('Тело ошибки:', errorText);
          throw new Error(`Ошибка удаления: ${response.status} - ${errorText || response.statusText}`);
        }
        
        // Обновляем список
        await loadWords();
        alert('Слово успешно удалено!');
      } catch (error) {
        console.error('Ошибка удаления слова:', error);
        alert('Ошибка при удалении слова: ' + error.message);
      }
    };
    
    // Сброс формы
    const resetForm = () => {
      formData.value = {
        word: '',
        meaning: ''
      };
      editingWord.value = null;
    };
    
    // Форматирование даты
    const formatDate = (dateString) => {
      const date = new Date(dateString);
      return date.toLocaleDateString('ru-RU', {
        year: 'numeric',
        month: 'long',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
      });
    };
    
    // Загружаем слова при монтировании компонента
    onMounted(() => {
      loadWords();
    });
    
    return {
      words,
      loading,
      editingWord,
      formData,
      errorMessage,
      saveWord,
      editWord,
      cancelEdit,
      deleteWord,
      formatDate
    };
  }
};
</script>

<style scoped>
.words-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.form-container {
  background: #f5f5f5;
  padding: 20px;
  border-radius: 8px;
  margin-bottom: 30px;
}

.form-group {
  margin-bottom: 15px;
}

.form-group label {
  display: block;
  margin-bottom: 5px;
  font-weight: bold;
}

.form-group input,
.form-group textarea {
  width: 100%;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 16px;
}

.form-actions {
  display: flex;
  gap: 10px;
  margin-top: 20px;
}

.error-message {
  background-color: #ffebee;
  color: #c62828;
  padding: 15px;
  border-radius: 4px;
  margin-bottom: 20px;
  border-left: 4px solid #c62828;
}

.words-list {
  background: #fff;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.loading {
  text-align: center;
  padding: 40px;
  font-size: 18px;
  color: #666;
}

.empty-state {
  text-align: center;
  padding: 40px;
  color: #999;
}

.words-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
  margin-top: 20px;
}

.word-card {
  background: #f9f9f9;
  padding: 20px;
  border-radius: 8px;
  border: 1px solid #ddd;
  transition: transform 0.2s;
}

.word-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0,0,0,0.1);
}

.word-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.word-header h3 {
  margin: 0;
  color: #333;
}

.word-actions {
  display: flex;
  gap: 5px;
}

.word-meaning {
  margin: 10px 0;
  color: #555;
  line-height: 1.6;
}

.word-meta {
  margin-top: 10px;
  padding-top: 10px;
  border-top: 1px solid #eee;
  color: #999;
  font-size: 14px;
}

.btn {
  padding: 10px 20px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 16px;
  transition: background-color 0.2s;
}

.btn-primary {
  background-color: #4CAF50;
  color: white;
}

.btn-primary:hover {
  background-color: #45a049;
}

.btn-secondary {
  background-color: #6c757d;
  color: white;
}

.btn-secondary:hover {
  background-color: #5a6268;
}

.btn-small {
  padding: 5px 10px;
  font-size: 14px;
}

.btn-edit {
  background-color: #2196F3;
  color: white;
}

.btn-edit:hover {
  background-color: #0b7dda;
}

.btn-delete {
  background-color: #f44336;
  color: white;
}

.btn-delete:hover {
  background-color: #da190b;
}
</style>
<template>
  <div class="home-container">
    <header class="home-header">
      <div class="header-content">
        <h1 class="app-title">Linguture</h1>
        <p class="app-subtitle">Изучайте языки легко и эффективно</p>
        <div class="header-buttons">
          <router-link to="/register" class="btn btn-primary btn-large">Присоединиться</router-link>
          <router-link to="/login" class="btn btn-secondary btn-large">Войти</router-link>
        </div>
      </div>
    </header>

    <section class="features-section">
      <div class="container">
        <h2 class="section-title">Почему выбирают нас?</h2>
        <div class="features-grid">
          <div class="feature-card">
            <div class="feature-icon">🎓</div>
            <h3>Интерактивные уроки</h3>
            <p>Учитесь через увлекательные упражнения и игры, которые делают процесс изучения языка интересным.</p>
          </div>
          <div class="feature-card">
            <div class="feature-icon">📊</div>
            <h3>Отслеживание прогресса</h3>
            <p>Следите за своим прогрессом, получайте достижения и мотивируйтесь достигать новых высот.</p>
          </div>
          <div class="feature-card">
            <div class="feature-icon">📚</div>
            <h3>Персональный словарь</h3>
            <p>Все выученные слова сохраняются в вашем личном словаре для дальнейшего повторения.</p>
          </div>
          <div class="feature-card">
            <div class="feature-icon">🎯</div>
            <h3>Адаптивное обучение</h3>
            <p>Система подстраивается под ваш уровень и темп обучения для максимальной эффективности.</p>
          </div>
        </div>
      </div>
    </section>

    <section class="cta-section">
      <div class="container">
        <h2 class="cta-title">Готовы начать свое языковое путешествие?</h2>
        <p class="cta-subtitle">Присоединяйтесь к тысячам студентов, которые уже учатся с нами</p>
        <div class="cta-buttons">
          <router-link to="/register" class="btn btn-primary btn-xlarge">Начать бесплатно</router-link>
        </div>
      </div>
    </section>

    <footer class="home-footer">
      <div class="container">
        <p>&copy; 2026 Linguture. Все права защищены.</p>
        <div class="footer-links">
          <!-- Скрытая ссылка для учителей -->
          <router-link to="/teacher/login" class="footer-link teacher-link">Для учителей</router-link>
        </div>
      </div>
    </footer>
  </div>
</template>

<script>
import { useRouter } from 'vue-router'

export default {
  name: 'HomeView',
  setup() {
    const router = useRouter()

    // Если пользователь авторизован, перенаправляем на уроки
    if (localStorage.getItem('token')) {
      router.push('/lessons')
    }

    return {}
  }
}
</script>

<style scoped>
.home-container {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.home-header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 100px 20px;
  text-align: center;
}

.header-content {
  max-width: 800px;
  margin: 0 auto;
}

.app-title {
  font-size: 48px;
  margin-bottom: 20px;
  font-weight: bold;
}

.app-subtitle {
  font-size: 24px;
  margin-bottom: 40px;
  opacity: 0.95;
}

.header-buttons {
  display: flex;
  gap: 20px;
  justify-content: center;
  flex-wrap: wrap;
}

.btn-large {
  padding: 16px 40px;
  font-size: 18px;
  border-radius: 50px;
  font-weight: 600;
}

.features-section {
  padding: 80px 20px;
  background: #f8f9fa;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
}

.section-title {
  text-align: center;
  font-size: 36px;
  margin-bottom: 60px;
  color: #333;
  font-weight: bold;
}

.features-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 30px;
}

.feature-card {
  background: white;
  padding: 35px 30px;
  border-radius: 16px;
  box-shadow: 0 4px 15px rgba(0,0,0,0.08);
  text-align: center;
  transition: all 0.3s ease;
}

.feature-card:hover {
  transform: translateY(-8px);
  box-shadow: 0 8px 25px rgba(0,0,0,0.12);
}

.feature-icon {
  font-size: 56px;
  margin-bottom: 20px;
}

.feature-card h3 {
  margin-bottom: 15px;
  color: #333;
  font-size: 22px;
}

.feature-card p {
  color: #666;
  line-height: 1.6;
  font-size: 16px;
}

.cta-section {
  padding: 100px 20px;
  background: linear-gradient(135deg, #4CAF50 0%, #45a049 100%);
  color: white;
  text-align: center;
}

.cta-title {
  font-size: 42px;
  margin-bottom: 20px;
  font-weight: bold;
}

.cta-subtitle {
  font-size: 22px;
  margin-bottom: 40px;
  opacity: 0.95;
}

.cta-buttons {
  display: flex;
  justify-content: center;
}

.btn-xlarge {
  padding: 20px 60px;
  font-size: 22px;
  border-radius: 50px;
  font-weight: 700;
  box-shadow: 0 6px 20px rgba(0,0,0,0.2);
}

.home-footer {
  background: #2c3e50;
  color: rgba(255,255,255,0.7);
  padding: 25px;
  text-align: center;
  margin-top: auto;
}

.footer-links {
  margin-top: 15px;
}

.footer-link {
  color: rgba(255,255,255,0.5);
  text-decoration: none;
  margin: 0 10px;
  font-size: 14px;
  transition: color 0.3s;
}

.footer-link:hover {
  color: white;
}

/* Скрытая ссылка для учителей */
.teacher-link {
  opacity: 0.3;
  font-size: 12px;
  transition: opacity 0.3s;
}

.teacher-link:hover {
  opacity: 1;
}
</style>
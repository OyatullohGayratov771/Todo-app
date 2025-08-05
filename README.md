# 📝 Todo-App (Microservices Architecture)

Bu loyiha **mikroservis arxitektura** asosida yaratilgan **Todo List Application** bo'lib, foydalanuvchilar uchun vazifalarni boshqarish imkoniyatini taqdim etadi. Loyihaning maqsadi – **kengaytiriladigan, barqaror va real-time** xususiyatlarga ega task management platformasini yaratishdir.

---

## 🚀 Xususiyatlar

- ✅ Vazifalarni qo'shish, o'chirish, tahrirlash
- 👤 Foydalanuvchi ro'yxatdan o'tkazish va autentifikatsiya
- 🔔 Kafka orqali real-time bildirishnomalar
- ⚡ Redis bilan tezkor keshlash
- 🌐 API Gateway orqali mikroservislarni birlashtirish
- 🐳 Docker Compose bilan oson deploy

---

## 🛠 Texnologiyalar

- **Backend:** Golang (Gin, GORM)
- **Frontend:** React (Vite yoki Create React App)
- **Database:** PostgreSQL
- **Cache:** Redis
- **Messaging Queue:** Kafka + Zookeeper
- **Gateway:** Go-based API Gateway (yoki Nginx)
- **Containerization:** Docker, Docker Compose

---

## 📦 Mikroservislar

1. **User Service**  
   - Foydalanuvchilarni ro'yxatdan o'tkazish va autentifikatsiya qilish
   - JWT asosida avtorizatsiya
   - PostgreSQL bilan ishlash

2. **Task Service**  
   - Vazifalarni yaratish, o'chirish, yangilash
   - Redis orqali keshlash
   - PostgreSQL bilan integratsiya

3. **Notification Service**  
   - Kafka orqali task yangilanishlaridan keyin xabarnoma yuborish
   - E-mail yoki WebSocket orqali foydalanuvchiga yetkazish

4. **API Gateway**  
   - Barcha servislarni bitta kirish nuqtasida birlashtirish
   - Request marshrutizatsiyasi va load balancing

5. **Frontend (React)**  
   - UI interfeys orqali foydalanuvchi bilan o‘zaro aloqa
   - REST API orqali gateway bilan ishlash

---

## 📂 Loyihaning tuzilishi
```
todo-app/
│── api-gateway/
│
│── user-service/
│
│── task-service/
│
│── notification-service/
│
│── frontend/
│ ├── src/
│ ├── public/
│ └── package.json
│
│── docker-compose.yml
│── README.md

```


---

## ⚙️ O‘rnatish va ishga tushirish

### 1️⃣ Repozitoriyani klonlash
```bash
git clone https://github.com/OyatullohGayratov771/Todo-app.git
cd todo-app
```

### 2️⃣ Muhit o'zgaruvchilarini sozlash
#### Har bir servis uchun ***.env*** faylini yarating.

### Misol (user-service/.env):
```
HTTP_HOST=user-service
HTTP_PORT=8081

DB_HOST=db
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=secret
DB_NAME=todo_app

REDIS_HOST=redis
REDIS_PORT=6379

KAFKA_HOST=kafka
KAFKA_PORT=9092
KAFKA_GROUP=notifaction-group
KAFKA_TOPIC=notifications

JWT_KEY=secret
```

### Misol (task-service/.env):
```env
HTTP_HOST=task-service
HTTP_PORT=8082

DB_HOST=db
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=secret
DB_NAME=todo_app

REDIS_HOST=redis
REDIS_PORT=6379
```

### Misol (notification-service/.env):
```env
HTTP_HOST=notification-service
HTTP_PORT=8083

KAFKA_HOST=kafka
KAFKA_PORT=9092
KAFKA_GROUP=notifaction-group
KAFKA_TOPIC=notifications
```
### Misol (api-gateway/.env):
```env
HTTP_HOST=api-gateway
HTTP_PORT=8080

USER_SERVICE_HOST=user-service
USER_SERVICE_PORT=8081

TASK_SERVICE_HOST=task-service
TASK_SERVICE_PORT=8082

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=secret
DB_NAME=todo_db
```

### 3️⃣ Docker Compose orqali barcha xizmatlarni ishga tushirish
```bash
docker-compose up --build
```


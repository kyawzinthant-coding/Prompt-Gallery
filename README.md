# **📚 PromptGallery API Documentation**
## **🎯 What is PromptGallery?**
**PromptGallery** is a **coding prompt management platform** that serves as a centralized repository for programming challenges, exercises, and coding prompts. Think of it as a **GitHub for coding problems** or a **LeetCode-style platform** where developers can:
- 📝 **Create and share coding prompts** across different programming languages
- 🔍 **Discover programming challenges** filtered by difficulty, language, and category
- 🏆 **Build a collection** of coding exercises for practice and interviews
- 🎓 **Learn and practice** with curated programming problems
- 📊 **Track popularity** and engagement of different coding challenges



### **🎯 Target Audience:**
- Software developers
- Coding bootcamp instructors and students
- Computer science educators
- Programming enthusiasts


## **🛣️ API Routes**
### **🏥 System Health**

| Method | Endpoint | Description |
| --- | --- | --- |
| `GET` | `/health` | Check if the API server is running |
### **📝 Prompt Management**

| Method | Endpoint | Description |
| --- | --- | --- |
| `GET` | `/api/v1/prompts` | Retrieve all coding prompts with filtering and pagination |
| `POST` | `/api/v1/prompts` | Create a new coding prompt |
| `GET` | `/api/v1/prompts/:id` | Get a specific prompt by ID |
| `DELETE` | `/api/v1/prompts/:id` |


## **🏗️ API Architecture**
**Base URL**: `http://localhost:8080`
**API Version**: `v1`
**Data Format**: JSON
**Framework**: Go Fiber
**Database**: PostgreSQL with GORM
**Route Patterns**:
- Health: `/health`
- API Base: `/api/v1`
- Resource: `/api/v1/prompts`
- Resource Item: `/api/v1/prompts/:id`
- Catch All: `*` (404 handler)

This API provides a **solid foundation** for building a coding prompt platform, supporting the core functionality needed for **creating, discovering, and managing programming challenges**. 🚀

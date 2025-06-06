# 🎓 Learning Lab

A personal lab for **concept-driven learning** — start with ideas, collaborate with AI to complete them, challenge assumptions, and verify knowledge through building.

This repo is for:
- **Starting with concepts and ideas** — not predetermined tutorials
- **Collaborating with AI** to turn ideas into working implementations  
- **Challenging and verifying** concepts through hands-on building
- **Learning by doing** — where curiosity drives the direction
- **Documenting the journey** from idea to implementation to validation

---

## 🧠 Learning Philosophy

> *"Start with curiosity, build with AI, verify through experience."*

### 🔄 The Learning Cycle

1. **💡 Concept/Idea** — You bring the curiosity and questions
2. **🤖 AI Collaboration** — AI helps design, implement, and explore
3. **⚡ Build & Test** — Create working implementations
4. **🔍 Challenge & Verify** — Test assumptions, explore edge cases
5. **📚 Document & Reflect** — Capture learnings and insights

### 🎯 Core Principles

**🚀 Idea-First Learning**
- Learning starts with **your curiosity** and concepts you want to explore
- No predefined curriculum — follow your interests and questions
- AI becomes your **thinking partner** and implementation assistant

**🔬 Hypothesis-Driven**
- Turn ideas into **testable hypotheses** 
- Build minimal implementations to **verify or challenge** concepts
- Learn through **experimentation** rather than memorization

**🛠️ Build to Understand**
- **Every concept becomes code** — no pure theory
- **Working examples** validate understanding
- **Breaking things** teaches as much as building them

**🤝 Human-AI Collaboration**
- You provide **direction, creativity, and critical thinking**
- AI provides **implementation speed, knowledge synthesis, and exploration**
- Together you **explore deeper** than either could alone

---

## 🧱 Structure

| Folder | Description |
|--------|-------------|
| `backend/` | **Backend concept explorations** — APIs, databases, caching, gRPC, streaming |
| `frontend/` | **Frontend concept explorations** *(future)* — React, state management, performance |
| `devops/` | **Infrastructure concept explorations** *(future)* — Docker, CI/CD, monitoring |

---

## 🛠️ Getting Started

### Start with Backend Concepts
```bash
git clone https://github.com/e6a5/learning.git
cd learning
```

**Pick a concept that makes you curious:**

| Question | Exploration | Status |
|----------|-------------|---------|
| "How do REST APIs work in practice?" | `backend/01-http-server/` | ✅ **Ready** |
| "How do apps talk to databases?" | `backend/02-mysql-crud/` | ✅ **Ready** |
| "How does caching improve performance?" | `backend/03-redis-intro/` | ✅ **Ready** |
| "What makes gRPC different from HTTP?" | `backend/04-grpc-basics/` | ✅ **Ready** |

### Collaborate with AI
1. **Pick a concept** that interests you from above
2. **Explore the existing implementation** — see how the question was answered
3. **Ask "What if...?" questions** — extend or modify with AI help
4. **Test your assumptions** — measure, break, and rebuild

### Example Exploration Flow
```bash
# Start with caching concepts
cd backend/03-redis-intro

# Run the existing implementation
go run main.go

# Then collaborate with AI:
# "What if I test this with 10,000 concurrent requests?"
# "How would this perform with PostgreSQL instead of MySQL?"
# "Can I add cache invalidation strategies?"
# "What happens if Redis goes down?"
```

---

## 📁 Project Organization

Each concept exploration follows this structure:
```
backend/concept-name/
├── README.md           # The question, hypothesis, and findings
├── main.go            # Core implementation
├── Makefile           # Build and test commands  
├── compose.yml        # Infrastructure setup (Docker)
├── go.mod/go.sum      # Dependencies
└── [variations/]      # Alternative approaches (when exploring)
```

---

## 🎯 Learning Goals

**Technical Skills**
- Learn technologies **in context** of real problems
- Understand **trade-offs** through hands-on experience
- Build **intuition** through experimentation

**Collaboration Skills**  
- Communicate ideas effectively to AI assistants
- Direct AI implementation while maintaining **creative control**
- **Synthesize AI suggestions** with your own critical thinking

**Problem-Solving Skills**
- Break down **abstract concepts** into testable implementations
- **Challenge assumptions** through measurement and testing
- Learn from **failure and unexpected results**

---

## 🚀 Next Steps

**Immediate Explorations:**
- Extend existing backend concepts with AI collaboration
- Performance test current implementations  
- Add monitoring and observability
- Explore failure scenarios and edge cases

**Future Concept Areas:**
- **Frontend:** "How do modern UIs handle complex state?"
- **DevOps:** "How do you deploy systems reliably?"
- **Architecture:** "How do you design for scale?"

---

## 📩 Contributions

Share your explorations:
- **Concepts you've explored** and what you discovered
- **Unexpected findings** that challenged your assumptions
- **Effective AI collaboration patterns** you've developed
- **Questions that led to interesting implementations**

*The best contributions show the journey from curiosity to understanding.*


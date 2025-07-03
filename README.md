# 🧵 Reddit Clone (Go) – Engine, REST API, Simulator & Digital Signature

A complete backend system for a Reddit-like platform built in **Go**, featuring:
- 🧠 Simulation Engine (Project 4.1)
- 🌐 REST API Interface & Client (Project 4.2)
- 🔐 Bonus: Digital Signature Verification

---

## 👥 Group Members
- **Mayank Garg**
- **Naman Tomar**

---

## 🔧 Features Overview

### ✅ Project 4.1: Reddit Engine & Zipf Simulator
- Account registration
- Subreddit creation, joining, leaving
- Post creation (text-only), reposting
- Hierarchical (nested) comments
- Upvotes/downvotes + Karma calculation
- Direct messaging (DMs) and replies
- Thousands of clients simulated using **Zipf distribution**
- Performance measurement (ops/sec, active users, response time)

### 🌐 Project 4.2: REST API & Client
- Fully RESTful API modeled after Reddit’s public API
- Endpoints for user management, subreddits, posts, comments, and DMs
- Simple REST client simulates interaction with server
- Logs of API interactions demonstrate communication

### 🔐 Bonus: Public Key-Based Digital Signatures
- Public key (RSA-2048 or ECC) submitted at registration
- Posts signed with user's private key
- API to retrieve user public keys
- Signature verified each time a post is accessed
- Uses standard Go cryptographic libraries

---

## 📁 Folder Structure

```
reddit-clone-2/
├── engine.go              # Core engine
├── main.go                # Simulator main
├── client.go              # Alternate simple client
├── server.go              # REST API backend
├── client_rest.go         # REST client simulation
├── tester.go              # Unit testing
├── crypto.go              # Digital signature (bonus)
├── go.mod                 # Module file
├── demo.mp4               # Project demo video
├── Project4_Report.pdf    # Report (submitted separately)
```

---

## 🚀 How to Run

1. **Install Go (1.19 or later)**

2. **Run Engine/REST Server**

```bash
go run server.go
```

3. **Run REST Client**

```bash
go run client_rest.go
```

4. **Run Load Simulator**

```bash
go run main.go
```

5. **Run Tests**

```bash
go run tester.go
```
---

## 📊 Performance Metrics (Sample)

| Users | Active | Posts | Comments | Votes | Throughput |
|-------|--------|-------|----------|-------|------------|
| 1000  | 700    | 5000  | 12520    | 20k+  | ~9500 ops/sec |

---

---

## 📜 License

Course: DOSP  

# Aron Project API Documentation

## Authentication

### User Signup
- **Endpoint:** `POST /auth/signup`
- **Description:** Register a new user.
- **Body:**
  ```json
  {
    "email": "user@example.com",
    "password": "password123",
    "role": "job_seeker" 
  }
  ```

### User Login
- **Endpoint:** `POST /auth/login`
- **Description:** Authenticate user and receive a JWT token.
- **Body:**
  ```json
  {
    "email": "user@example.com",
    "password": "password123"
  }
  ```

### Google Login
- **Endpoint:** `GET /auth/google/login`
- **Description:** Redirects to Google OAuth consent screen.

### Get Profile
- **Endpoint:** `GET /api/profile`
- **Headers:** `Authorization: Bearer <token>`
- **Description:** Get current user's basic profile info (ID, email, role).

---

## Job Provider Personal Profile

### Create/Update Provider Profile
- **Endpoint:** `POST /api/job-provider/profile`
- **Headers:** `Authorization: Bearer <token>`
- **Body:**
  ```json
  {
    "first_name": "Jane",
    "last_name": "Smith",
    "title": "HR Manager",
    "contact_number": "+1234567890",
    "bio": "Experienced recruiter..."
  }
  ```

### Get My Profile
- **Endpoint:** `GET /api/job-provider/profile`
- **Headers:** `Authorization: Bearer <token>`

### Get Specific Provider Profile
- **Endpoint:** `GET /api/job-provider/profile/:user_id`
- **Headers:** `Authorization: Bearer <token>`

---

## Company Management (Job Providers)

### Create Company Profile
- **Endpoint:** `POST /api/companies`
- **Headers:** `Authorization: Bearer <token>`
- **Body:**
  ```json
  {
    "company_name": "Tech Corp",
    "company_size": "50-100",
    "location": "New York",
    "industry": "IT",
    "services": "Software Development",
    "logo_url": "http://..."
  }
  ```

### Get My Companies
- **Endpoint:** `GET /api/companies`
- **Headers:** `Authorization: Bearer <token>`

### Get Company Details
- **Endpoint:** `GET /api/companies/:id`
- **Headers:** `Authorization: Bearer <token>`

### Update Company
- **Endpoint:** `PUT /api/companies/:id`
- **Headers:** `Authorization: Bearer <token>`

### Add Product to Company
- **Endpoint:** `POST /api/companies/:id/products`
- **Body:**
  ```json
  {
    "name": "SaaS Platform",
    "description": "Cloud solution",
    "price": 99.99
  }
  ```

---

## Job Management (Job Providers)

### Create Job Posting
- **Endpoint:** `POST /api/companies/:company_id/jobs`
- **Body:**
  ```json
  {
    "title": "Backend Developer",
    "description": "Go developer needed...",
    "requirements": "3+ years exp",
    "field": "Engineering",
    "location": "Remote",
    "type": "Full-time",
    "salary_range": "$80k-$120k"
  }
  ```

### Get Company Jobs
- **Endpoint:** `GET /api/companies/:company_id/jobs`

### Update Job
- **Endpoint:** `PUT /api/jobs/:job_id`

### Delete Job
- **Endpoint:** `DELETE /api/jobs/:job_id`

### View Applicants
- **Endpoint:** `GET /api/jobs/:job_id/applicants`
- **Description:** Get list of users who applied to this job.

### Update Application Status
- **Endpoint:** `PUT /api/jobs/applications/:application_id/status`
- **Body:**
  ```json
  {
    "status": "interviewed" // applied, reviewing, interviewed, rejected, hired
  }
  ```

---

## Job Seekers

### Create/Update Profile
- **Endpoint:** `POST /api/job-seeker/profile`
- **Body:**
  ```json
  {
    "first_name": "John",
    "last_name": "Doe",
    "skills": "Go, Python",
    "is_open_to_work": true
  }
  ```

### Add Internship
- **Endpoint:** `POST /api/job-seeker/internships`
- **Body:**
  ```json
  {
    "company": "Startup Inc",
    "role": "Intern",
    "start_date": "2023-01-01T00:00:00Z"
  }
  ```

### Job Discovery (Swipe Deck)
- **Endpoint:** `GET /api/job-seeker/jobs/discovery`
- **Description:** Get a list of jobs to swipe on.

### Swipe (Apply/Pass)
- **Endpoint:** `POST /api/job-seeker/jobs/swipe`
- **Body:**
  ```json
  {
    "job_id": 123,
    "action": "like" // "like" (Apply) or "pass"
  }
  ```

### Search Jobs
- **Endpoint:** `GET /api/job-seeker/jobs/search?field=IT&location=Remote`

---

## Chat

### WebSocket Connection
- **Endpoint:** `GET /api/chat/ws`
- **Description:** Connect to WebSocket server for real-time chat.

### Send Message (HTTP)
- **Endpoint:** `POST /api/chat/send`
- **Body:**
  ```json
  {
    "receiver_id": 45,
    "content": "Hello, regarding the job..."
  }
  ```

### Get Chat History
- **Endpoint:** `GET /api/chat/history/:user_id`

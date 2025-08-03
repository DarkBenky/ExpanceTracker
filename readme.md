# Expense Tracker API

## API Endpoints

### 1. User Registration
**POST** `/register`

Register a new user and receive a JWT token.

**Request:**
```json
{
  "username": "john_doe",
  "password": "secure_password123"
}
```

**Response:**
```json
{
  "message": "User registered successfully",
  "user_id": 1,
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

---

### 2. User Login
**POST** `/login`

Authenticate user and receive a JWT token.

**Request:**
```json
{
  "username": "john_doe",
  "password": "secure_password123"
}
```

**Response:**
```json
{
  "message": "Login successful",
  "user_id": 1,
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

---

### 3. Create Group
**POST** `/groups`

Create a new expense group. The creator becomes the owner and is automatically added as a member.

**Request:**
```json
{
  "token": "your-jwt-token-here",
  "name": "Family Expenses"
}
```

**Response:**
```json
{
  "message": "Group created successfully",
  "group_id": 1
}
```

---

### 4. Add User to Group
**POST** `/groups/members`

Add a user to an existing group. Only the group owner can perform this action.

**Request:**
```json
{
  "token": "group-owner-jwt-token",
  "group_id": 1,
  "user_id": 2
}
```

**Response:**
```json
{
  "message": "User added to group successfully"
}
```

---

### 5. Add Expense
**POST** `/expenses`

Add a new expense to a group. Only group members can add expenses.

**Request:**
```json
{
  "token": "your-jwt-token-here",
  "user_id": 1,
  "group_id": 1,
  "description": "Lunch at restaurant",
  "amount": 25.50,
  "category": "Food",
  "date": "2025-08-04"
}
```

**Response:**
```json
{
  "message": "Expense added successfully"
}
```

---

### 6. Remove Expense
**DELETE** `/expenses`

Remove an expense from a group. Only group members can remove expenses.

**Request:**
```json
{
  "token": "your-jwt-token-here",
  "user_id": 1,
  "group_id": 1,
  "expense_id": 5
}
```

**Response:**
```json
{
  "message": "Expense removed successfully"
}
```
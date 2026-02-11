# Plane Service Documentation

The Plane Service provides a comprehensive API for managing aircraft and their parts, including usage tracking and maintenance monitoring.

## Table of Contents

- [Overview](#overview)
- [API Endpoints](#api-endpoints)
  - [Planes](#planes)
  - [Plane Parts](#plane-parts)
  - [Maintenance](#maintenance)
- [Usage Examples](#usage-examples)
- [Error Handling](#error-handling)

---

## Overview

The Plane Service allows users to:
- Register and manage aircraft (planes)
- Add and track parts installed on each plane
- Monitor usage hours and maintenance thresholds
- Get alerts for parts requiring maintenance

All endpoints require JWT authentication except for the initial setup.



---

## API Endpoints

### Planes

#### Create a Plane

**Endpoint:** `POST /api/planes`

**Request Body:**
```json
{
  "tail_number": "N12345",
  "model": "Boeing 737-800"
}
```

**Response (201 Created):**
```json
{
  "id": 1,
  "tail_number": "N12345",
  "model": "Boeing 737-800",
  "created_at": "2024-01-15T10:30:00Z"
}
```

**Validation Rules:**
- `tail_number`: Required, 2-50 characters, must be unique
- `model`: Required, 2-100 characters

---

#### Get All Planes

**Endpoint:** `GET /api/planes`

**Response (200 OK):**
```json
[
  {
    "id": 1,
    "tail_number": "N12345",
    "model": "Boeing 737-800",
    "created_at": "2024-01-15T10:30:00Z"
  },
  {
    "id": 2,
    "tail_number": "N67890",
    "model": "Airbus A320",
    "created_at": "2024-01-16T14:20:00Z"
  }
]
```

---

#### Get Plane by ID

**Endpoint:** `GET /api/planes/:id`

**Response (200 OK):**
```json
{
  "id": 1,
  "tail_number": "N12345",
  "model": "Boeing 737-800",
  "created_at": "2024-01-15T10:30:00Z"
}
```

---

#### Get Plane by Tail Number

**Endpoint:** `GET /api/planes/tail/:tail_number`

**Example:** `GET /api/planes/tail/N12345`

**Response (200 OK):**
```json
{
  "id": 1,
  "tail_number": "N12345",
  "model": "Boeing 737-800",
  "created_at": "2024-01-15T10:30:00Z"
}
```

---

#### Update a Plane

**Endpoint:** `PUT /api/planes/:id`

**Request Body:**
```json
{
  "tail_number": "N12346",
  "model": "Boeing 737-900"
}
```

**Response (200 OK):**
```json
{
  "id": 1,
  "tail_number": "N12346",
  "model": "Boeing 737-900",
  "created_at": "2024-01-15T10:30:00Z"
}
```

---

#### Delete a Plane

**Endpoint:** `DELETE /api/planes/:id`

**Response:** `204 No Content`

---

#### Get Plane with All Parts

**Endpoint:** `GET /api/planes/:id/with-parts`

**Response (200 OK):**
```json
{
  "plane": {
    "id": 1,
    "tail_number": "N12345",
    "model": "Boeing 737-800",
    "created_at": "2024-01-15T10:30:00Z"
  },
  "parts": [
    {
      "id": 1,
      "plane_id": 1,
      "part_name": "Engine Fan Blade",
      "serial_number": "SN-ENG-001",
      "category": "engine",
      "usage_hours": 1250.5,
      "usage_limit_hours": 5000,
      "usage_percent": 25.01,
      "installed_at": "2024-01-15T10:30:00Z"
    }
  ]
}
```

---

### Plane Parts

#### Add a Part to a Plane

**Endpoint:** `POST /api/planes/:planeId/parts`

**Request Body:**
```json
{
  "plane_id": 1,
  "part_name": "Engine Fan Blade",
  "serial_number": "SN-ENG-001",
  "category": "engine",
  "usage_hours": 0,
  "usage_limit_hours": 5000
}
```

**Validation Rules:**
- `plane_id`: Required, must reference an existing plane
- `part_name`: Required, 2-255 characters
- `serial_number`: Required, 2-100 characters, must be unique
- `category`: Required, 2-150 characters
- `usage_hours`: Optional, default 0
- `usage_limit_hours`: Required, must be greater than 0

**Response (201 Created):**
```json
{
  "id": 1,
  "plane_id": 1,
  "part_name": "Engine Fan Blade",
  "serial_number": "SN-ENG-001",
  "category": "engine",
  "usage_hours": 0,
  "usage_limit_hours": 5000,
  "usage_percent": 0,
  "installed_at": "2024-01-15T10:30:00Z"
}
```

---

#### Get All Parts for a Plane

**Endpoint:** `GET /api/planes/:planeId/parts`

**Query Parameters:**
- `category` (optional): Filter by category

**Example:** `GET /api/planes/1/parts?category=engine`

**Response (200 OK):**
```json
[
  {
    "id": 1,
    "plane_id": 1,
    "part_name": "Engine Fan Blade",
    "serial_number": "SN-ENG-001",
    "category": "engine",
    "usage_hours": 1250.5,
    "usage_limit_hours": 5000,
    "usage_percent": 25.01,
    "installed_at": "2024-01-15T10:30:00Z"
  }
]
```

---

#### Get All Parts (Global)

**Endpoint:** `GET /api/planes/parts`

**Response (200 OK):**
```json
[
  {
    "id": 1,
    "plane_id": 1,
    "part_name": "Engine Fan Blade",
    "serial_number": "SN-ENG-001",
    "category": "engine",
    "usage_hours": 1250.5,
    "usage_limit_hours": 5000,
    "usage_percent": 25.01,
    "installed_at": "2024-01-15T10:30:00Z"
  }
]
```

---

#### Get Part by ID

**Endpoint:** `GET /api/planes/parts/:partId`

**Response (200 OK):**
```json
{
  "id": 1,
  "plane_id": 1,
  "part_name": "Engine Fan Blade",
  "serial_number": "SN-ENG-001",
  "category": "engine",
  "usage_hours": 1250.5,
  "usage_limit_hours": 5000,
  "usage_percent": 25.01,
  "installed_at": "2024-01-15T10:30:00Z"
}
```

---

#### Update Part Details

**Endpoint:** `PUT /api/planes/parts/:partId`

**Request Body:**
```json
{
  "part_name": "High-Performance Fan Blade",
  "category": "engine",
  "serial_number": "SN-ENG-001-UPDATED"
}
```

**Response (200 OK):**
```json
{
  "id": 1,
  "plane_id": 1,
  "part_name": "High-Performance Fan Blade",
  "serial_number": "SN-ENG-001-UPDATED",
  "category": "engine",
  "usage_hours": 1250.5,
  "usage_limit_hours": 5000,
  "usage_percent": 25.01,
  "installed_at": "2024-01-15T10:30:00Z"
}
```

---

#### Update Part Usage Hours

**Endpoint:** `PUT /api/planes/parts/:partId/usage`

**Request Body:**
```json
{
  "usage_hours": 2500.75
}
```

**Validation:**
- `usage_hours`: Must be greater than or equal to 0, cannot exceed `usage_limit_hours`

**Response (200 OK):**
```json
{
  "id": 1,
  "plane_id": 1,
  "part_name": "Engine Fan Blade",
  "serial_number": "SN-ENG-001",
  "category": "engine",
  "usage_hours": 2500.75,
  "usage_limit_hours": 5000,
  "usage_percent": 50.02,
  "installed_at": "2024-01-15T10:30:00Z"
}
```

---

#### Delete a Part

**Endpoint:** `DELETE /api/planes/parts/:partId`

**Response:** `204 No Content`

---

### Maintenance

#### Get Parts Needing Maintenance

**Endpoint:** `GET /api/planes/maintenance/alerts`

**Query Parameters:**
- `threshold` (optional): Percentage threshold, default 80

**Example:** `GET /api/planes/maintenance/alerts?threshold=70`

**Response (200 OK):**
```json
[
  {
    "id": 2,
    "plane_id": 1,
    "part_name": "Brake Pad Set",
    "serial_number": "SN-BRAKE-005",
    "category": "brakes",
    "usage_hours": 450,
    "usage_limit_hours": 500,
    "usage_percent": 90,
    "installed_at": "2024-01-10T08:00:00Z"
  },
  {
    "id": 5,
    "plane_id": 2,
    "part_name": "Tire Assembly",
    "serial_number": "SN-TIRE-012",
    "category": "landing_gear",
    "usage_hours": 360,
    "usage_limit_hours": 500,
    "usage_percent": 72,
    "installed_at": "2024-01-12T12:00:00Z"
  }
]
```


## Usage Examples

### Complete Workflow

#### 1. Login to get JWT token
```bash
POST /api/users/login
{
  "name": "maintenance_manager",
  "password": "your-password"
}
```

#### 2. Create a plane
```bash
POST /api/planes
Authorization: Bearer <token>
{
  "tail_number": "N737MAX",
  "model": "Boeing 737 MAX 8"
}
```

#### 3. Add parts to the plane
```bash
POST /api/planes/1/parts
Authorization: Bearer <token>
{
  "plane_id": 1,
  "part_name": "Left Engine",
  "serial_number": "SN-ENG-2024-001",
  "category": "engine",
  "usage_hours": 0,
  "usage_limit_hours": 10000
}
```

#### 4. Track usage over time
```bash
# After 1000 flight hours
PUT /api/planes/parts/1/usage
Authorization: Bearer <token>
{
  "usage_hours": 1000
}
```

#### 5. Check maintenance alerts
```bash
GET /api/planes/maintenance/alerts?threshold=80
Authorization: Bearer <token>
```

---

### cURL Examples

**Create Plane:**
```bash
curl -X POST http://localhost:8080/api/planes \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"tail_number": "N12345", "model": "Boeing 737-800"}'
```

**Add Part:**
```bash
curl -X POST http://localhost:8080/api/planes/1/parts \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "plane_id": 1,
    "part_name": "Engine Fan Blade",
    "serial_number": "SN-ENG-001",
    "category": "engine",
    "usage_hours": 0,
    "usage_limit_hours": 5000
  }'
```

**Update Usage:**
```bash
curl -X PUT http://localhost:8080/api/planes/parts/1/usage \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"usage_hours": 2500.5}'
```

**Get Maintenance Alerts:**
```bash
curl -X GET "http://localhost:8080/api/planes/maintenance/alerts?threshold=80" \
  -H "Authorization: Bearer <token>"
```

---

## Error Handling

### Common Error Responses

| Status | Error | Description |
|--------|-------|-------------|
| 400 | invalid plane ID | Invalid ID parameter |
| 400 | invalid threshold value | Invalid query parameter |
| 401 | unauthorized | Missing or invalid JWT token |
| 404 | plane not found | Plane does not exist |
| 404 | plane part not found | Part does not exist |
| 409 | plane with this tail number already exists | Duplicate tail number |
| 409 | plane part with this serial number already exists | Duplicate serial number |
| 500 | internal server error | Server error |

### Error Response Format
```json
{
  "error": "error message description"
}
```

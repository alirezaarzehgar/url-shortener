# API Design

From user perspective this system just should have NNN endpoints.

## Endpoints

### 1. Create a Short Link

**POST** `/new`

**Request Body** (JSON):

```json
{
  "link": "https://example.com/long-url"
}
```

**Responses:**

* `200 OK`

```json
{
  "key": "a1B2c3"  // 6-character Base64-encoded string
}
```

* `500 Internal Server Error` – Something went wrong

#### Security and abuse protection

The endpoint is rate-limited by IP address to prevent flooding and DDOS.

---

### 2. Redirect to Original Link

**GET** `/{key}`

**Responses:**

* `301 Moved Permanently` – Redirects to original link
* `404 Not Found` – Short link not found
* `500 Internal Server Error` – Something went wrong

---

### Notes

* No authentication required.
* Short links are permanent.
* Maximum original link length: 1KB.

# 🖼️ Go Image Host

A super-simple self-hosted image upload app built with Go.  
Includes:

- 🔐 Password-protected login (bcrypt hash only, no username)
- 🖼️ Image uploads via browser
- 🔗 Shareable direct URLs (`/i/filename.jpg`)
- 🗑️ Delete images from web UI
- 🐳 Docker support (with volume for persistent storage)

---

## 🚀 Getting Started

### 🔧 Requirements

- Go 1.23+ (or Docker)
- HTML5 browser
- (Optional) Tailscale / Nginx for remote access or HTTPS

---

### 🏁 Run locally (without Docker)

```bash
go mod tidy
go run main.go
```

Visit: [http://localhost:8080](http://localhost:8080)  
Login with password: `supersecret` (default hash set in `main.go`)

---

### 🐳 Run with Docker

```bash
docker build -t go-img .
docker run -d \
  --restart=always \
  -p 8765:8765 \
  -v "$PWD/uploads:/app/uploads" \
  go-img
```

Visit: [http://localhost:8765](http://localhost:8765)

---

## 🔐 Change the login password

1. Generate a new bcrypt hash:

   ```bash
   go run golang.org/x/crypto/bcrypt@latest
   ```

2. Replace the hash in `main.go`:

   ```go
   const hash = "your-new-bcrypt-hash-here"
   ```

3. Rebuild the app.

---

## 📁 Folder Structure

```
.
├── main.go           # Core app logic
├── login.html        # Login page
├── manage.html       # Upload + image gallery page
├── uploads/          # Uploaded image files (mounted in Docker)
├── Dockerfile
├── .dockerignore
└── .gitignore
```

---

## ✅ Features

- ✅ Login-protected UI
- ✅ Upload `.png`, `.jpg`, `.jpeg`, `.gif`
- ✅ Direct URLs for images
- ✅ Delete with one click
- ✅ No external database or frameworks

---

## 🛡️ Security Notes

- Password is stored as a bcrypt hash (no plaintext risk)
- Still super basic: no CSRF, no rate limits, no HTTPS — run behind a reverse proxy or VPN (e.g. [Tailscale](https://tailscale.com/)) for production use
- Add authentication middleware if you plan to go public

---

## 📜 License

MIT — use it, ship it, tweak it. You own it.

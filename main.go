// go run main.go
package main

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

const (
	addr       = ":8765"
	hash       = "$2a$12$WLsgMPA3aAp7dMd/zJZcvO1lL0FI3Db50zzNFo0O3OWnury3DcqLG" // <— bcrypt hash
	cookieName = "auth"
)

var (
	loginTmpl  = template.Must(template.ParseFiles("login.html"))
	manageTmpl = template.Must(template.ParseFiles("manage.html"))
)

func main() {
	if err := os.MkdirAll("uploads", 0o755); err != nil {
		log.Fatal(err)
	}

	http.Handle("/i/", http.StripPrefix("/i/", http.FileServer(http.Dir("uploads"))))
	http.HandleFunc("/", root)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/manage", manage)
	http.HandleFunc("/upload", upload)
	http.HandleFunc("/delete", del)

	log.Printf("🚀 running on http://localhost%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

// ---------- auth helpers ----------
func authorised(r *http.Request) bool {
	c, err := r.Cookie(cookieName)
	return err == nil && c.Value == "1"
}

func requireAuth(w http.ResponseWriter, r *http.Request) bool {
	if !authorised(r) {
		http.Redirect(w, r, "/login", http.StatusFound)
		return false
	}
	return true
}

// ---------- handlers ----------
func root(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/manage", http.StatusFound)
}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		_ = loginTmpl.Execute(w, r.URL.Query().Get("bad") == "1")
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(r.FormValue("password"))); err == nil {
		http.SetCookie(w, &http.Cookie{Name: cookieName, Value: "1", Path: "/", HttpOnly: true})
		http.Redirect(w, r, "/manage", http.StatusFound)
		return
	}
	http.Redirect(w, r, "/login?bad=1", http.StatusFound)
}

func logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{Name: cookieName, Value: "", Path: "/", MaxAge: -1})
	http.Redirect(w, r, "/login", http.StatusFound)
}

func manage(w http.ResponseWriter, r *http.Request) {
	if !requireAuth(w, r) {
		return
	}
	entries, _ := os.ReadDir("uploads")
	var images []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		l := strings.ToLower(name)
		if strings.HasSuffix(l, ".png") || strings.HasSuffix(l, ".jpg") ||
			strings.HasSuffix(l, ".jpeg") || strings.HasSuffix(l, ".gif") {
			images = append(images, name)
		}
	}
	_ = manageTmpl.Execute(w, images)
}

func upload(w http.ResponseWriter, r *http.Request) {
	if !requireAuth(w, r) {
		return
	}
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/manage", http.StatusFound)
		return
	}
	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "bad upload", http.StatusBadRequest)
		return
	}
	defer file.Close()

	dstPath := filepath.Join("uploads", header.Filename)
	dst, err := os.Create(dstPath)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	defer dst.Close()
	if _, err = io.Copy(dst, file); err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/manage", http.StatusFound)
}

func del(w http.ResponseWriter, r *http.Request) {
	if !requireAuth(w, r) {
		return
	}
	name := r.URL.Query().Get("img")
	if name == "" {
		http.Redirect(w, r, "/manage", http.StatusFound)
		return
	}
	if err := os.Remove(filepath.Join("uploads", filepath.Clean(name))); err != nil {
		log.Println("delete error:", err)
	}
	http.Redirect(w, r, "/manage", http.StatusFound)
}

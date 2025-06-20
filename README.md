# 🎂 BirthdayMate

**BirthdayMate** is a Go-based web application that lets users:

* Enter their date of birth
* View their current age
* Discover celebrities who share their birthday

---

## 🧰 Tech Stack

* **Language:** Go (Golang)
* **Containerization:** Docker
* **CI:** GitHub Actions
* **Docker Registry:** Docker Hub

---

## 📁 Project Structure

```
.
├── Dockerfile                 # Docker build instructions
├── go.mod                     # Go module file
├── main.go                    # Application source code
├── Makefile                   # (Optional) Build commands
└── .github/workflows/ci.yaml  # CI pipeline config
```

---

## 🚀 Steps Followed to Achieve Working CI (GitHub Actions)

Here's what was done step-by-step to make CI fully functional:

### 1. ✅ Initialize Go Project

```bash
go mod init birthdaymate
go mod tidy
```

### 2. ✅ Write Application Code

Created `main.go` to accept user input and print age + matching celebrity birthdays.

### 3. ✅ Dockerize the App

Created a multi-stage `Dockerfile` to build the Go binary and copy it into a minimal Alpine image:

```dockerfile
# Build Stage
FROM golang:1.23 AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN go build -o app main.go

# Final Image
FROM alpine:latest
WORKDIR /root
COPY --from=builder /app/app .
EXPOSE 9090
CMD ["./app"]
```

### 4. ✅ Build Docker Image Locally

```bash
docker build -t vinaypdb/birthdaymate:latest .
```

### 5. ✅ Push Image to Docker Hub

```bash
docker login
docker push vinaypdb/birthdaymate:latest
```

### 6. ✅ Set Up CI with GitHub Actions

Created `.github/workflows/ci.yaml` with these steps:

* Checkout source code
* Set up Go environment
* Build & test Go app
* Log in to Docker Hub
* Build and push Docker image

```yaml
name: CI Pipeline

on:
  push:
    branches: [main]

jobs:
  build-and-push:
    runs-on: ubuntu-latest

    steps:
    - name: Checkout source code
      uses: actions/checkout@v3

    - name: Set up Go
      uses: actions/setup-go@v5
      with:
        go-version: '1.21'

    - name: Build Go app
      run: go build -v ./...

    - name: Run Go tests
      run: go test -v ./...

    - name: Log in to DockerHub
      uses: docker/login-action@v3
      with:
        username: ${{ secrets.DOCKER_USERNAME }}
        password: ${{ secrets.DOCKER_PASSWORD }}

    - name: Build Docker image
      run: docker build -t ${{ secrets.DOCKER_USERNAME }}/birthdaymate:latest .

    - name: Push Docker image
      run: docker push ${{ secrets.DOCKER_USERNAME }}/birthdaymate:latest
```

### 7. ✅ Add GitHub Secrets

In the GitHub repository → **Settings → Secrets → Actions**, added:

* `DOCKER_USERNAME`
* `DOCKER_PASSWORD`

### 8. ✅ Triggered CI by Pushing to `main`

```bash
git add .
git commit -m "✅ Setup complete: Go app + Docker + GitHub Actions"
git push origin main
```

️➡️ CI ran automatically and pushed the Docker image to Docker Hub successfully!

---

## 📆 Docker Hub

Pull the built image from Docker Hub:

```bash
docker pull vinaypdb/birthdaymate:latest
```

---

## ✅ Next Steps (Optional)

* [ ] Setup Helm chart for Kubernetes deployment
* [ ] Configure Argo CD for GitOps delivery
* [ ] Deploy to Amazon EKS

---

## 🙌 Author

**Vinay Pedapuri**
[Docker Hub](https://hub.docker.com/u/vinaypdb) • [GitHub](https://github.com/vinaypdb)


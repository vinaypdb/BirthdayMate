
# 🎂 BirthdayMate

**BirthdayMate** is a Go‑based web application that lets users:

* Enter their date of birth
* View their current age
* know how many times their birthday has fallen on a sundays so far
* Discover celebrities who share their birthday

---

## 🧰 Tech Stack

* **Language:** Go (Golang)
* **Containerization:** Docker
* **CI:** GitHub Actions
* **Security Scans:** Trivy & OWASP Dependency‑Check
* **Registry:** Docker Hub

---

## 📁 Project Structure

```text
.
├── Dockerfile                 # Docker build instructions
├── go.mod                     # Go module file
├── main.go                    # Application source code
├── Makefile                   # (Optional) build/run helpers
└── .github/workflows/ci.yaml  # CI pipeline config
```

---

## 🚀 How We Built a Passing CI Pipeline

Below is the exact sequence of steps we followed.

### 1  ✅ Create & Initialise Go Module

```bash
mkdir BirthdayMate && cd BirthdayMate
go mod init birthdaymate
go mod tidy
```

### 2  ✅ Write Application Code

`main.go` accepts a birth‑date, calculates age and lists celebrity “birthday twins”.

### 3  ✅ Dockerise the App

Multi‑stage `Dockerfile` (Alpine final image, 19 MB):

```dockerfile
# ---------- Build stage ----------
FROM golang:1.23 AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN go build -o app main.go

# ---------- Final stage ----------
FROM alpine:latest
WORKDIR /root
COPY --from=builder /app/app .
EXPOSE 9090
CMD ["./app"]
```

### 4  ✅ Local Image Build & Test

```bash
docker build -t vinaypdb/birthdaymate:latest .
```

### 5  ✅ Push Image to Docker Hub (Manual First Push)

```bash
docker login
docker push vinaypdb/birthdaymate:latest
```

### 6  ✅ Add CI Workflow (`.github/workflows/ci.yaml`)

Main stages:

* Checkout ⬇️
* Go build & unit tests ✅
* **OWASP Dependency‑Check** (Go modules CVE scan) 🛡
* Docker login & image build 🐳
* **Trivy image scan** (fail on HIGH/CRITICAL) 🔍
* Push to Docker Hub ☁️

```yaml
name: CI Pipeline

on:
  push:
    branches: [main]

jobs:
  build-and-push:
    runs-on: ubuntu-latest

    steps:
    - name: Checkout code
      uses: actions/checkout@v3

    - name: Set up Go
      uses: actions/setup-go@v5
      with:
        go-version: "1.21"

    - name: Build & test
      run: |
        go build -v ./...
        go test  -v ./...

    - name: OWASP Dependency‑Check
      uses: dependency-check/Dependency-Check_Action@main
      with:
        project: "BirthdayMate"
        path: "."
        format: "SARIF"
        out: "dependency-check-report"

    - name: Log in to DockerHub
      uses: docker/login-action@v3
      with:
        username: ${{ secrets.DOCKER_USERNAME }}
        password: ${{ secrets.DOCKER_PASSWORD }}

    - name: Build Docker image
      run: docker build -t ${{ secrets.DOCKER_USERNAME }}/birthdaymate:latest .

    - name: Trivy Image Scan
      uses: aquasecurity/trivy-action@0.12.0
      with:
        image-ref: ${{ secrets.DOCKER_USERNAME }}/birthdaymate:latest
        format: table
        exit-code: 1          # fail on HIGH/CRITICAL
        ignore-unfixed: true

    - name: Push Docker image
      run: docker push ${{ secrets.DOCKER_USERNAME }}/birthdaymate:latest
```

### 7  ✅ Add GitHub Secrets

`Settings → Secrets → Actions`

* `DOCKER_USERNAME`
* `DOCKER_PASSWORD`

### 8  ✅ Commit & Push — CI Passes

```bash
git add .
git commit -m "🎉 Fully‑automated CI with security scans"
git push origin main
```

GitHub Actions now builds, scans and publishes the image **automatically** on every push to `main`.

---

## 🐳 Docker Hub

```bash
docker pull vinaypdb/birthdaymate:latest
```

---

## 📌 Next Steps (Part 2)

* [ ] Scaffold Helm chart for Kubernetes deployments
* [ ] Configure Argo CD for GitOps CD
* [ ] Deploy to Amazon EKS (Terraform)

---

## 🙌 Author

**Vinay Pedapuri**  ⋅  [Docker Hub](https://hub.docker.com/u/vinaypdb) ⋅  [GitHub](https://github.com/vinaypdb)


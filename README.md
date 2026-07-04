# HTTP Reverse Proxy & Load Balancer

A production-inspired HTTP reverse proxy and load balancer built in Go.

It distributes incoming requests across multiple backend servers using a thread-safe round-robin algorithm, performs active health checks, and automatically removes unhealthy servers until they recover.

## Features

- HTTP Reverse Proxy
- Round Robin Load Balancing
- Active Health Checks
- Automatic Failover
- Request Logging
- Panic Recovery
- YAML Configuration
- Graceful Shutdown
- Concurrent Request Handling

---

## Architecture

```
            Client
               │
               ▼
      Reverse Proxy / Load Balancer
         │         │         │
         ▼         ▼         ▼
     Backend1  Backend2  Backend3
```

---

## Tech Stack

- Go
- net/http
- httputil
- sync/atomic
- YAML

---

## Project Structure

```
cmd/
internal/
backend/
configs/
```

---

## Run

```bash
go mod tidy
go run cmd/main.go
```

Start multiple backend servers.

```bash
PORT=8081 go run backend/main.go
PORT=8082 go run backend/main.go
PORT=8083 go run backend/main.go
```

Then

```bash
curl http://localhost:8080
```

---

## Future Improvements

- Docker
- AWS EC2 Deployment
- HTTPS Support
- Metrics Dashboard
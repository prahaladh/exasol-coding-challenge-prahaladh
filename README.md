# Exasol Proof-of-Work Client (Go)

## Overview

This project implements a high-performance client in **Go** to solve Exasol’s Proof-of-Work (PoW) challenge over a **TLS-based TCP protocol**.

The client:

* Establishes a secure TLS connection
* Handles a custom line-based protocol
* Solves a computational PoW challenge
* Submits user details upon successful validation

---

## Key Concepts Used

* **Concurrency (Goroutines)**
  Parallel brute-force search using CPU cores

* **Networking (TCP + TLS)**
  Secure communication using `tls.Dial`

* **Cryptography (SHA1)**
  Hash-based PoW validation

* **Performance Optimization**

  * Work distribution across goroutines
  * Buffer reuse to reduce GC overhead

---

## How It Works

1. Connect to Exasol server via TLS
2. Perform handshake (`HELO`)
3. Solve PoW:

   * Find suffix such that:

     ```
     SHA1(authdata + suffix) starts with N zeros
     ```
4. Respond to server challenges:

   * NAME, MAIL, COUNTRY, etc.
5. Submit and receive `END` confirmation

---

## Proof-of-Work Strategy

* Uses **`runtime.NumCPU()` workers**
* Each worker processes:

  ```
  i = start + k * workers
  ```
* Ensures:

  * No duplicate work
  * Maximum CPU utilization

---

## TLS Configuration

* Uses client certificate authentication
* `.pem` file contains both:

  * Certificate
  * Private key

---

## Running the Project

### 1. Set environment variables

#### Windows (PowerShell)

```
$env:TLS_CERT="client.pem"
$env:TLS_KEY="client.pem"
```

#### Linux / Mac

```
export TLS_CERT=client.pem
export TLS_KEY=client.pem
```

---

### 2. Install the dependencies and Run the Code

```
go mod tidy
go run main.go
```

---

## Project Structure

```
.
├── main.go
├── go.mod
├── README.md
└── client.pem (NOT committed)
```

---

## Security Note

Sensitive files are excluded using `.gitignore`:

```
*.pem
*.key
*.crt
```

> Never commit private keys to version control.

---

## Performance Notes

* Efficiently handles difficulty levels up to **6+**
* Uses:

  * Buffer reuse (`buf[:0]`)
  * Parallel hashing
* Designed to minimize GC overhead

---

## Example Output

```
Server: HELO
Server: POW ... 6
Solving POW...
Server: END
✅ Completed successfully
```

---

## Why Go?

* Lightweight concurrency (goroutines)
* True parallel execution (no GIL)
* Efficient memory management
* Ideal for network + CPU-bound workloads

---

## Author

**Prahaladh HN**
SRE | Platform Engineering | Kubernetes | Go

---

## Conclusion

This project demonstrates:

* Strong problem-solving skills
* Understanding of distributed systems
* Performance-oriented engineering

---

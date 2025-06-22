#!/bin/bash

# DePIN Compatibility API Setup Script
# This script automates the complete setup process

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
PROJECT_NAME="depin-compatibility-api"
GO_MIN_VERSION="1.21"
DEFAULT_PORT="8080"

# Print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Compare version numbers
version_gt() {
    test "$(printf '%s\n' "$@" | sort -V | head -n 1)" != "$1"
}

# Banner
print_banner() {
    echo -e "${BLUE}"
    echo "================================================================"
    echo "           DePIN Compatibility API Setup Script"
    echo "================================================================"
    echo -e "${NC}"
}

# Check prerequisites
check_prerequisites() {
    print_status "Checking prerequisites..."
    
    # Check if Go is installed
    if ! command_exists go; then
        print_error "Go is not installed. Please install Go $GO_MIN_VERSION or higher."
        echo "Visit: https://golang.org/doc/install"
        exit 1
    fi
    
    # Check Go version
    GO_VERSION=$(go version | cut -d' ' -f3 | tr -d 'go')
    if version_gt $GO_MIN_VERSION $GO_VERSION; then
        print_error "Go version $GO_VERSION is too old. Please upgrade to $GO_MIN_VERSION or higher."
        exit 1
    fi
    
    print_success "Go $GO_VERSION found"
    
    # Check for make (optional)
    if command_exists make; then
        print_success "Make found"
    else
        print_warning "Make not found. You'll need to run commands manually."
    fi
    
    # Check for Docker (optional)
    if command_exists docker; then
        print_success "Docker found"
        DOCKER_AVAILABLE=true
    else
        print_warning "Docker not found. Container deployment will be unavailable."
        DOCKER_AVAILABLE=false
    fi
    
    # Check for curl (for testing)
    if command_exists curl; then
        print_success "curl found"
    else
        print_warning "curl not found. API testing will be limited."
    fi
}

# Create project structure
create_project_structure() {
    print_status "Creating project structure..."
    
    # Create directories
    mkdir -p internal/{api,models,service,data}
    mkdir -p {scripts,examples,docs,tests,data}
    mkdir -p examples/{curl,python,javascript}
    
    print_success "Project structure created"
}

# Create sample data
create_sample_data() {
    print_status "Creating sample DePIN data..."
    
    # Check if user provided depin_specifications_final.csv exists
    if [ -f "depin_specifications_final.csv" ]; then
        print_success "Using existing depin_specifications_final.csv file"
        cp depin_specifications_final.csv data/depin_specs.csv
        return 0
    fi
    
    print_status "Creating sample CSV based on depin_specifications_final.csv structure..."
    
    cat > data/depin_specs.csv << 'EOF'
project_name,project_type,node_type,cpu_cores_min,cpu_architecture,ram_gb_min,ram_gb_recommended,storage_gb_min,storage_type,gpu_required,gpu_vram_gb_min,gpu_requirements,network_speed_mbps_min,network_type,supported_os,blockchain_network,token_symbol,estimated_monthly_cost_usd_min,estimated_monthly_cost_usd_max,cost_category,additional_requirements,home_friendly,last_updated
Filecoin station,Storage,Light Node,4,Any,8,16,500,SSD,FALSE,0,None,20,Broadband,"Linux,Windows,macOS",Filecoin,FIL,10,50,Low,For wallet and blockchain interaction only,TRUE,2025-01-01
AIOZ,CDN/AI,Standard Node,1,Any,0.5,2,50,Any,FALSE,0,None,20,Broadband,"Linux,Windows,macOS",Ethereum,AIOZ,5,20,Very Low,Additional 20GB for AI tasks,TRUE,2025-01-01
Render Network,Compute,GPU Node,4,Any,8,16,250,SSD,TRUE,8,RTX 3070 or better,100,Broadband,"Linux,Windows",Ethereum,RNDR,50,200,Medium,High-end GPU required for 3D rendering,TRUE,2025-01-01
Livepeer,Video,Transcoder,8,Any,16,32,1000,SSD,TRUE,6,GPU with H.264 encoding,200,Broadband,Linux,Ethereum,LPT,100,500,High,Professional video transcoding setup,FALSE,2025-01-01
Helium IoT,IoT,Hotspot,1,ARM/x86,2,4,32,Any,FALSE,0,None,10,Broadband,Linux,Helium,HNT,20,80,Low,IoT hotspot for network coverage,TRUE,2025-01-01
Akash Network,Compute,Provider,8,Any,32,64,500,SSD,FALSE,0,None,100,Broadband,Linux,Cosmos,AKT,80,300,Medium,Kubernetes provider node,FALSE,2025-01-01
Theta Network,Video,Edge Node,4,Any,8,16,200,SSD,FALSE,0,None,50,Broadband,"Linux,Windows,macOS",Theta,THETA,30,120,Medium,Video streaming edge caching,TRUE,2025-01-01
Storj,Storage,Storage Node,2,Any,4,8,500,Any,FALSE,0,None,25,Broadband,"Linux,Windows,macOS",Ethereum,STORJ,15,60,Low,Decentralized cloud storage node,TRUE,2025-01-01
EOF
    
    print_success "Sample data created (replace with your depin_specifications_final.csv)"
    print_warning "To use your actual data, copy depin_specifications_final.csv to the project root before running setup"
}

# Install Go dependencies
install_dependencies() {
    print_status "Installing Go dependencies..."
    
    # Initialize Go module if it doesn't exist
    if [ ! -f go.mod ]; then
        go mod init depin-compatibility-api
    fi
    
    # Install main dependencies
    go get github.com/gin-gonic/gin@latest
    go get golang.org/x/time@latest
    
    # Tidy up dependencies
    go mod tidy
    
    print_success "Dependencies installed"
}

# Create Makefile
create_makefile() {
    print_status "Creating Makefile..."
    
    cat > Makefile << 'EOF'
.PHONY: build run test clean dev docker docker-compose help

# Default target
help:
	@echo "Available commands:"
	@echo "  make run          - Run the application"
	@echo "  make dev          - Run with hot reload (requires air)"
	@echo "  make build        - Build the application"
	@echo "  make build-prod   - Build optimized for production"
	@echo "  make test         - Run tests"
	@echo "  make test-coverage - Run tests with coverage"
	@echo "  make clean        - Clean build artifacts"
	@echo "  make docker       - Build Docker image"
	@echo "  make docker-run   - Run Docker container"
	@echo "  make compose-up   - Start with Docker Compose"
	@echo "  make compose-down - Stop Docker Compose"

# Application name and version
APP_NAME := depin-api
VERSION := 1.0.0
BUILD_TIME := $(shell date +%Y-%m-%d_%H:%M:%S)
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Build flags
LDFLAGS := -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT)"

# Run the application
run:
	go run main.go

# Run with hot reload (requires air: go install github.com/cosmtrek/air@latest)
dev:
	@if command -v air >/dev/null 2>&1; then \
		air; \
	else \
		echo "Hot reload requires 'air'. Install with: go install github.com/cosmtrek/air@latest"; \
		make run; \
	fi

# Build the application
build:
	go build $(LDFLAGS) -o bin/$(APP_NAME) main.go

# Build optimized for production
build-prod:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo $(LDFLAGS) -o bin/$(APP_NAME) main.go

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# Docker commands
docker:
	docker build -t $(APP_NAME):$(VERSION) .
	docker tag $(APP_NAME):$(VERSION) $(APP_NAME):latest

docker-run:
	docker run -d --name $(APP_NAME) -p 8080:8080 $(APP_NAME):latest

# Docker Compose commands
compose-up:
	docker-compose up -d

compose-down:
	docker-compose down

# Install development tools
install-tools:
	go install github.com/cosmtrek/air@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Lint code
lint:
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found. Install with: make install-tools"; \
	fi

# Format code
format:
	go fmt ./...
	@if command -v goimports >/dev/null 2>&1; then \
		goimports -w .; \
	fi
EOF
    
    print_success "Makefile created"
}

# Create Dockerfile
create_dockerfile() {
    print_status "Creating Dockerfile..."
    
    cat > Dockerfile << 'EOF'
# Build stage
FROM golang:1.21-alpine AS builder

# Install git and ca-certificates
RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates

# Create appuser
ENV USER=appuser
ENV UID=10001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download
RUN go mod verify

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo -o app main.go

# Final stage
FROM scratch

# Import from builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

# Copy our static executable and data
COPY --from=builder /app/app /app
COPY --from=builder /app/data /data

# Use an unprivileged user
USER appuser:appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD ["/app", "-health-check"] || exit 1

# Run the binary
ENTRYPOINT ["/app"]
EOF
    
    print_success "Dockerfile created"
}

# Create Docker Compose file
create_docker_compose() {
    print_status "Creating Docker Compose configuration..."
    
    cat > docker-compose.yml << 'EOF'
version: '3.8'

services:
  depin-api:
    build: .
    container_name: depin-compatibility-api
    ports:
      - "8080:8080"
    environment:
      - PORT=8080
      - HOST=0.0.0.0
      - DATA_PATH=/data/depin_specs.csv
      - LOG_LEVEL=info
      - GIN_MODE=release
    volumes:
      - ./data:/data:ro
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "/app", "-health-check"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s
    networks:
      - depin-network

  # Optional: Nginx reverse proxy
  nginx:
    image: nginx:alpine
    container_name: depin-nginx
    ports:
      - "80:80"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf:ro
    depends_on:
      - depin-api
    restart: unless-stopped
    networks:
      - depin-network
    profiles:
      - with-proxy

networks:
  depin-network:
    driver: bridge
EOF
    
    print_success "Docker Compose configuration created"
}

# Create example files
create_examples() {
    print_status "Creating usage examples..."
    
    # cURL examples
    cat > examples/curl/examples.sh << 'EOF'
#!/bin/bash

# DePIN Compatibility API - cURL Examples

API_URL="http://localhost:8080"

echo "🚀 DePIN Compatibility API - cURL Examples"
echo "=========================================="

# Health check
echo "1. Health Check:"
curl -s "$API_URL/api/v1/health" | jq .
echo -e "\n"

# High-end gaming system
echo "2. High-End Gaming System:"
curl -s -X POST "$API_URL/api/v1/predict" \
  -H "Content-Type: application/json" \
  -d '{
    "system": {
      "cpu_cores": 12,
      "ram_gb": 32,
      "storage_gb": 1000,
      "has_ssd": true,
      "has_gpu": true,
      "gpu_vram_gb": 16,
      "network_mbps": 500,
      "os": "Windows"
    }
  }' | jq .
echo -e "\n"

# Budget laptop
echo "3. Budget Laptop:"
curl -s -X POST "$API_URL/api/v1/predict" \
  -H "Content-Type: application/json" \
  -d '{
    "system": {
      "cpu_cores": 4,
      "ram_gb": 8,
      "storage_gb": 256,
      "has_ssd": true,
      "has_gpu": false,
      "gpu_vram_gb": 0,
      "network_mbps": 50,
      "os": "macOS"
    }
  }' | jq .
echo -e "\n"

# List all projects
echo "4. All DePIN Projects:"
curl -s "$API_URL/api/v1/projects" | jq '.summary'
echo -e "\n"

echo "✅ Examples completed!"
EOF

    chmod +x examples/curl/examples.sh
    
    # Python client
    cat > examples/python/client.py << 'EOF'
#!/usr/bin/env python3
"""
DePIN Compatibility API - Python Client Example
"""

import requests
import json
from typing import Dict, Any

class DePINClient:
    def __init__(self, base_url: str = "http://localhost:8080"):
        self.base_url = base_url.rstrip('/')
        self.session = requests.Session()
        self.session.headers.update({'Content-Type': 'application/json'})
    
    def health_check(self) -> Dict[str, Any]:
        """Check API health status"""
        response = self.session.get(f"{self.base_url}/api/v1/health")
        response.raise_for_status()
        return response.json()
    
    def predict_compatibility(self, system_spec: Dict[str, Any]) -> Dict[str, Any]:
        """Predict DePIN compatibility for a system"""
        payload = {"system": system_spec}
        response = self.session.post(f"{self.base_url}/api/v1/predict", json=payload)
        response.raise_for_status()
        return response.json()
    
    def list_projects(self) -> Dict[str, Any]:
        """Get all DePIN projects"""
        response = self.session.get(f"{self.base_url}/api/v1/projects")
        response.raise_for_status()
        return response.json()

def main():
    """Example usage"""
    client = DePINClient()
    
    # Health check
    print("🏥 Health Check:")
    health = client.health_check()
    print(f"Status: {health['status']}")
    print(f"Projects Loaded: {health['projects_loaded']}")
    print()
    
    # Example system specifications
    systems = [
        {
            "name": "High-End Gaming PC",
            "spec": {
                "cpu_cores": 12,
                "ram_gb": 32,
                "storage_gb": 1000,
                "has_ssd": True,
                "has_gpu": True,
                "gpu_vram_gb": 16,
                "network_mbps": 500,
                "os": "Windows"
            }
        },
        {
            "name": "Budget Laptop",
            "spec": {
                "cpu_cores": 4,
                "ram_gb": 8,
                "storage_gb": 256,
                "has_ssd": True,
                "has_gpu": False,
                "gpu_vram_gb": 0,
                "network_mbps": 50,
                "os": "macOS"
            }
        }
    ]
    
    # Test each system
    for system in systems:
        print(f"🖥️ Testing {system['name']}:")
        result = client.predict_compatibility(system['spec'])
        
        summary = result['summary']
        print(f"  System Rating: {summary['system_rating']}")
        print(f"  Compatible Projects: {summary['compatible_count']}/{summary['total_projects']}")
        print(f"  Compatibility Rate: {summary['compatibility_rate']:.1f}%")
        
        if result['compatible_projects']:
            best_project = result['compatible_projects'][0]
            print(f"  Best Match: {best_project['name']} ({best_project['performance_rating']})")
        
        if result['recommendations']:
            print(f"  Recommendation: {result['recommendations'][0]}")
        
        print()

if __name__ == "__main__":
    main()
EOF

    # JavaScript client
    cat > examples/javascript/client.js << 'EOF'
#!/usr/bin/env node
/**
 * DePIN Compatibility API - Node.js Client Example
 */

const https = require('https');
const http = require('http');
const { URL } = require('url');

class DePINClient {
    constructor(baseUrl = 'http://localhost:8080') {
        this.baseUrl = baseUrl.replace(/\/$/, '');
    }

    async request(method, path, data = null) {
        return new Promise((resolve, reject) => {
            const url = new URL(path, this.baseUrl);
            const client = url.protocol === 'https:' ? https : http;
            
            const options = {
                method,
                headers: {
                    'Content-Type': 'application/json',
                    'User-Agent': 'DePIN-Client/1.0.0'
                }
            };

            if (data) {
                const jsonData = JSON.stringify(data);
                options.headers['Content-Length'] = Buffer.byteLength(jsonData);
            }

            const req = client.request(url, options, (res) => {
                let body = '';
                res.on('data', chunk => body += chunk);
                res.on('end', () => {
                    try {
                        const parsed = JSON.parse(body);
                        if (res.statusCode >= 200 && res.statusCode < 300) {
                            resolve(parsed);
                        } else {
                            reject(new Error(`HTTP ${res.statusCode}: ${parsed.error || body}`));
                        }
                    } catch (err) {
                        reject(new Error(`Parse error: ${err.message}`));
                    }
                });
            });

            req.on('error', reject);

            if (data) {
                req.write(JSON.stringify(data));
            }

            req.end();
        });
    }

    async healthCheck() {
        return this.request('GET', '/api/v1/health');
    }

    async predictCompatibility(systemSpec) {
        return this.request('POST', '/api/v1/predict', { system: systemSpec });
    }

    async listProjects() {
        return this.request('GET', '/api/v1/projects');
    }
}

async function main() {
    const client = new DePINClient();

    try {
        // Health check
        console.log('🏥 Health Check:');
        const health = await client.healthCheck();
        console.log(`Status: ${health.status}`);
        console.log(`Projects Loaded: ${health.projects_loaded}`);
        console.log();

        // Example systems
        const systems = [
            {
                name: 'High-End Workstation',
                spec: {
                    cpu_cores: 16,
                    ram_gb: 64,
                    storage_gb: 2000,
                    has_ssd: true,
                    has_gpu: true,
                    gpu_vram_gb: 24,
                    network_mbps: 1000,
                    os: 'Linux'
                }
            },
            {
                name: 'Gaming Laptop',
                spec: {
                    cpu_cores: 8,
                    ram_gb: 16,
                    storage_gb: 512,
                    has_ssd: true,
                    has_gpu: true,
                    gpu_vram_gb: 8,
                    network_mbps: 100,
                    os: 'Windows'
                }
            }
        ];

        // Test each system
        for (const system of systems) {
            console.log(`🖥️ Testing ${system.name}:`);
            const result = await client.predictCompatibility(system.spec);

            const summary = result.summary;
            console.log(`  System Rating: ${summary.system_rating}`);
            console.log(`  Compatible: ${summary.compatible_count}/${summary.total_projects} projects`);
            console.log(`  Compatibility Rate: ${summary.compatibility_rate.toFixed(1)}%`);

            if (result.compatible_projects.length > 0) {
                const best = result.compatible_projects[0];
                console.log(`  Best Match: ${best.name} (${best.performance_rating})`);
            }

            if (result.recommendations.length > 0) {
                console.log(`  Recommendation: ${result.recommendations[0]}`);
            }

            console.log();
        }

    } catch (error) {
        console.error('❌ Error:', error.message);
    }
}

if (require.main === module) {
    main();
}

module.exports = DePINClient;
EOF

    print_success "Usage examples created"
}

# Create Git configuration
create_git_config() {
    print_status "Creating Git configuration..."
    
    # .gitignore
    cat > .gitignore << 'EOF'
# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib
bin/
dist/

# Test binary, built with `go test -c`
*.test

# Output of the go coverage tool
*.out
coverage.html

# Go workspace file
go.work

# IDE and editor files
.vscode/
.idea/
*.swp
*.swo
*~

# OS generated files
.DS_Store
.DS_Store?
._*
.Spotlight-V100
.Trashes
ehthumbs.db
Thumbs.db

# Log files
*.log
logs/

# Environment files
.env
.env.local
.env.production

# Database files
*.db
*.sqlite

# Temporary files
tmp/
temp/

# Docker
.dockerignore

# Air (hot reload)
tmp/
EOF
    
    print_success "Git configuration created"
}

# Create documentation
create_documentation() {
    print_status "Creating documentation..."
    
    # API Documentation
    cat > docs/API.md << 'EOF'
# DePIN Compatibility API Documentation

## Overview

The DePIN Compatibility API analyzes system specifications and predicts compatibility with various DePIN (Decentralized Physical Infrastructure Network) projects.

## Base URL

```
http://localhost:8080/api/v1
```

## Endpoints

### POST /predict

Predicts DePIN compatibility for a given system.

**Request Body:**
```json
{
  "system": {
    "cpu_cores": 8,
    "ram_gb": 16,
    "storage_gb": 512,
    "has_ssd": true,
    "has_gpu": true,
    "gpu_vram_gb": 8,
    "network_mbps": 100,
    "os": "Windows"
  }
}
```

**Response:**
```json
{
  "compatible_projects": [...],
  "incompatible_projects": [...],
  "summary": {
    "total_projects": 8,
    "compatible_count": 6,
    "compatibility_rate": 75.0,
    "system_rating": "High-End"
  },
  "recommendations": [...],
  "generated_at": "2024-01-15T10:30:00Z"
}
```

### GET /health

Health check endpoint.

### GET /projects

Lists all available DePIN projects.

### GET /docs

API documentation (this page).

### GET /metrics

Service metrics and statistics.

## System Specifications

| Field | Type | Description | Range |
|-------|------|-------------|-------|
| `cpu_cores` | int | Number of CPU cores | 1-64 |
| `ram_gb` | int | RAM in GB | 1-128 |
| `storage_gb` | int | Storage in GB | 32-8192 |
| `has_ssd` | bool | SSD storage | true/false |
| `has_gpu` | bool | Dedicated GPU | true/false |
| `gpu_vram_gb` | int | GPU VRAM in GB | 0-48 |
| `network_mbps` | int | Network speed in Mbps | 1-10000 |
| `os` | string | Operating system | Windows/Linux/macOS |

## Compatibility Scores

- **Excellent (0.9-1.0)**: System exceeds requirements
- **Good (0.7-0.89)**: System meets requirements well
- **Fair (0.5-0.69)**: System meets minimum requirements
- **Poor (0.0-0.49)**: System has limitations

## Error Responses

```json
{
  "error": "Error type",
  "message": "Detailed error message",
  "code": 400,
  "timestamp": "2024-01-15T10:30:00Z"
}
```

## Rate Limiting

- 10 requests per second per IP
- Burst up to 20 requests
- Headers: `X-RateLimit-*`

## Examples

See the `examples/` directory for complete client implementations in:
- cURL
- Python
- Node.js
EOF
    
    print_success "Documentation created"
}

# Build application
build_application() {
    print_status "Building application..."
    
    if command_exists make; then
        make build
    else
        go build -o bin/depin-api main.go
    fi
    
    print_success "Application built successfully"
}

# Run tests
run_tests() {
    print_status "Running basic validation..."
    
    # Test if the application compiles
    if go build -o /tmp/depin-api-test main.go; then
        rm -f /tmp/depin-api-test
        print_success "Compilation test passed"
    else
        print_error "Compilation test failed"
        return 1
    fi
    
    # Test data loading
    if [ -f "data/depin_specs.csv" ]; then
        print_success "Data file validation passed"
    else
        print_error "Data file missing"
        return 1
    fi
}

# Start application
start_application() {
    print_status "Starting application..."
    
    # Start in background
    if command_exists make; then
        nohup make run > app.log 2>&1 &
    else
        nohup go run main.go > app.log 2>&1 &
    fi
    
    APP_PID=$!
    echo $APP_PID > app.pid
    
    # Wait for startup
    print_status "Waiting for application to start..."
    for i in {1..30}; do
        if curl -s http://localhost:$DEFAULT_PORT/api/v1/health > /dev/null 2>&1; then
            print_success "Application started successfully on port $DEFAULT_PORT"
            return 0
        fi
        sleep 1
    done
    
    print_error "Application failed to start"
    return 1
}

# Test API
test_api() {
    print_status "Testing API endpoints..."
    
    # Test health endpoint
    if curl -s http://localhost:$DEFAULT_PORT/api/v1/health | grep -q "healthy"; then
        print_success "Health endpoint working"
    else
        print_error "Health endpoint failed"
        return 1
    fi
    
    # Test prediction endpoint
    TEST_PAYLOAD='{"system":{"cpu_cores":8,"ram_gb":16,"storage_gb":512,"has_ssd":true,"has_gpu":true,"gpu_vram_gb":8,"network_mbps":100,"os":"Linux"}}'
    
    if curl -s -X POST http://localhost:$DEFAULT_PORT/api/v1/predict \
       -H "Content-Type: application/json" \
       -d "$TEST_PAYLOAD" | grep -q "compatible_projects"; then
        print_success "Prediction endpoint working"
    else
        print_error "Prediction endpoint failed"
        return 1
    fi
    
    print_success "API tests passed"
}

# Cleanup function
cleanup() {
    print_status "Cleaning up..."
    
    if [ -f app.pid ]; then
        PID=$(cat app.pid)
        if kill -0 $PID 2>/dev/null; then
            kill $PID
            print_status "Stopped application (PID: $PID)"
        fi
        rm -f app.pid
    fi
    
    rm -f app.log
}

# Main setup function
main() {
    print_banner
    
    # Set trap for cleanup
    trap cleanup EXIT
    
    print_status "Starting DePIN Compatibility API setup..."
    echo
    
    # Run setup steps
    check_prerequisites
    create_project_structure
    create_sample_data
    install_dependencies
    create_makefile
    create_dockerfile
    create_docker_compose
    create_examples
    create_git_config
    create_documentation
    build_application
    run_tests
    
    echo
    print_success "Setup completed successfully!"
    echo
    
    # Ask user if they want to start the application
    echo -e "${BLUE}Would you like to start the API now? (y/n)${NC}"
    read -r START_NOW
    
    if [[ $START_NOW =~ ^[Yy]$ ]]; then
        echo
        start_application
        
        if [ $? -eq 0 ]; then
            test_api
            echo
            print_success "🎉 DePIN Compatibility API is now running!"
            echo
            echo "📋 Quick Start:"
            echo "  • Health Check: curl http://localhost:$DEFAULT_PORT/api/v1/health"
            echo "  • API Docs: curl http://localhost:$DEFAULT_PORT/api/v1/docs"
            echo "  • Test Examples: cd examples/curl && ./examples.sh"
            echo
            echo "📁 Files created:"
            echo "  • Application: bin/depin-api"
            echo "  • Data: data/depin_specs.csv"
            echo "  • Examples: examples/"
            echo "  • Documentation: docs/"
            echo
            echo "🐳 Docker commands:"
            echo "  • Build: make docker"
            echo "  • Run: make docker-run"
            echo "  • Compose: make compose-up"
            echo
            echo "🛠️ Development:"
            echo "  • Build: make build"
            echo "  • Test: make test"
            echo "  • Hot reload: make dev"
            echo
            print_warning "The application is running in the background."
            print_warning "Check app.log for logs, or use 'kill \$(cat app.pid)' to stop."
        fi
    else
        echo
        print_status "Setup complete! You can start the API later with:"
        echo "  make run    # or go run main.go"
        echo
        print_status "Next steps:"
        echo "  1. Review the configuration in the generated files"
        echo "  2. Customize the DePIN data in data/depin_specs.csv"
        echo "  3. Check the examples in examples/"
        echo "  4. Read the documentation in docs/"
    fi
}

# Run main function
main "$@"
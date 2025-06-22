# DePIN Compatibility API

> **Simple, fast, and reliable DePIN compatibility prediction for consumer systems**

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Docker](https://img.shields.io/badge/Docker-Ready-green.svg)](https://docker.com)

A lightweight REST API that predicts DePIN (Decentralized Physical Infrastructure Network) compatibility based on consumer system specifications. Built with Go for high performance and easy deployment.

## ✨ Features

- 🚀 **Fast & Lightweight** - Pure Go implementation, no heavy ML dependencies
- 🎯 **Consumer Focused** - Optimized for desktop, laptop, and gaming systems
- 📊 **Smart Scoring** - Advanced rule-based compatibility algorithm
- 🐳 **Docker Ready** - Complete containerization with Docker Compose
- 🔧 **Easy Setup** - One-command installation and deployment
- 📖 **Well Documented** - Comprehensive API documentation and examples
- 🧪 **Tested** - Built-in test suite with coverage reports

## 🚀 Quick Start

### Option 1: One-Command Setup
```bash
curl -sSL https://raw.githubusercontent.com/yourusername/depin-compatibility-api/main/scripts/setup.sh | bash
```

### Option 2: Manual Setup
```bash
# Clone the repository
git clone https://github.com/yourusername/depin-compatibility-api.git
cd depin-compatibility-api

# (Optional) Copy your actual DePIN data
# cp /path/to/depin_specifications_final.csv ./

# Run setup script
chmod +x scripts/setup.sh
./scripts/setup.sh

# Start the API
make run
```

### Option 3: Docker
```bash
# Using Docker Compose (recommended)
docker-compose up -d

# Or build and run manually
docker build -t depin-api .
docker run -p 8080:8080 depin-api
```

## 📋 API Usage

### Check System Compatibility
```bash
curl -X POST http://localhost:8080/api/v1/predict \
  -H "Content-Type: application/json" \
  -d '{
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
  }'
```

### Response
```json
{
  "compatible_projects": [
    {
      "name": "Filecoin Light Node",
      "compatibility_score": 0.95,
      "performance_rating": "Excellent",
      "estimated_cost": "$10-50/month",
      "missing_requirements": [],
      "recommended_upgrades": []
    }
  ],
  "summary": {
    "total_projects": 8,
    "compatible_count": 6,
    "compatibility_rate": 75.0
  },
  "recommendations": [
    "Your system has excellent compatibility with most DePIN projects!"
  ]
}
```

## 🎯 System Requirements Analysis

### Supported Systems
- **Desktop Computers** - Windows, Linux, macOS
- **Gaming Systems** - High-performance desktops with GPUs
- **Laptops** - Consumer laptops and mobile workstations
- **Workstations** - Content creation and development systems

### Hardware Ranges
| Component | Entry Level | Mid-Range | High-End |
|-----------|-------------|-----------|----------|
| **CPU** | 2-4 cores | 6-8 cores | 10-12 cores |
| **RAM** | 4-8 GB | 16 GB | 32 GB |
| **Storage** | 256 GB | 512 GB | 1 TB |
| **GPU** | Integrated | 6-8 GB | 12-16 GB |
| **Network** | 25 Mbps | 100 Mbps | 500 Mbps |

## 🛠️ Development

### Prerequisites
- Go 1.21+
- Make (optional, for build automation)
- Docker & Docker Compose (optional)

### Local Development
```bash
# Install dependencies
go mod download

# Run tests
make test

# Run with hot reload
make dev

# Build binary
make build

# Run linting
make lint
```

### Project Structure
```
├── main.go              # Application entry point
├── internal/            # Private application code
│   ├── api/            # HTTP layer (handlers, middleware)
│   ├── models/         # Data structures
│   ├── service/        # Business logic
│   └── data/           # Data access layer
├── data/               # CSV data files
├── scripts/            # Utility scripts
├── examples/           # Usage examples
└── docs/               # Documentation
```

## 📚 Examples

### Python Client
```python
import requests

response = requests.post('http://localhost:8080/api/v1/predict', json={
    'system': {
        'cpu_cores': 8,
        'ram_gb': 16,
        'storage_gb': 512,
        'has_ssd': True,
        'has_gpu': True,
        'gpu_vram_gb': 8,
        'network_mbps': 100,
        'os': 'Linux'
    }
})

result = response.json()
print(f"Compatible projects: {result['summary']['compatible_count']}")
```

### JavaScript Client
```javascript
const response = await fetch('http://localhost:8080/api/v1/predict', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    system: {
      cpu_cores: 8,
      ram_gb: 16,
      storage_gb: 512,
      has_ssd: true,
      has_gpu: true,
      gpu_vram_gb: 8,
      network_mbps: 100,
      os: 'macOS'
    }
  })
});

const result = await response.json();
console.log(`Compatible: ${result.summary.compatible_count} projects`);
```

## 🔧 Configuration

Environment variables:
```bash
# Server configuration
PORT=8080                    # Server port (default: 8080)
HOST=localhost              # Server host (default: localhost)

# Data configuration
DATA_PATH=./data            # Path to data files
CSV_FILE=depin_specs.csv    # DePIN specifications file

# Logging
LOG_LEVEL=info              # Log level (debug, info, warn, error)
LOG_FORMAT=json             # Log format (json, text)
```

### Using Your Own DePIN Data

To use your actual DePIN project data:

1. **Place your CSV file** in the project root as `depin_specifications_final.csv`
2. **Run the setup script** - it will automatically detect and use your file
3. **Or manually copy** your CSV to `data/depin_specs.csv`

The API will load your actual DePIN projects with their real specifications, requirements, and cost estimates.

## 📊 API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/v1/predict` | Predict DePIN compatibility |
| `GET` | `/api/v1/health` | Health check |
| `GET` | `/api/v1/projects` | List all DePIN projects |
| `GET` | `/api/v1/metrics` | Prometheus metrics |

For detailed API documentation, see [docs/API.md](docs/API.md).

## 🧪 Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run integration tests
make test-integration

# Benchmark tests
make benchmark
```

## 🚀 Deployment

### Production Deployment
```bash
# Build optimized binary
make build-prod

# Deploy with Docker
docker-compose -f docker-compose.prod.yml up -d

# Deploy to cloud (example)
make deploy-aws
```

### Health Monitoring
The API includes built-in health checks and Prometheus metrics:
```bash
# Health check
curl http://localhost:8080/api/v1/health

# Metrics
curl http://localhost:8080/api/v1/metrics
```

## 🤝 Contributing

We welcome contributions! Please see [docs/CONTRIBUTING.md](docs/CONTRIBUTING.md) for guidelines.

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

---
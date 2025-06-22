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

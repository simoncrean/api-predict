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

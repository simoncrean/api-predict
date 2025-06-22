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

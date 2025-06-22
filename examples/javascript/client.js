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

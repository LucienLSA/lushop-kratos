# Jaeger Deployment

## Overview
Jaeger distributed tracing system for monitoring and troubleshooting microservices.

## Components
- **Jaeger All-in-One**: Complete Jaeger installation with UI, collector, and agent
- **Memory Storage**: In-memory span storage (suitable for development)
- **Multiple Protocols**: Supports Jaeger, Zipkin, and OTLP protocols

## Access
- **Jaeger UI**: http://[NODE-IP]:31686
- **Collector (Jaeger)**: Port 14268
- **Collector (Zipkin)**: Port 9411
- **Collector (OTLP)**: Port 4317

## Features
- **Distributed Tracing**: Track requests across microservices
- **Service Dependencies**: Visualize service call graphs
- **Performance Analysis**: Identify bottlenecks and latency issues
- **Error Tracking**: Monitor and debug failed requests

## Usage
```bash
# Deploy
kubectl apply -k .

# Access UI
kubectl port-forward svc/jaeger 16686:16686

# Test collector
curl -X POST http://localhost:14268/api/traces \
  -H "Content-Type: application/x-protobuf" \
  -d @trace-data.pb
```

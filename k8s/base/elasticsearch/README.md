# Elasticsearch + Kibana Deployment

## Overview
Elasticsearch search and analytics engine with Kibana visualization dashboard for Lushop log aggregation and analysis.

## Components
- **Elasticsearch**: Search and analytics engine
- **Kibana**: Data visualization and exploration
- **Init Container**: System tuning for Elasticsearch
- **Persistent Storage**: Data persistence across restarts

## Access
- **Elasticsearch API**: http://[NODE-IP]:9200
- **Kibana UI**: http://[NODE-IP]:30561

## Features
- **Full-text Search**: Advanced search capabilities
- **Real-time Analytics**: Live data analysis and aggregation
- **Kibana Dashboards**: Interactive data visualizations
- **Log Aggregation**: Centralized logging solution
- **RESTful API**: HTTP-based API for data operations

## Usage
```bash
# Deploy
kubectl apply -k .

# Access Kibana
kubectl port-forward svc/kibana 5601:5601

# Test Elasticsearch
curl http://localhost:9200/_cluster/health
```

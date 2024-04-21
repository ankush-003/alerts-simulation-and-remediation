# Alerts Simulation & Remediation 🔔
![Workflow](https://github.com/ankush-003/alerts-simulation-and-remediation/actions/workflows/main.yml/badge.svg)
![Simulator Docker Pulls](https://img.shields.io/docker/pulls/ankush003/simulator)
![GitHub commit activity](https://img.shields.io/github/commit-activity/t/ankush-003/alerts-simulation-and-remediation)
![Vercel Deploy](https://deploy-badge.vercel.app/vercel/alerts-simulation-and-remediation)

`Alert Simulation and Remediation` is an advanced monitoring and alerting system designed to help manage alerts from deployments effectively. This project aims to provide a comprehensive solution for simulating various system environments, evaluating alerts, providing remediation recommendations, and delivering real-time notifications and insights.

## Features ✨
### Simulation Environment 🌲
- `Simulator`: Simulates various system environments, such as high CPU load, network load, low memory availability, and high disk usage, by creating multiple goroutines, sending HTTP requests, allocating memory, and writing files.

### Alert Management 📢
- `Rule Engine`: Evaluates alerts based on predefined rules and provides remediation recommendations.
- `Prometheus and Grafana Stack`: Fetches and visualizes system metrics using Prometheus and Grafana.
- `Kafka Integration`: Utilizes Kafka for communication between the rule engine and simulator.

### Notification and Insights 📣
- `Mail Server`: Sends email notifications for critical alerts.
- `Real-time Notifications`: Leverages Redis Streams and Server-Sent Events (SSE) to deliver real-time alert notifications to the frontend dashboard.
- `ASMR QueryBot`: Implements a chatbot powered by Large Language Models (LLMs) to provide interactive insights and answer user queries related to alerts and system performance.
- `MongoDB Vector Search`: Stores alert data as vectors using MongoDB, enabling efficient searching and querying with LLMs using the LlamaIndex framework.

### Deployment and Scalability 🌍
- `Dockerized`: The entire project is containerized using Docker, ensuring consistent deployment across different environments.

## Architecture 🗺️
![image](https://github.com/ankush-003/alerts-simulation-and-remediation/assets/94037471/c652d953-9bcb-4dac-baa8-438c6fffb7ac)

## Workflow 🖥️
![image](https://github.com/ankush-003/alerts-simulation-and-remediation/assets/94037471/583eea7c-77df-4bc7-8db4-0e5a963c0ea5)

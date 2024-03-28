import json
import random

def generate_metrics():
  """Generates random system metrics and returns them as a JSON object."""

  metrics = {
    "hardware": {
      "cpu": {
        "utilization": random.randint(0, 100),  # CPU utilization (%)
        "temperature": random.randint(30, 80),  # CPU temperature (°C)
        "cores": random.randint(1, 16),        # Number of CPU cores
        "wait_time": random.randint(0, 20),     # CPU wait time (%)
      },
      "memory": {
        "usage": random.randint(0, 99),        # Memory usage (%)
        "free": None,                          # Free memory (MB or GB)
        "swap_usage": random.randint(0, 20),     # Swap usage (%)
      },
      "disk": {
        "read_write": random.randint(10, 100),  # Disk read/write throughput (MB/s)
        "utilization": random.randint(0, 90),    # Disk utilization (%)
        "available": None,                     # Available disk space (MB or GB)
      },
      "network": {
        "traffic_in": random.randint(1, 20),      # Network traffic in (MB/s)
        "traffic_out": random.randint(1, 20),     # Network traffic out (MB/s)
        "packets_in": random.randint(1000, 5000),  # Packets in (pps)
        "packets_out": random.randint(1000, 5000), # Packets out (pps)
        "latency": random.randint(1, 50),        # Network latency (ms)
      },
      "other": {
        "fan_speed": random.randint(800, 1500),  # Fan speeds (RPM)
        "uptime": random.randint(3600, 86400),   # System uptime (seconds)
        "voltage": random.uniform(11.5, 12.5),   # Voltage levels (V)
      }
    },
    "software": {
      "processes": {
        "running": random.randint(10, 50),       # Number of running processes
      },
      "applications": {
        "response_time": random.randint(50, 200), # Response time (ms)
        "active_users": random.randint(0, 100),    # Number of active users
        "error_rate": random.uniform(0, 0.05),     # Error rate (%)
      },
      "os": {
        "load_average": random.uniform(0, 5),     # System load average
        "disk_queue_length": random.randint(0, 10), # Disk queue length
        "kernel_version": f"6.{random.randint(0, 7)}.{random.randint(0, 10)}"  # Kernel version (example format)
      }
    }
  }

  # Calculate free memory based on usage percentage (assuming total memory is 16GB)
  metrics["hardware"]["memory"]["free"] = int(16000 * (1 - metrics["hardware"]["memory"]["usage"] / 100))

  # Calculate available disk space based on a random utilization (assuming total disk is 1TB)
  metrics["hardware"]["disk"]["available"] = int(1000000 * (1 - metrics["hardware"]["disk"]["utilization"] / 100))

  return json.dumps(metrics, indent=4)

if __name__ == "__main__":
  metrics_json = generate_metrics()
  print(metrics_json)

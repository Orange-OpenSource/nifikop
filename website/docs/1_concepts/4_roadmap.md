---
id: 4_roadmap
title: Roadmap
sidebar_label: Roadmap
---

## Available

### NiFi cluster installation

|                       |           |
| --------------------- | --------- |
| Status                | Done      |
| Priority              | High      |
| Targeted Start date   | Jan 2020  |

### Graceful NiFi Cluster Scaling

|                       |           |
| --------------------- | --------- |
| Status                | Done      |
| Priority              | High      |
| Targeted Start date   | Jan 2020  |

Apache NiFi is a good candidate to create an operator, because everything is made to orchestrate it through REST Api calls. With this comes automation of actions such as scaling, following all required steps : https://nifi.apache.org/docs/nifi-docs/html/administration-guide.html#decommission-nodes.

### Communication via SSL

|                       |          |
| --------------------- | -------- |
| Status                | Done     |
| Priority              | High     |
| Targeted Start date   | May 2020 |


The operator fully automates NiFi's SSL support.
The operator can provision the required secrets and certificates for you, or you can provide your own.

### Dataflow lifecycle management via CRD

|                       |           |
| --------------------- | --------- |
| Status                | Done      |
| Priority              | High      |
| Targeted Start date   | Aug 2020 |

## Backlog

### Monitoring via Prometheus

|                       |          |
| --------------------- | -------- |
| Status                | To Do    |
| Priority              | High     |
| Targeted Start date   | Oct 2020 |

The NiFi operator exposes NiFi JMX metrics to Prometheus.

### Authentification management

|                       |       |
| --------------------- | ----- |
| Status                | To Do |
| Priority              | High  |
| Targeted Start date   | -     |


### Reacting on Alerts

|                       |       |
| --------------------- | ----- |
| Status                | To Do |
| Priority              | Low   |
| Targeted Start date   | -     |

The NiFi Operator acts as a **Prometheus Alert Manager**. It receives alerts defined in Prometheus, and creates actions based on Prometheus alert annotations.

Currently, there are three actions expected :
- upscale cluster (add a new Node)
- downscale cluster (remove a Node)
- add additional disk to a Node

### Seamless Istio mesh support

|                       |       |
| --------------------- | ----- |
| Status                | To Do |
| Priority              | Low   |
| Targeted Start date   | -     |

- Operator allows to use ClusterIP services instead of Headless, which still works better in case of Service meshes.
- To avoid too early nifi initialization, which might lead to unready sidecar container. The operator will use a small script to
mitigate this behaviour. All NiFi image can be used the only one requirement is an available **wget** command.
- To access a NiFi cluster which runs inside the mesh. Operator will supports creating Istio ingress gateways.
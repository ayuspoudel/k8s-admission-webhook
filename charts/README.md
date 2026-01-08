# Helm Help

This directory contains Helm charts for deploying various applications and services. Each subdirectory corresponds to a specific Helm chart.

To see the schema for values run:

```bash
helm plugin install https://github.com/krak3n/helm-values-schema-generator
helm values-schema .

```
This will generate values.schema.json against which your values file can be tested and so on.
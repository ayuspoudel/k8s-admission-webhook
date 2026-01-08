
### Frequently Used Pipelines
---
[![Helm Chart Tests](https://github.com/ayuspoudel/sentinel-sre/actions/workflows/helm-test.yaml/badge.svg)](https://github.com/ayuspoudel/sentinel-sre/actions/workflows/helm-test.yaml)
[![Helm Chart Release](https://github.com/ayuspoudel/sentinel-sre/actions/workflows/helm-release.yaml/badge.svg)](https://github.com/ayuspoudel/sentinel-sre/actions/workflows/helm-release.yaml)
[![Helm Snapshot Update (Manual)](https://github.com/ayuspoudel/sentinel-sre/actions/workflows/helm-snapshot-release.yaml/badge.svg)](https://github.com/ayuspoudel/sentinel-sre/actions/workflows/helm-snapshot-release.yaml)
[![Helm Deploy Gh Pages](https://github.com/ayuspoudel/sentinel-sre/actions/workflows/pages/pages-build-deployment/badge.svg)](https://github.com/ayuspoudel/sentinel-sre/actions/workflows/pages/pages-build-deployment)
[![Sentinel Agent CI/CD](https://github.com/ayuspoudel/sentinel-sre/actions/workflows/agent-cicd.yaml/badge.svg)](https://github.com/ayuspoudel/sentinel-sre/actions/workflows/agent-cicd.yaml)
[![Helm Repo E2E Test for Sentinel Agent](https://github.com/ayuspoudel/sentinel-sre/actions/workflows/helm-repo-e2e-sentinel-agent.yaml/badge.svg)](https://github.com/ayuspoudel/sentinel-sre/actions/workflows/helm-repo-e2e-sentinel-agent.yaml)
[![Helm Chart Version Bump](https://github.com/ayuspoudel/sentinel-sre/actions/workflows/helm-version-update.yaml/badge.svg)](https://github.com/ayuspoudel/sentinel-sre/actions/workflows/helm-version-update.yaml)

---

# K8s Admission Webhook (Go)
---
This repository is a minimalistic but complete implementation of a Kubernetes validating admission webhook written in Go. It is packaged both as raw Kubernetes manifests under the `deploy/` directory and as a Helm chart under `charts/`.

This is not a policy engine, nor a wrapper around an existing controller. The intent of this project was to understand how admission webhooks, webhook servers, and the Kubernetes admission pipeline actually work, end to end, without relying on higher-level abstractions.

While working on this, one core idea became very clear: Kubernetes is not just a container orchestration platform. It is a distributed API system that stores a desired state and continuously tries to make reality match it. Admission webhooks sit directly on that API boundary. A correct implementation here can fundamentally change how users experience Kubernetes, far beyond just scheduling containers.



## Understanding the admission pipeline
---
Every request to the Kubernetes API server does not immediately become reality. Before an object is persisted into etcd, it flows through a series of admission phases. These phases decide whether a request is allowed, denied, or modified.

An admission webhook is not an async callback or a best-effort hook. It is inline control-plane code. When the API server calls a webhook, it is blocked until a response is returned or a timeout is hit. This means a webhook is effectively part of the API server’s request path.

In this project, the webhook receives real `AdmissionReview` objects from the API server, decodes them, and returns valid `AdmissionResponse` objects. The logs in this repository are not simulated — they come from real admission calls triggered by Kubernetes operations.



## Why a webhook server
---
Writing a webhook server in Go forces you to confront details that are easy to ignore when using policy engines:

* the API server never talks directly to Pods, only Services
* TLS trust is explicit and one-directional
* the API server does not read cluster Secrets
* failure modes can block the entire cluster

By implementing the webhook server directly, it becomes obvious why things like timeouts, idempotency, dry-run handling, and side-effects are taken so seriously in production controllers.

This project intentionally keeps the webhook logic simple so the focus stays on correctness, wiring, and trust — the parts that usually fail silently.



## Challenges encountered
---
### TLS and cert-manager

TLS was the most critical and least forgiving part of the setup. The webhook server certificate must contain correct Service DNS SANs, and the API server will only trust what is embedded in the `caBundle` field of the webhook configuration.

The API server does not “discover” trust from cluster Secrets. If the CA bundle is wrong, the webhook is unreachable even though Pods are healthy. Cert-manager made certificate lifecycle manageable, but understanding where trust actually lives was essential.



### Helm, schemas, and validation

Helm values validation and schema enforcement surfaced mistakes early. Defining a strict `values.schema.json` helped prevent invalid configurations from ever rendering, especially for TLS and webhook parameters.

Using Helm also reinforced how important release-aware naming and namespace scoping are when dealing with control-plane components.



### Kubectl dry-run

Server-side dry runs turned out to be one of the best debugging tools while developing this webhook. They trigger admission without mutating cluster state and made it easy to verify that validation logic, decoding, and responses were correct without risking side effects.

Dry-run awareness also shaped how the webhook should behave: validation is always allowed, but mutation or side effects should never occur during dry-run requests.



### Go struct tags and API types

Working directly with Kubernetes API types highlighted how powerful Go struct tags and JSON annotations are. Correctly matching Kubernetes’ expected wire format is non-negotiable — small mismatches lead to silent failures.

Using the official admission API types made response correctness explicit and eliminated guesswork around serialization.



## Current state
---
At its current stage, this webhook:

* is deployed via Helm and argocd
* serves HTTPS with valid TLS
* is trusted by the Kubernetes API server
* receives and responds to real admission requests
* logs request UIDs for traceability

The webhook currently allows all requests. This is intentional. The goal was to finish with a correct, production-grade admission pipeline, not to rush into policy logic.




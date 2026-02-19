# StreamSpace - Next Steps Guide

## Current Status ✅
- All pods running successfully in `streamspace` namespace
- Web UI accessible via port-forward on `http://localhost:3000`
- API backend running and connected to PostgreSQL
- NATS event bus operational
- Kubernetes controller running

## Issues Identified 🔍

### 1. Template Schema Mismatch (CRITICAL)
**Problem**: Template YAML files use legacy `kasmvnc` field, but the CRD expects `vnc` field.

**Error**: `unknown field "spec.kasmvnc"` when applying templates

**Root Cause**: The codebase is in transition:
- Go types (controller) use modern `VNC` struct ✅
- Template manifests still use legacy `kasmvnc` field ❌
- 36+ template files need migration

### 2. Web UI API Errors
**404 Errors**:
- `/api/v1/catalog/repositories/1/sync`
- `/api/v1/catalog/repositories/2/sync`

**Reason**: Repository sync endpoints not implemented in catalog handler

**500 Errors**:
- `/api/v1/admin/nodes`
- `/api/v1/admin/nodes/stats`

**Reason**: Likely Kubernetes API access issue (RBAC permissions or connection)

### 3. Missing Prometheus
**Error**: `streamspace-prometheus` service not found

**Reason**: Monitoring not enabled in deployment (optional feature)

---

## Solution: Apply Corrected Template

### Step 1: Apply the Corrected Firefox Template

I've created a corrected template file: `firefox-template-corrected.yaml`

This template uses the modern `vnc` field instead of `kasmvnc`.

**Apply it:**
```bash
kubectl apply -f firefox-template-corrected.yaml
```

**Verify it was created:**
```bash
kubectl get templates -n streamspace
```

You should see:
```
NAME              DISPLAYNAME            CATEGORY       IMAGE                                    AGE
firefox-browser   Firefox Web Browser    Web Browsers   lscr.io/linuxserver/firefox:latest      5s
```

### Step 2: Create a Session (Test via kubectl first)

Create a session YAML file: `firefox-session.yaml`

```yaml
apiVersion: stream.space/v1alpha1
kind: Session
metadata:
  name: test-firefox
  namespace: streamspace
spec:
  user: admin
  template: firefox-browser
  state: running
  resources:
    requests:
      memory: "1Gi"
      cpu: "500m"
    limits:
      memory: "2Gi"
      cpu: "1000m"
  persistentHome: false
  idleTimeout: 30m
```

**Apply it:**
```bash
kubectl apply -f firefox-session.yaml
```

**Watch the pod start:**
```bash
kubectl get pods -n streamspace -w
```

Wait for the session pod to reach `Running` state (~30-60 seconds for image pull).

### Step 3: Access the Session

**Find the session pod:**
```bash
kubectl get pods -n streamspace | grep firefox
```

**Port forward to the session:**
```bash
kubectl port-forward -n streamspace <pod-name> 3001:3000
```

**Open in browser:**
```
http://localhost:3001
```

You should see Firefox running in a desktop environment!

---

## Web UI Issues - Investigation Steps

### Check API Logs

**View API logs:**
```bash
kubectl logs -n streamspace -l app.kubernetes.io/name=streamspace-api --tail=100
```

Look for errors related to:
- Kubernetes API connection
- RBAC permissions
- Database queries
- Repository sync

### Check Controller Logs

**View controller logs:**
```bash
kubectl logs -n streamspace -l app.kubernetes.io/name=streamspace-controller --tail=100
```

Look for:
- Template reconciliation errors
- Session creation events
- VNC configuration parsing

### Test API Endpoints Directly

**Test nodes endpoint:**
```bash
curl http://localhost:8080/api/v1/admin/nodes
```

**Test templates endpoint:**
```bash
curl http://localhost:8080/api/v1/templates
```

(You'll need to port-forward the API service first)

---

## Understanding the Template Migration

### Legacy Format (WRONG - Don't use)
```yaml
spec:
  kasmvnc:
    enabled: true
    port: 3000
```

### Modern Format (CORRECT - Use this)
```yaml
spec:
  vnc:
    enabled: true
    port: 3000
    protocol: websocket
    encryption: false
```

### Why the Change?

StreamSpace is migrating from proprietary KasmVNC to 100% open source VNC stack:

**Current (Temporary)**:
- LinuxServer.io containers with KasmVNC
- Port 3000
- WebSocket protocol

**Future (Phase 6)**:
- StreamSpace containers with TigerVNC + noVNC
- Port 5900 (standard RFB)
- Open source stack

The `vnc` field is VNC-agnostic and supports both implementations.

---

## Repository Templates Issue

The Web UI is trying to sync templates from repositories, but:

1. **Repository sync endpoints not implemented** in `api/internal/handlers/catalog.go`
2. **Templates from external repos** use legacy `kasmvnc` field
3. **Need to migrate** 195 templates from your added repository

### Workaround: Apply Templates Manually

For now, apply templates manually using corrected YAML files like the Firefox example.

### Long-term Solution: Implement Repository Sync

The repository sync feature needs:
1. Sync endpoint implementation in catalog handler
2. Template migration script to convert `kasmvnc` → `vnc`
3. Validation before applying templates to cluster

---

## How to Stop StreamSpace

Based on your successful cleanup, here are the verified commands:

### Method 1: Quick Stop (Recommended)

```bash
# Scale down all deployments to 0 replicas
kubectl scale deployment --all --replicas=0 -n streamspace

# Delete all sessions
kubectl delete sessions --all -n streamspace

# Delete the namespace
kubectl delete namespace streamspace
```

### Method 2: Complete Cleanup (Nuclear)

```bash
# Delete namespace (removes everything)
kubectl delete namespace streamspace

# Delete CRDs
kubectl delete crd templates.stream.space
kubectl delete crd sessions.stream.space
kubectl delete crd connections.stream.space
kubectl delete crd templaterepositories.stream.space
```

### Stop Monitoring Stack

If you have Prometheus/Grafana running:

```bash
# Find monitoring resources
kubectl get all --all-namespaces | findstr grafana
kubectl get all --all-namespaces | findstr prometheus

# Delete by label
kubectl delete all -l app.kubernetes.io/name=grafana --all-namespaces
kubectl delete all -l app.kubernetes.io/name=prometheus --all-namespaces
kubectl delete all -l app=kube-prometheus-stack --all-namespaces

# Delete monitoring deployments and statefulsets
kubectl delete deployment --all -n monitoring
kubectl delete statefulset --all -n monitoring
kubectl delete daemonset --all -n monitoring
```

### Stop Port Forwards

**Windows (Git Bash):**
```bash
# Find kubectl processes
tasklist | findstr kubectl

# Kill specific process
taskkill /PID <process-id> /F

# Or kill all kubectl
taskkill /IM kubectl.exe /F
```

### What Stays Running

The `k8s_local-path-provisioner` container is part of Kubernetes infrastructure and is safe to leave running. It uses minimal resources (0.17% CPU) and is needed for persistent volumes.

---

## Next Actions

### Immediate (Test Basic Functionality)
1. ✅ Apply corrected Firefox template
2. ✅ Create a test session via kubectl
3. ✅ Port-forward and access Firefox in browser
4. ✅ Test hibernation (stop session, restart session)

### Short-term (Fix Web UI)
1. Check API logs for specific errors
2. Verify RBAC permissions for nodes API
3. Test API endpoints directly
4. Fix repository sync endpoints (or disable in UI)

### Long-term (Template Migration)
1. Create migration script for all templates
2. Update CRD documentation
3. Migrate database schema (`kasmvnc_*` → `vnc_*` columns)
4. Update external repository templates

---

## Useful Commands

### Check Everything
```bash
# All pods
kubectl get pods -n streamspace

# All templates
kubectl get templates -n streamspace

# All sessions
kubectl get sessions -n streamspace

# Services
kubectl get svc -n streamspace

# Events (troubleshooting)
kubectl get events -n streamspace --sort-by='.lastTimestamp'
```

### Port Forwarding
```bash
# UI
kubectl port-forward -n streamspace svc/streamspace-ui 3000:80

# API
kubectl port-forward -n streamspace svc/streamspace-api 8080:8080

# Session (replace pod-name)
kubectl port-forward -n streamspace <session-pod-name> 3001:3000
```

### Logs
```bash
# API logs
kubectl logs -n streamspace -l app.kubernetes.io/name=streamspace-api -f

# Controller logs
kubectl logs -n streamspace -l app.kubernetes.io/name=streamspace-controller -f

# UI logs
kubectl logs -n streamspace -l app.kubernetes.io/name=streamspace-ui -f
```

### Cleanup
```bash
# Delete a session
kubectl delete session test-firefox -n streamspace

# Delete a template
kubectl delete template firefox-browser -n streamspace

# Delete everything (start over)
kubectl delete namespace streamspace
```

---

## Success Criteria

You'll know everything is working when:

1. ✅ Template applies without errors
2. ✅ Session pod starts and reaches Running state
3. ✅ Port-forward connects successfully
4. ✅ Browser shows Firefox desktop environment
5. ✅ Can interact with Firefox (browse websites)
6. ✅ Session persists after hibernation/wake

---

## Questions to Answer

After testing, we need to understand:

1. **Does the corrected template apply successfully?**
2. **Does the session pod start without errors?**
3. **Can you access Firefox via port-forward?**
4. **What specific errors appear in API logs?**
5. **What specific errors appear in controller logs?**

Share the results and we'll fix any remaining issues!

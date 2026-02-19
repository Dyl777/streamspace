# StreamSpace Design Documentation

**Version:** v2.0-beta
**Last Updated:** 2025-11-26
**Status:** Comprehensive architecture and design documentation for StreamSpace

---

##  Quick Start

### For New Contributors

Start here to understand the system and coding practices:

- **[C4 Architecture Diagrams](architecture/c4-diagrams.md)** - Visual system overview (Context, Container, Component, Code)
- **[Coding Standards](coding-standards.md)** - Go, React/TypeScript, SQL, and Git style guide
- **[Component Library](ux/component-library.md)** - Reusable UI components and patterns

### For Architects & Tech Leads

Understand the key architectural decisions that shape the system:

- **[ADR Log](architecture/adr-log.md)** - All architecture decision records
- **[ADR-004: Multi-Tenancy](architecture/adr-004-multi-tenancy-org-scoping.md)** -  **CRITICAL** - Org-scoped RBAC (Issues #211, #212)
- **[ADR-005: WebSocket Dispatch](architecture/adr-005-websocket-command-dispatch.md)** - Command dispatch architecture
- **[ADR-006: Database Source of Truth](architecture/adr-006-database-source-of-truth.md)** - Database-first design pattern
- **[ADR-007: Agent Outbound WebSocket](architecture/adr-007-agent-outbound-websocket.md)** - Firewall-friendly agent connections
- **[ADR-008: VNC Proxy](architecture/adr-008-vnc-proxy-control-plane.md)** - Centralized VNC access control
- **[ADR-009: Helm Deployment](architecture/adr-009-helm-deployment-no-operator.md)** - Deployment strategy (no K8s Operator)

### For Product Managers

Understand feature lifecycle and acceptance criteria:

- **[Product Lifecycle](product/product-lifecycle.md)** - API versioning, feature maturity, deprecation policies
- **[Acceptance Criteria Guide](acceptance-criteria-guide.md)** - Feature definition with Given-When-Then format
- **[Information Architecture](ux/information-architecture.md)** - UI navigation and page hierarchy

### For SREs & Operations

Production deployment, scaling, and operational procedures:

- **[Load Balancing & Scaling](operations/load-balancing-and-scaling.md)** - Production operations guide (1,000+ sessions)
- **[Industry Compliance](compliance/industry-compliance.md)** - SOC 2, HIPAA, FedRAMP readiness
- **[Vendor Assessment](vendor-assessment.md)** - Third-party risk evaluation template

### For Security Engineers

Security architecture, compliance, and risk management:

- **[ADR-004: Multi-Tenancy](architecture/adr-004-multi-tenancy-org-scoping.md)** - Org isolation and security boundaries
- **[ADR-001: VNC Token Auth](architecture/adr-001-vnc-token-auth.md)** - VNC authentication mechanism
- **[Industry Compliance](compliance/industry-compliance.md)** - Compliance controls mapping (SOC 2, HIPAA)
- **[Vendor Assessment](vendor-assessment.md)** - Security assessment checklist

### For QA & Test Engineers

Testing standards and acceptance criteria:

- **[Acceptance Criteria Guide](acceptance-criteria-guide.md)** - Feature testing with scenarios
- **[Coding Standards](coding-standards.md)** - Testing conventions and coverage requirements

---

## � Directory Structure

```
docs/design/
├── README.md                           # This file - documentation index
│
├── architecture/                       # Architecture Decision Records (ADRs)
│   ├── adr-log.md                     # Index of all ADRs
│   ├── adr-template.md                # Template for new ADRs
│   ├── adr-001-vnc-token-auth.md      # VNC authentication
│   ├── adr-002-cache-layer.md         # Redis caching strategy
│   ├── adr-003-agent-heartbeat-contract.md  # Agent health protocol
│   ├── adr-004-multi-tenancy-org-scoping.md # CRITICAL: Multi-tenancy security
│   ├── adr-005-websocket-command-dispatch.md # WebSocket vs NATS
│   ├── adr-006-database-source-of-truth.md   # Database-first architecture
│   ├── adr-007-agent-outbound-websocket.md   # Agent connection pattern
│   ├── adr-008-vnc-proxy-control-plane.md    # VNC proxy architecture
│   ├── adr-009-helm-deployment-no-operator.md # Deployment strategy
│   └── c4-diagrams.md                 # System architecture visualizations
│
├── ux/                                # User Experience & UI design
│   ├── information-architecture.md    # Site map, navigation, URL structure
│   └── component-library.md           # Reusable UI components
│
├── operations/                        # Production operations
│   └── load-balancing-and-scaling.md  # Scaling guide, capacity planning
│
├── compliance/                        # Regulatory compliance
│   └── industry-compliance.md         # SOC 2, HIPAA, FedRAMP
│
├── product/                           # Product management
│   └── product-lifecycle.md           # Feature maturity, API versioning
│
├── acceptance-criteria-guide.md       # Feature definition standards
├── coding-standards.md                # Go, React/TS, SQL, Git conventions
├── retrospective-template.md          # Sprint retrospective format
└── vendor-assessment.md               # Third-party risk evaluation
```

---

##  ADR Quick Reference

Architecture Decision Records (ADRs) document significant architectural choices:

| ADR | Status | Priority | Description |
|-----|--------|----------|-------------|
| [ADR-001](architecture/adr-001-vnc-token-auth.md) |  Accepted | High | VNC token authentication mechanism |
| [ADR-002](architecture/adr-002-cache-layer.md) |  Accepted | Medium | Redis cache layer for session metadata |
| [ADR-003](architecture/adr-003-agent-heartbeat-contract.md) |  In Progress | High | Agent heartbeat & health check protocol |
| [ADR-004](architecture/adr-004-multi-tenancy-org-scoping.md) |  Accepted |  **CRITICAL** | Multi-tenancy via org-scoped RBAC |
| [ADR-005](architecture/adr-005-websocket-command-dispatch.md) |  Accepted | High | WebSocket command dispatch (vs NATS) |
| [ADR-006](architecture/adr-006-database-source-of-truth.md) |  Accepted | High | Database as source of truth |
| [ADR-007](architecture/adr-007-agent-outbound-websocket.md) |  Accepted | High | Agent outbound WebSocket connections |
| [ADR-008](architecture/adr-008-vnc-proxy-control-plane.md) |  Accepted | High | VNC proxy via Control Plane |
| [ADR-009](architecture/adr-009-helm-deployment-no-operator.md) |  Accepted | Medium | Helm chart deployment (no Operator) |

**Legend:**
-  **Accepted** - Decision implemented and in production
-  **In Progress** - Decision made, implementation underway
-  **Proposed** - Under review, not yet implemented
-  **CRITICAL** - P0 priority, security or system-critical

---

## � Document Types

### Architecture Decision Records (ADRs)

**Purpose:** Document significant architectural decisions with context, alternatives, and consequences.

**Format:** Structured markdown with status, date, context, decision, alternatives, consequences.

**Location:** `architecture/adr-*.md`

**Process:**
1. Copy `architecture/adr-template.md`
2. Fill in context, decision, alternatives, consequences
3. Submit PR for review
4. Merge when accepted

### Design Documents

**Purpose:** Comprehensive design specifications for features, systems, or processes.

**Format:** Free-form markdown with clear structure.

**Location:** Various directories (ux, operations, compliance, product)

**Examples:**
- C4 Architecture Diagrams (visual system overview)
- Load Balancing & Scaling (operational guide)
- Industry Compliance (regulatory mapping)

### Standards & Guidelines

**Purpose:** Project-wide conventions and best practices.

**Format:** Reference documentation with examples.

**Examples:**
- Coding Standards (Go, React/TypeScript, SQL, Git)
- Acceptance Criteria Guide (feature definition)
- Retrospective Template (team process)

---

## � External Resources

### Full Design & Governance Documentation

**Private Repository:** `streamspace-dev/streamspace-design-governance`

Contains comprehensive design documentation including:
- Stakeholder requirements
- System design specifications
- UX mockups and wireframes
- Delivery plans and timelines
- Risk and governance documentation
- Security and compliance deep dives

**Access:** Internal team only (contains sensitive planning and vendor assessments)

### Public Documentation

**User-Facing Documentation:** See `/docs/` in main repository
- [ARCHITECTURE.md](../ARCHITECTURE.md) - High-level system overview
- [DEPLOYMENT.md](../../DEPLOYMENT.md) - Installation and deployment guide
- [FEATURES.md](../../FEATURES.md) - Feature status and roadmap
- [TROUBLESHOOTING.md](../TROUBLESHOOTING.md) - Common issues and solutions

---

##  Contributing to Documentation

### When to Create an ADR

Create an ADR when making decisions that:
- Affect multiple components or teams
- Have significant consequences (performance, security, cost)
- Involve trade-offs between alternatives
- Need to be explained to future contributors

**Examples:**
- Choosing a database (PostgreSQL vs MySQL)
- Authentication mechanism (JWT vs session cookies)
- Deployment model (Operator vs Helm chart)

**Not ADR-worthy:**
- Library choice for minor feature (just use best practice)
- Code refactoring (use PR description)
- Bug fixes (use commit message)

### How to Update Existing Documentation

1. **Read the document** - Understand current state
2. **Make changes** - Update content, add sections, fix errors
3. **Update metadata** - Change "Last Updated" date
4. **Submit PR** - Include rationale for changes
5. **Tag reviewers** - Assign relevant stakeholders

### Documentation Review Process

**Design Docs & ADRs:** Reviewed in PRs (1 approval required)

**Reviewers:**
- Architects: All ADRs, architecture changes
- Product: Product lifecycle, acceptance criteria
- SRE: Operations, scaling, compliance
- Security: ADRs with security impact

---

##  Documentation Quality Standards

### Good Documentation Is:

- **Accurate** - Reflects current state of system
- **Complete** - Covers all necessary details
- **Concise** - No unnecessary information
- **Well-structured** - Clear headings, logical flow
- **Up-to-date** - Last Updated date within 6 months
- **Discoverable** - Linked from index, easy to find

### Documentation Checklist

- [ ] Clear title and purpose
- [ ] Metadata (version, date, status, owner)
- [ ] Table of contents (for docs >500 lines)
- [ ] Code examples (where applicable)
- [ ] Diagrams (architecture, flows, sequences)
- [ ] References to related docs
- [ ] Last Updated date

---

##  Finding Documentation

### By Role

Use the [Quick Start](#-quick-start) section above - organized by role (Developer, Architect, PM, SRE, Security, QA).

### By Topic

| Topic | Documents |
|-------|-----------|
| **Architecture** | ADR-001 to ADR-009, C4 Diagrams |
| **Multi-Tenancy** | ADR-004 |
| **Authentication** | ADR-001 (VNC tokens), ADR-004 (org RBAC) |
| **Caching** | ADR-002 |
| **Agents** | ADR-003 (heartbeat), ADR-007 (WebSocket), ADR-009 (deployment) |
| **VNC** | ADR-001 (auth), ADR-008 (proxy) |
| **Scaling** | Load Balancing & Scaling |
| **Compliance** | Industry Compliance Matrix |
| **UI/UX** | Information Architecture, Component Library |
| **Testing** | Acceptance Criteria Guide |
| **Operations** | Load Balancing & Scaling, Product Lifecycle |

### By GitHub Issue

ADRs are linked to relevant GitHub issues:
- Issue #211 → ADR-004 (WebSocket org scoping)
- Issue #212 → ADR-004 (Org context & RBAC)
- Issue #214 → ADR-002 (Cache layer)
- Issue #215 → ADR-003 (Agent heartbeat)

---

## � Documentation Maintenance

### Review Schedule

- **ADRs:** Review on implementation or annually
- **Design Docs:** Review quarterly or on major version
- **Standards:** Review semi-annually

### Deprecation Process

When architectural decisions change:
1. Update ADR status to "Superseded"
2. Add "Superseded By" section linking to new ADR
3. Keep original ADR for historical context
4. Do NOT delete superseded ADRs

### Feedback

**Questions or issues with documentation?**
- Open a GitHub issue with label `documentation`
- Tag with relevant area (architecture, ux, operations)
- Assign to documentation owner if known

---

## � Documentation Stats

**Current Status (v2.0-beta):**
- **Total ADRs:** 9 (9 Accepted, 0 Proposed)
- **Design Docs:** 10 (Phase 1 + Phase 2 complete)
- **Total Lines:** ~7,600 lines
- **Last Major Update:** 2025-11-26 (Documentation Sprint)

**Coverage:**
-  Architecture: Comprehensive (9 ADRs)
-  Operations: Complete (scaling, compliance)
-  Development: Complete (coding standards, components)
-  Product: Complete (lifecycle, acceptance criteria)
- ⏳ UX: Good (IA, components) - wireframes in private repo

---

## � Contact & Support

**Documentation Questions:**
- GitHub Issues: Tag with `documentation` label
- Team Channel: #documentation (Slack/Discord)
- Email: architecture@streamspace.dev

**Maintainers:**
- Architecture: Agent 1 (Architect)
- Operations: SRE Team
- Product: Product Management
- UX: Design Team

**Next Documentation Review:** Q1 2026 (post v2.0 GA)

---

**Last Updated:** 2025-11-26
**Version:** 1.0 (v2.0-beta documentation sprint)
**Changelog:**
- 2025-11-26: Initial comprehensive documentation index created
- 2025-11-26: Added 9 ADRs and 10 design documents

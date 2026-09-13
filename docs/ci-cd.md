# CrashBench CI/CD Integration Guide

CrashBench provides enterprise-grade automated pipeline gating, vulnerability scanning, and test reporting.

---

## 1. Automated Pipeline Gating (`crashbench check`)

Run `crashbench check` in your CI/CD workflows to prevent fragile or insecure agents from merging into production:

```bash
crashbench check \
  --target "npx claude-code" \
  --name "Claude Code (CI)" \
  --fail-under 75 \
  --sarif results.sarif \
  --junit junit.xml \
  --summary-md pr-comment.md
```

### Exit Codes
- `0`: The agent's CrashBench Resilience Index (CRI) meets or exceeds `--fail-under`.
- `1`: The agent scored below `--fail-under` (pipeline blocked).

---

## 2. GitHub Code Scanning Alerts (SARIF 2.1.0)

CrashBench exports standard **SARIF 2.1.0** reports, uploading detected terminal vulnerabilities directly to the **GitHub Security > Code scanning alerts** dashboard:

```yaml
- name: Upload SARIF Security Alerts
  if: always()
  uses: github/codeql-action/upload-sarif@v3
  with:
    sarif_file: results.sarif
```

Each failed scenario generates an alert tagged with its corresponding **CWE** (e.g. `CWE-200`, `CWE-835`) and **OWASP GenAI** category.

---

## 3. JUnit XML Test Dashboards

Integrate CrashBench results into standard test report visualizers across Jenkins, CircleCI, and GitLab CI:

```yaml
- name: Publish Test Results
  if: always()
  uses: EnricoMi/publish-unit-test-result-action@v2
  with:
    files: junit.xml
```

---

## 4. One-Command CI Bootstrapping

Generate a complete GitHub Actions workflow file in seconds:

```bash
crashbench init-ci
```

This generates `.github/workflows/crashbench.yml`.

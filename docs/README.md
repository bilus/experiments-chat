# Experiments Chat Docs

This directory captures ACP spike findings that directly affect the planned chat requirements.

The goal is not to preserve the full spike history. The goal is to keep the product-relevant decisions close to the prototype:

- how each target provider is launched
- which ACP capabilities looked portable
- where provider-specific adaptation is still required
- how the prototype should structure its internal provider contract
- what logging should capture once real ACP calls replace the static provider

## Files

- [agent-requirements.md](agent-requirements.md) documents EARS requirements, GitHub issue shape, and project tracking.
- [agent-testing.md](agent-testing.md) documents TDD, Gherkin, Godog, and Go unit test expectations.
- [agent-workflow.md](agent-workflow.md) documents the `development` branch, `gps`, and patch-stack flow.
- [agent-environment.md](agent-environment.md) documents devbox, direnv, and command conventions.
- [provider-launch.md](provider-launch.md) records how Claude, Copilot, and Codex were started during the spike.
- [compatibility-findings.md](compatibility-findings.md) summarizes the matrix results that matter for requirements.
- [agent-provider-contract.md](agent-provider-contract.md) explains the prototype boundary between Agent and Provider.
- [logging-and-observability.md](logging-and-observability.md) captures the logging requirements inherited from the spike.
- [planning-guardrails.md](planning-guardrails.md) preserves the stage boundaries that should shape future issues.

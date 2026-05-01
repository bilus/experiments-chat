# Planning Guardrails

The ACP spike established stage boundaries that should continue to shape this repo.

## Stage 1

Stage 1 is compatibility and architecture de-risking.

Allowed work:

- minimal terminal chat
- ACP provider adapters
- compatibility tests
- structured logs
- small harnesses that prove provider behavior

Avoid expanding Stage 1 into the full wiki workflow.

## Stage 2

Stage 2 is the ACP-native Wikigen successor.

It should support:

- provider-agnostic execution
- customizable prompts
- local markdown output
- regular single-project repos
- monorepos with explicit project selection
- durable logs and run artifacts

Stage 2 must be thoroughly discussed and designed before substantial implementation starts.

## Stage 3

Stage 3 is only a sketch.

The current idea is an LLM wiki MCP layer for GitHub Wiki, inspired by the existing OtterWiki MCP direction and Karpathy's LLM-maintained wiki idea.

Do not let Stage 3 drive near-term architecture more strongly than Stage 1 and the eventually approved Stage 2 design.


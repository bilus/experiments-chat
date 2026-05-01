# ACP Provider Launch Findings

The prototype should treat ACP as the integration layer, but it must not assume every target exposes ACP the same way.

## Provider Startup

| Provider | Spike startup path | Finding for requirements |
| --- | --- | --- |
| Claude | `claude-code-acp` or `npx -y @zed-industries/claude-code-acp@latest` | Claude is reached through an ACP wrapper around the local Claude CLI. It must not use the direct Anthropic API path. |
| Copilot | `copilot --acp --stdio` | Copilot exposes ACP through the installed Copilot CLI over stdio. |
| Codex | `codex-acp` | Codex is reached through an ACP wrapper around the local Codex CLI. |

## Go Client

The spike used `github.com/coder/acp-go-sdk` as the Go ACP client library.

That makes it the first candidate for the real provider adapter, but the app should keep this dependency behind a small internal adapter so later ACP library changes do not leak into chat and workflow code.

## Requirement Impact

The provider configuration should eventually include:

- provider name
- command and arguments
- working directory
- environment overrides
- provider version when available
- whether the provider is native ACP or an ACP wrapper

The first terminal chat iteration can keep this static, but the boundary should be ready for real provider commands.


# Section Template

Use this template when adding or refactoring sections in `STYLE.MD`.

## Template

```md
### Section name

> [!IMPORTANT]
> Optional short callout when needed.

**📌 Rule**:

State the rule directly and concretely.

**🤔 Why**:

Explain why the rule exists and what kind of errors, complexity, or ambiguity it prevents.

**❌ Bad**:

Explain what is wrong with the anti-pattern.

```go
// bad example
```

**✅ Good**:

Explain why the preferred version is better.

```go
// good example
```

**📝 Exceptions / Notes**:

Add this block only when the rule needs caveats, scope limits, or explicit non-goals.

**🔗 Related rules**:

- [Some other rule](#some-other-rule)
- [Another rule](#another-rule)
```

## Recommendations

**Use the smallest example that proves the point**:

- Prefer one-file examples by default.
- Use multi-file examples only when package boundaries, exported APIs, or cross-type composition are part of the rule.
- If a large example is kept on purpose, make sure the extra size teaches something that a smaller example cannot show.

**Mark intentionally broken examples explicitly**:

- If a bad example is not supposed to compile, say so near the example.
- Use that sparingly and only when the brokenness itself is part of the lesson.

**Do not force every block into every section**:

- `Rule` and `Why` are the default minimum.
- Add `Bad` and `Good` when contrast helps.
- Add `Exceptions / Notes` only when there is a real caveat.
- Add `Related rules` when the reader would benefit from jumping to nearby concepts.

**Prefer doc comments over inline comments inside code**:

- If the explanation belongs to the symbol, put it in the doc comment.
- Avoid comments inside function bodies in preferred examples.

**Keep the tone directive, but readable**:

- Write rules as project policy, not as brainstorming.
- Keep explanations short unless extra detail materially improves onboarding.

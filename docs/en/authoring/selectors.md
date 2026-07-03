# Selectors

Scenaria uses **Playwright** locators. In steps, the selector is a quoted string: `нажимаю "#login"`, `вижу "text=Save"`.

## Recording priority

During live recording (`internal/selector/strategies.go`):

| Strategy | Example | Prefer when |
|----------|---------|-------------|
| **testid** | `[data-testid="submit"]` | You control `data-testid` in the app |
| **id** | `#login-form` | Unique stable `id` |
| **aria / role** | `role=button[name="Sign in"]` | Accessible buttons and links |
| **label** | `label=Email` | Inputs with `<label>` |
| **placeholder** | `[placeholder="Search…"]` | Fields without label |
| **text** | `button:has-text("Catalog")` | Visible text with tag hint |

Clicks: short text selector first; use `>>` chains when multiple matches exist.  
Fields: label → placeholder → aria → name → testid → id.

## Chains (`>>`)

Nested elements use Playwright chaining:

```gherkin
Когда нажимаю "nav >> text=Products"
Тогда вижу "main >> h1"
```

Scenaria resolves via `ResolveChainedLocator`: container first, then target.

## Best practices

1. Add **`data-testid`** on buttons, forms, and key blocks.
2. Prefer **`role` + `name`** over long CSS paths.
3. Shorten auto-recorded selectors manually — avoid bare `div`/`img`.
4. Scenaria UI is marked `data-scenaria-ui` and filtered during recording.
5. Hover menus: add `навожу` before click or chain to menu item.

## Variables `{{name}}`

```gherkin
Когда нажимаю "{{button_id}}"
И заполняю "{{email}}" в "input[type=email]"
```

- Nesting up to 16 passes; cycles error out.
- Use `запоминаю … как "имя"` for complex values.

## Validation

| Method | How |
|--------|-----|
| IDE | **On-page selectors** or **Validate in browser** |
| CLI | `scenaria validate ./features --browser` |

URL assertion steps wait up to ~5s for SPA navigation.

## Related

- [Recording](../gui/recording.md)
- [Gherkin](gherkin.md)

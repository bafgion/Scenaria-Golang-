# Bundled examples

The `examples/` folder ships with Scenaria and uses [example.com](https://example.com) — no login required.

## Open in IDE

1. **Start** → **Open example scenarios**
2. Select a file in the catalog
3. **Run test** (Ctrl+Enter) or **Dry-run** first

Or **Project → Open project…** → select the `examples` folder.

## Files

| File | Demonstrates |
|------|--------------|
| `01-pervaya-proverka.feature` | Open page, assert heading, close browser (`@smoke`) |
| `02-perehod-po-ssylke.feature` | Click link, pause, back, assert |
| `03-uslovnyy-shag.feature` | `Если вижу …` conditional block |
| `04-tablica-primerov.feature` | Scenario outline + `Примеры:` table (2 runs) |
| `05-testclient-kontekst.feature` | `Контекст:` + TestClient `DemoUser` |

## TestClient demo

`05-testclient-kontekst.feature` loads `examples/.scenaria/test_clients/DemoUser.json` (empty demo session). For real projects:

1. **Record → TestClient…** — log in, save profile
2. Add `Контекст:` with `Дано я подключаю TestClient "name"`

## Outline table

`04-tablica-primerov.feature` runs the template once per table row; catalog shows badge **2 examples**. Placeholders like `<url>` map to columns.

## Tips

- Replace URLs with your app when ready
- Full step list: **Ctrl+Space** or **F1**
- Record from scratch: **Ctrl+B** → **Start recording**

## Related

- [First project](../getting-started/first-project.md)
- [Gherkin](gherkin.md)

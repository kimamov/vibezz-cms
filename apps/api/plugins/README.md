# Plugins

Plugins add custom content types with their own Go logic (models, HTTP handlers).

## Loading plugins

Plugins are loaded at startup. The list is in [registry.go](registry.go). To enable a plugin, add it to the `All()` return slice.

## Adding a custom plugin

1. **Create a subpackage** under `plugins/`, e.g. `plugins/myplugin/`.

2. **Implement the `plugin.Plugin` interface** (see [internal/plugin/plugin.go](../internal/plugin/plugin.go)):
   - `Name() string` – display name
   - `Slug() string` – unique slug (content type slug and URL prefix)
   - `ContentType() plugin.ContentTypeDefinition` – name and fields for the content type
   - `Register(deps *plugin.Deps)` – mount admin and public routes, run init logic

3. **Use `plugin.Deps`** in `Register`:
   - `ContentTypeID` – your plugin’s content type UUID (resolved by the loader)
   - `EntryService`, `ContentTypeService`, `MediaService` – core services
   - `Admin` – authenticated admin router group (`/api/admin`)
   - `Public` – public router group (`/api/public`)

4. **Register in the registry**: in [registry.go](registry.go), import your package and append `myplugin.New()` to the slice.

## Example: News plugin

The [news](news/) plugin adds a “News” content type with:

- **Model** ([model.go](news/model.go)) – `NewsItem` struct and `FromEntry()` for typed access to fields
- **Content type** – headline, excerpt, body, featured (via [plugin.go](news/plugin.go) `ContentType()`)
- **Admin API** – `GET/POST /api/admin/news`, `GET/PATCH/DELETE /api/admin/news/:id`, publish/unpublish
- **Public API** – `GET /api/public/news`, `GET /api/public/news/:slug`

Entries use a path prefix (`/news/...`) so they don’t collide with pages. For path-based plugins you can set `PathPrefix` when creating entries via `content.CreateEntryInput.PathPrefix`.

---
name: fyne-ui
description: Use this skill when building, improving, refactoring, or integrating Go/Fyne desktop UI screens, dashboards, forms, reports, themes, widgets, navigation, CRUD pages, and polished UI/UX flows.
---

# Fyne UI Skill

You are working on a Go desktop application using Fyne.

## Main goal

Build beautiful, functional, maintainable, and integrated Fyne interfaces. Do not generate ugly tutorial-style UIs. Prefer production-style administrative interfaces with clear layout, spacing, visual hierarchy, reusable components, and real integration with the application's services.

## Before coding

1. Inspect the existing project structure.
2. Identify where UI, services, domain models, repositories, and reports are located.
3. Do not put everything in `main.go`.
4. Preserve the existing architecture.
5. If there is no structure, propose and create a simple structure:

```txt
cmd/desktop/main.go
internal/ui/window.go
internal/ui/theme.go
internal/ui/pages
internal/ui/components
internal/domain
internal/service
internal/repository
internal/report
```


## Dependency injection and main composition root

This project uses a composition-root pattern.

The `main.go` file is responsible for initializing all concrete dependencies used by the application and passing them into the app through constructors.

The `main.go` may instantiate:

- Config.
- Logger.
- Database connection.
- Repositories.
- Services.
- Report generators.
- UI dependencies.
- The Fyne app instance.
- The main application/window.

The `main.go` must wire dependencies, but it must not contain business rules, page layout details, report generation logic, database queries, or UI component implementation.

Correct responsibility for `main.go`:

```go
func main() {
    cfg := config.Load()

    db := database.MustOpen(cfg.DatabaseURL)
    defer db.Close()

    employeeRepo := repository.NewEmployeeRepository(db)
    employeeService := service.NewEmployeeService(employeeRepo)

    reportGenerator := report.NewGenerator()

    fyneApp := app.NewWithID("br.com.example.app")
    fyneApp.Settings().SetTheme(ui.NewAppTheme())

    application := ui.NewApplication(ui.ApplicationDeps{
        FyneApp:          fyneApp,
        EmployeeService:  employeeService,
        ReportGenerator:  reportGenerator,
    })

    application.Run()
}
```

Incorrect responsibility for `main.go`:

```go
func main() {
    // Do not build all screens directly here.
    // Do not place all widget creation here.
    // Do not query the database directly from buttons.
    // Do not implement report generation here.
    // Do not put business rules here.
}
```

Use dependency injection through constructors.

Example:

```go
type ApplicationDeps struct {
    FyneApp         fyne.App
    EmployeeService service.EmployeeService
    ReportGenerator report.Generator
}

func NewApplication(deps ApplicationDeps) *Application {
    return &Application{
        app:             deps.FyneApp,
        employeeService: deps.EmployeeService,
        reportGenerator: deps.ReportGenerator,
    }
}
```

Pages should also receive dependencies explicitly:

```go
type EmployeePageDeps struct {
    EmployeeService service.EmployeeService
}

func NewEmployeePage(deps EmployeePageDeps) fyne.CanvasObject {
    // Build the employee page using deps.EmployeeService.
}
```

Rules:

1. `main.go` initializes concrete dependencies.
2. `main.go` passes dependencies to the application.
3. The application passes dependencies to pages/components when needed.
4. UI pages depend on service interfaces, not concrete repositories.
5. Repositories depend on database connections or clients.
6. Services contain business rules.
7. Pages contain presentation logic.
8. Components contain reusable visual structure.
9. Report generators live outside UI code.
10. No hidden global dependency initialization inside pages.

If the existing project already has a dependency injection pattern, preserve it and extend it instead of replacing it.

## UI principles

Every screen must have:

1. Clear title.
2. Short description.
3. Primary action.
4. Consistent spacing.
5. Loading state when data is fetched.
6. Empty state when there is no data.
7. Error state when something fails.
8. Success feedback after saving/exporting.
9. Separation between layout, component, and business logic.

## Layout rules

Prefer:

- `container.NewBorder` for app shell.
- Sidebar on the left.
- Header/topbar at the top.
- Main content in the center.
- `container.NewPadded` for breathing room.
- `container.NewScroll` for long content.
- `container.NewGridWithColumns` or `container.NewGridWrap` for cards.

Avoid:

- Giant vertical stacks without scroll.
- Mixing business logic directly inside widgets.
- Hardcoded repeated UI blocks.
- Unlabeled inputs.
- Tiny buttons with unclear purpose.
- Creating one new window for each screen without a strong reason.

## Visual style

Use a custom theme when the app needs a branded look.

Create reusable components for:

- Stat cards.
- Navigation items.
- Section headers.
- Status badges.
- Empty states.
- Error panels.
- Loading panels.
- Table toolbar.
- Report toolbar.
- Form sections.
- Confirmation dialogs.

Prefer Fyne theme colors when possible. Use `canvas` only when default widgets are not enough.

The app should feel like a serious internal system, not a tutorial demo.

## App shell

For administrative apps, prefer this structure:

```txt
App window
├── Header/topbar
├── Sidebar navigation
└── Main content area
```

The shell should be reusable and should not know the details of each page.

The shell may receive:

- App title.
- Navigation items.
- Router/content container.
- Current user information.
- Theme settings.
- Logout callback.

## Navigation

Use a small router abstraction that swaps the central content container.

Expected behavior:

1. Sidebar item is clicked.
2. Router replaces central page.
3. Selected sidebar item updates visual state.
4. Page receives dependencies through constructor.

Example constructor style:

```go
type Router struct {
    content *fyne.Container
}

func NewRouter() *Router {
    return &Router{
        content: container.NewStack(),
    }
}

func (r *Router) Show(page fyne.CanvasObject) {
    r.content.Objects = []fyne.CanvasObject{page}
    r.content.Refresh()
}

func (r *Router) View() fyne.CanvasObject {
    return r.content
}
```

Do not use global mutable navigation state unless the existing project already follows that pattern.

## Forms

Forms must be grouped by meaning.

For example:

- Personal data.
- Contract data.
- Address data.
- System access.
- Actions.

Rules:

1. Validate fields before saving.
2. Show validation errors clearly.
3. Disable save button while saving.
4. Show success feedback after save.
5. Keep cancel/back behavior obvious.
6. Do not make huge flat forms with dozens of ungrouped fields.
7. Use placeholders only as examples, not as replacements for labels.

## Tables

For lists and reports, prefer `widget.Table`.

Tables should support:

- Column widths.
- Header row.
- Filtering when requested.
- Row action buttons when needed.
- Refresh after data changes.
- Empty state if there is no data.
- Error state if loading fails.

Table screens should usually have:

```txt
Title
Description
Toolbar with search/filter/actions
Summary cards when useful
Table
Pagination or footer when useful
```

Avoid raw, cramped tables with no surrounding context.

## Dashboards

A dashboard should summarize work and guide action.

Prefer:

- Stat cards.
- Pending items.
- Alerts.
- Recent activity.
- Quick actions.
- Small charts only when they help understanding.

Avoid:

- Random decorative charts.
- Too many numbers without explanation.
- Filling the screen with widgets just to look busy.

## Reports

Report screens should include:

- Title.
- Period/filter area.
- Summary cards.
- Result table.
- Export PDF action.
- Export XLSX action when applicable.
- Generated-at metadata when exporting.
- Clear empty state.

Do not make reports only as raw tables. Always include summary and filters when useful.

PDF and XLSX generation should live outside the UI layer, preferably in `internal/report`.

## Integration

UI must call service interfaces, not repositories directly.

All concrete dependencies must be initialized in `main.go` or in a composition function called by `main.go`.

Prefer dependency injection:

```go
type EmployeePageDeps struct {
    EmployeeService service.EmployeeService
}

func NewEmployeePage(deps EmployeePageDeps) fyne.CanvasObject {
    // Build screen here.
}
```

If the real service does not exist yet, create a small interface and a mock implementation only when needed.

Do not couple Fyne widgets directly to database code.

## Concurrency

Do not block the UI thread with long operations.

For loading data:

1. Show loading state.
2. Run the service call in a goroutine.
3. Update the UI after completion.
4. Show error or data state.

For long exports:

1. Show progress/loading feedback.
2. Disable the export button while running.
3. Re-enable when finished.
4. Show success or error dialog.

Avoid freezing the window during I/O, database operations, HTTP requests, or report generation.

## Error handling

Every important user action must handle errors.

Examples:

- Failed to load data.
- Failed to save form.
- Failed to export report.
- Invalid required field.
- Service unavailable.
- Empty result after filtering.

Error messages should be understandable to a normal system user.

Bad:

```txt
failed with code -1
```

Good:

```txt
Não foi possível carregar os funcionários. Verifique a conexão ou tente novamente.
```

## Empty states

Every list/report page should have a good empty state.

Examples:

```txt
Nenhum funcionário encontrado.
Cadastre um novo funcionário ou ajuste os filtros.
```

```txt
Nenhuma guia encontrada para o período.
Altere o filtro de competência ou gere uma nova guia.
```

Empty states should usually include a next action.

## Code organization

Generated code must:

- Compile.
- Use clear names.
- Avoid global mutable state.
- Keep components small.
- Keep UI files organized.
- Use Go idioms.
- Avoid magic numbers when a theme/component can solve it.
- Include only useful comments.
- Avoid copy/paste UI blocks.
- Prefer constructors for pages and components.

Recommended structure:

```txt
internal/ui/
├── window.go
├── theme.go
├── router.go
├── components/
│   ├── stat_card.go
│   ├── sidebar.go
│   ├── section.go
│   ├── badge.go
│   ├── empty_state.go
│   ├── loading.go
│   └── toolbar.go
└── pages/
    ├── dashboard.go
    ├── employees.go
    ├── reports.go
    └── settings.go
```

## Custom components

Create custom reusable components when default widgets are not enough.

Useful components:

- `StatCard`
- `Sidebar`
- `NavItem`
- `StatusBadge`
- `SectionHeader`
- `EmptyState`
- `ErrorState`
- `LoadingState`
- `ReportToolbar`
- `TableToolbar`
- `FormSection`

For visual customization, use:

- `canvas.Rectangle`
- `canvas.Text`
- `canvas.Line`
- `container.NewStack`
- `container.NewPadded`
- `theme` colors and sizes

For deeper widgets:

- Embed `widget.BaseWidget`.
- Call `ExtendBaseWidget`.
- Implement `CreateRenderer`.
- Keep renderer code isolated.

## Theming

Use a custom theme when asked for a branded or polished app.

The theme may define:

- Primary color.
- Background color.
- Foreground color.
- Button color.
- Padding size.
- Text size.
- Icon behavior.
- Font fallback.

Do not hardcode colors everywhere if the theme can handle them.

## Fyne-specific preferences

Use:

- `fyne.io/fyne/v2/app`
- `fyne.io/fyne/v2/container`
- `fyne.io/fyne/v2/widget`
- `fyne.io/fyne/v2/canvas`
- `fyne.io/fyne/v2/theme`
- `fyne.io/fyne/v2/dialog`
- `fyne.io/fyne/v2/data/binding` when useful

Prefer:

- `container.NewBorder` for app shell.
- `container.NewGridWithColumns` for dashboard cards.
- `container.NewGridWrap` for responsive cards.
- `container.NewStack` for custom visual overlays.
- `container.NewScroll` for long pages.
- `widget.Table` for large data lists.
- `widget.Form` for simple forms, or custom layouts for complex forms.

## When asked to create a new screen

Deliver:

1. New page file.
2. Reusable components if needed.
3. Wiring in navigation/sidebar.
4. Service interface if integration is needed.
5. Mock service only if real service does not exist.
6. Clear TODOs only where integration needs missing user information.

The screen must have:

- Title.
- Description.
- Main action.
- Loading state.
- Empty state.
- Error state.
- Main content state.

## When asked to improve an existing screen

Do not rewrite the whole app blindly.

First:

1. Inspect current screen.
2. Identify layout, UX, and integration problems.
3. Refactor incrementally.
4. Preserve behavior unless asked to change it.
5. Improve visual hierarchy, spacing, and components.
6. Keep the code compiling.

## When asked for reports

Create or update:

- Report page.
- Filters.
- Summary cards.
- Result table.
- Export action.
- Report service or report generator.

PDF/XLSX export must not be implemented directly inside button callbacks except for very small glue code.

## When asked for integration

Connect UI to existing services.

Rules:

1. Search existing service interfaces.
2. Reuse existing domain models.
3. Do not duplicate models unnecessarily.
4. Do not call repositories directly from UI.
5. Keep UI independent from storage details.
6. Use dependency injection through page constructors.

## Quality bar

Before finishing, check:

- Does the code compile?
- Is the screen readable?
- Is the primary action clear?
- Is spacing consistent?
- Are errors handled?
- Is loading handled?
- Is empty data handled?
- Did injected all the dependencies on `main.go` file?
- Did you preserve existing architecture?
- Did you create reusable components instead of copy/paste?
- Does the result look like a real administrative system?

## Final expectation

The result should look like a serious internal system, not a default demo app.

Prioritize:

1. Polished admin UI.
2. Maintainable Go code.
3. Reusable Fyne components.
4. Clean service integration.
5. Clear user feedback.
6. Good report screens.
7. Consistent theme and spacing.
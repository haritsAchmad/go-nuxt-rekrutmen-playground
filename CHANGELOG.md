# Changelog

This changelog is based on the Git history from 2026-06-18 through
2026-06-30. No project changes were recorded from 2026-06-18 through
2026-06-22; development recorded in this repository started on 2026-06-23.

## 2026-06-23

### Backend

- Initialize the Go backend module and HTTP server.
- Add a health-check endpoint.
- Add the initial Lowongan list and dummy CRUD endpoints.
- Separate Lowongan logic into domain, handler, use case, and repository layers.
- Add an in-memory Lowongan repository with dummy data.
- Move HTTP route registration into the internal route package.

### Frontend

- Initialize the Nuxt application.
- Add the initial application page, Nuxt configuration, public assets, and
  TypeScript configuration.

## 2026-06-24

### Backend

- Add CORS and preflight request handling for the Nuxt application.
- Complete Lowongan create, update, delete, and status update operations.
- Add keyword search and status filtering.
- Add reusable JSON response and query ID request helpers.
- Add bulk status update and bulk delete endpoints.

### Frontend

- Connect the Nuxt application to the Lowongan API and display Lowongan data.
- Add create, edit, delete, and status update interactions.
- Add keyword search, status filtering, apply, and reset controls.
- Add row selection, select-all, bulk status update, and bulk delete actions.

## 2026-06-25

### Backend

- Add the Lowongan detail repository operation, use case, handler, and route.

### Frontend

- Add the Lowongan detail view.
- Add and refine the reusable Lowongan table header component.
- Add a reusable button component with label, type, and color variants.
- Add a reusable input component with `v-model`, configurable input types,
  file handling, validation attributes, and numeric input restrictions.

## 2026-06-26

### Backend

- Replace manually declared dummy records with generated in-memory Lowongan
  data.
- Generate 1,000 dummy Lowongan records.
- Add page and limit request parameters.
- Add paginated Lowongan responses with page, limit, total, and total-page
  metadata.
- Implement in-memory offset and limit pagination.

### Frontend

- Add success, error, and loading feedback for bulk actions.
- Add server-driven pagination, previous and next page controls, page-size
  selection, and visible record ranges.
- Refine bulk action buttons and Lowongan data presentation.

## 2026-06-29

### Backend

- Consolidate request parsing and JSON response helpers into the Lowongan
  handler.
- Remove the separate request and response packages.
- Add environment configuration helpers with string, integer, and boolean
  fallback support.
- Add application and PostgreSQL configuration, DSN generation, and a
  configurable HTTP port.
- Add `.env` dependency support, an environment example, and rules to keep the
  local `.env` file out of Git.
- Add PostgreSQL 18 connectivity through `pgx/v5` and a connection pool.
- Add a PostgreSQL Lowongan repository.
- Implement SQL `LIMIT` and `OFFSET` pagination with total-record counting.
- Add a startup repository check that reads paginated dummy data from
  PostgreSQL.
- Run Go module tidying, remove unused database dependencies, and normalize Go
  source formatting.

### Frontend

- No frontend changes were recorded.

## 2026-06-30

### Backend

- Wire the Lowongan use case and handler through repository interfaces and
  dependency injection.
- Use the PostgreSQL Lowongan repository in the API.
- Add reusable HTTP middleware for CORS, request logging, panic recovery, JWT
  authentication, role checks, and method-based role rules.
- Apply global middleware and protect Lowongan routes with authentication and
  role authorization.
- Add auth domain models, token validation, request context helpers, use case,
  PostgreSQL repository, handler, routes, and configuration.
- Add users table seed SQL and an `AUTH_SECRET` environment example.
- Allow the `Authorization` header in CORS responses.
- Add a backend devtool command for creating users.
- Extend Lowongan records and create/update requests with opening date, closing
  date, and description fields.
- Persist the new Lowongan detail fields in PostgreSQL and add the corresponding
  database migration.
- Add basic date parsing and validation, including rejecting closing dates that
  precede opening dates.

### Frontend

- Add a cookie-based login flow to the Nuxt app.
- Send the auth token from Nuxt API requests.
- Enable Pinia and move authentication state, session restoration, login,
  logout, and authorization headers into a dedicated auth store.
- Split the application into Nuxt pages, add a dedicated login page, and protect
  the Lowongan page with route middleware.
- Add opening date, closing date, description, and editable status fields to the
  Lowongan form and detail views.
- Add client-side date-range validation and success or error feedback for
  Lowongan actions and filter changes.

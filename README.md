# agree

Using simple comment strings, _agree_ can track what
your schemas are and then check for any changes that might need to be echoed in
other spots.

For example, you might have:

- FastAPI backend
  - SQLAlchemy ORM
  - Pydantic for validation / DTOs
- NextJS frontend
  - Typescript / zod for validation / contracts from api

In an ideal world, you would have a single source of truth ie. ur models and
generate OpenAPI specs from that to get your other schemas, but this isn't always
possible. In this case, you need something that enables manual tracking in multiple places
and perhaps across multiple languages.

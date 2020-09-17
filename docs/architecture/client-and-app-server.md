---
title: Client and App server
parent: Architecture
nav_order: 2

---

# Client and App server

The client is what the user sees. It's the frontend of the application.
Client interacts with the app server to do tasks such as:

- Register or Login
- Create checks and status pages
- Update incident status
- Add alerts to checks

This happens through the REST API provided by the app server. The only
job of App server is to do all the CRUD operations of the application.
After receiving a request, it authenticates, authorizes and updates the
app database.


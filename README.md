## go-user-tokens

go-user-tokens is intended to serve as a modifiable base of boilerplate code for web application servers. Specifically, go-user-tokens implements a basic token creating and authorizing service for applications which need user-managed sessions.

go-user-tokens uses MongoDB as a permanent data store and hashes user's password with SHA1. go-user-tokens has predefined endpoints for login and register behaviors. An example /authEndpoint exists to get you started. Static files can also be served using serveFiles.

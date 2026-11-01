# Permit

**Permit is a lightweight, self-hosted temporary access management service that provides secure, time-limited access to protected HTTP resources using single-use tokens.**

It allows you to grant temporary access to specific services without exposing your infrastructure globally, distributing permanent credentials, or creating additional user accounts.

For example, you can use Permit to provide temporary access to a VPN configuration download, an internal dashboard, an administrative interface, or any other HTTP resource that should only be accessible to authorized users.

Permit integrates with Nginx through its `auth_request` module, allowing access to existing applications to be controlled without modifying those applications.

## How it works

Permit uses a two-stage authorization mechanism:

1. **Create a permit.** An administrator generates a temporary access link specifying the target website, resource location, token expiration time, and session lifetime.
2. **Share the link.** The generated link can be sent to the intended user through any communication channel.
3. **Activate the permit.** The user opens the link and clicks the *Open* button. Permit validates the single-use token and creates a temporary session.
4. **Access the resource.** The user is redirected to the protected resource. Permit issues an HTTP-only session cookie.
5. **Verify requests.** Nginx sends an authorization subrequest to Permit before proxying requests to the protected application.
6. **Expire automatically.** Once the session expires, subsequent requests are denied.

The original activation token can only be used once. Permit stores SHA-256 hashes of tokens rather than their original values.

### Access flow

```text
                  Administrator
                       |
                       |
                 Create permit
                       |
                       v
                +-------------+
                |    Permit   |
                |             |
                |  Admin API  |
                +-------------+
                       |
                       |
                Temporary URL
                       |
                       v
                 End user
                       |
                       |
                 Open link
                       |
                       v
              +----------------+
              |  Embed website |
              |                |
              |  Activate     |
              |  single-use    |
              |  token         |
              +----------------+
                       |
                       |
                Session cookie
                       |
                       v
              +----------------+
              | Target website |
              |                |
              | Nginx          |
              |                |
              | auth_request   |
              +----------------+
                       |
                       | Authorization check
                       v
                +-------------+
                |   Permit    |
                |             |
                | Verify      |
                | session     |
                +-------------+
                       |
                       |
                  HTTP 200
                       |
                       v
              Protected resource
```

## Features

- **Single-use activation tokens** — cryptographically random 256-bit tokens that cannot be reused.
- **Time-limited access** — independently configurable activation token and session expiration times.
- **Session-based authorization** — activation tokens are exchanged for temporary session credentials.
- **Nginx integration** — native integration through the standard `auth_request` mechanism.
- **Multiple protected websites** — permits can be associated with different embed and target websites.
- **Resource-level redirection** — redirect users to a specific path on the target website.
- **Web administration interface** — create, inspect, and revoke permits.
- **OIDC authentication** — protect the administration interface using an external identity provider.
- **Automatic cleanup** — expired permits are periodically removed from the database.
- **PostgreSQL storage** — persistent storage of permits and session metadata.
- **Docker support** — ready-to-use container images for AMD64 and ARM64.

## Requirements

- PostgreSQL 17 or later.
- Nginx with the `ngx_http_auth_request_module` module.
- Docker (optional).
- An OpenID Connect identity provider for administrator authentication (recommended).

## Installation

### Docker

Images are published to GitHub Container Registry.

```bash
docker pull ghcr.io/zerodotfive/permit:latest
```

Supported architectures:

- `linux/amd64`
- `linux/arm64`

### Docker Compose

The repository includes a `compose.yml` file containing the application and PostgreSQL services.

Create a `.env` file:

```dotenv
# Permit
PUB_LISTEN=127.0.0.1:8090
API_LISTEN=127.0.0.1:8091

DOCS_ENABLED=false
AUTH_ENABLED=true

# PostgreSQL
POSTGRES_URI=postgres://postgres:password@127.0.0.1:5432/permit?sslmode=disable

# Database migrations
MIGRATE=true

# OpenID Connect
AUTH_CLIENT_ID=permit
AUTH_CLIENT_SECRET=your-client-secret
AUTH_REDIRECT_URL=https://permit.example1.com/auth/callback
OPENID_CONFIGURATION_ISSUER=https://auth.example.com/application/o/permit/
JWT_KEY=replace-with-a-random-secret-of-at-least-32-bytes
ADMIN_GROUP=administrators
```

Start the application:

```bash
docker compose up -d
```

The database migrations are applied automatically when `MIGRATE=true`.

**Important:** The example PostgreSQL password and JWT key must be replaced with production values. Keep the application ports private and expose them only through a trusted reverse proxy.

### Configuration

Permit is configured using environment variables.

| Variable | Default | Description |
|---|---|---|
| `PUB_LISTEN` | `127.0.0.1:8090` | Address for public activation and verification endpoints. |
| `API_LISTEN` | `127.0.0.1:8091` | Address for the administration API and web interface. |
| `PUB_LOCATION_ACTIVATE` | `/.p/a` | Activation endpoint path. |
| `PUB_LOCATION_CHAIN` | `/.p/c` | Session chaining endpoint path. |
| `PUB_LOCATION_VERIFY` | `/.p/v` | Session verification endpoint path. |
| `DOCS_ENABLED` | `false` | Enable OpenAPI documentation. |
| `AUTH_ENABLED` | `true` | Enable administrator authentication. |
| `POSTGRES_URI` | See source | PostgreSQL connection string. |
| `MIGRATE` | `true` | Automatically apply database migrations. |
| `AUTH_CLIENT_ID` | — | OIDC client ID. |
| `AUTH_CLIENT_SECRET` | — | OIDC client secret. |
| `AUTH_REDIRECT_URL` | — | OIDC callback URL. |
| `OPENID_CONFIGURATION_ISSUER` | — | OIDC discovery issuer URL. |
| `JWT_KEY` | — | JWT signing key (minimum 32 bytes). |
| `ADMIN_GROUP` | — | Group allowed to administer permits. |

The public and administration listeners are separate:

- `PUB_LISTEN` handles activation, session chaining, and session verification.
- `API_LISTEN` handles administration, authentication, and the web interface.

The public listener should only be reachable through your reverse proxy.

## Nginx configuration

Permit is designed to work with Nginx's `auth_request` module.

The following example demonstrates three virtual hosts:

| Hostname | Purpose |
|---|---|
| `permit.example1.com` | Permit administration interface and API. |
| `authprovider.example2.com` | Website used to embed and activate temporary access links. |
| `authtarget.example3.com` | Protected application or resource. |

The embed and target websites can be different domains.

### 1. Permit administration

The administration interface is available on port `8091`.

It is recommended to restrict access to the administration interface using an IP allowlist in addition to the application's OIDC authentication.

```nginx
server {
    listen 80;
    server_name permit.example1.com;

    return 301 https://$host$request_uri;
}

server {
    listen 443 ssl;
    server_name permit.example1.com;

    ssl_certificate /etc/ssl/certs/permit.crt;
    ssl_certificate_key /etc/ssl/private/permit.key;

    location / {
        allow 192.0.2.10;
        deny all;

        proxy_pass http://127.0.0.1:8091;
        proxy_set_header Host $http_host;
        proxy_set_header X-Forwarded-For $remote_addr;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

Replace `192.0.2.10` with your administrator's public IP address.

### 2. Embed website

This website is used to activate temporary access links.

```nginx
server {
    listen 443 ssl;
    server_name authprovider.example2.com;

    server_tokens off;

    ssl_certificate /etc/ssl/certs/authprovider.crt;
    ssl_certificate_key /etc/ssl/private/authprovider.key;

    location / {
        satisfy any;

        allow 192.0.2.10;
        deny all;

        auth_request /.p/v;

        # Your application upstream
        proxy_pass http://127.0.0.1:9000;

        proxy_set_header Host $http_host;
        proxy_set_header X-Forwarded-For $remote_addr;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # Permit public endpoints
    location ^~ /.p/ {
        proxy_pass http://127.0.0.1:8090;

        proxy_pass_request_body off;
        proxy_set_header Content-Length "";

        proxy_set_header Host $http_host;
        proxy_set_header Referer $http_referer;
        proxy_set_header X-Forwarded-For $remote_addr;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### 3. Target website

This is the actual protected application or resource.

```nginx
server {
    listen 443 ssl;
    server_name authtarget.example3.com;

    server_tokens off;

    ssl_certificate /etc/ssl/certs/authtarget.crt;
    ssl_certificate_key /etc/ssl/private/authtarget.key;

    location / {
        satisfy any;

        allow 192.0.2.10;
        deny all;

        auth_request /.p/v;

        # Your application upstream
        proxy_pass http://127.0.0.1:9001;

        proxy_set_header Host $http_host;
        proxy_set_header X-Forwarded-For $remote_addr;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # Permit public endpoints
    location ^~ /.p/ {
        proxy_pass http://127.0.0.1:8090;

        proxy_pass_request_body off;
        proxy_set_header Content-Length "";

        proxy_set_header Host $http_host;
        proxy_set_header Referer $http_referer;
        proxy_set_header X-Forwarded-For $remote_addr;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### How Nginx authorization works

The important directive is:

```nginx
auth_request /.p/v;
```

Before forwarding a request to the application, Nginx sends an internal authorization subrequest to Permit.

Permit checks the session cookie and verifies that:

- The session exists.
- The session has not expired.
- The requested hostname matches the website associated with the permit.

The response determines whether Nginx forwards the original request:

| Permit response | Nginx behavior |
|---|---|
| `200 OK` | Allow the request. |
| `401 Unauthorized` | Deny the request. |
| `403 Forbidden` | Deny the request. |
| Other errors | Return an authorization error. |

The following directives are also important:

```nginx
satisfy any;

allow 192.0.2.10;
deny all;
```

They allow explicitly whitelisted IP addresses to bypass the authorization subrequest. All other clients must have a valid Permit session.

If you do not need IP-based exceptions, you can remove `satisfy`, `allow`, and `deny`, leaving `auth_request` as the only access control mechanism.

**Important Nginx considerations:**

- The `/.p/` location must be accessible without the normal `auth_request` protection, otherwise authorization will recursively call itself.
- `proxy_pass` must not append a replacement URI. The configuration above preserves the original Permit endpoint path.
- `proxy_pass_request_body off` and an empty `Content-Length` prevent the original request body from being sent to the authorization service.
- `$remote_addr` is used instead of forwarding an arbitrary client-supplied `X-Forwarded-For` value.
- The `Host` header must be preserved because Permit validates the requested hostname against the permit configuration.
- Do not expose ports `8090` and `8091` directly to the Internet.

After configuring Nginx, validate and reload:

```bash
nginx -t && systemctl reload nginx
```

## Creating permits

Open the administration interface:

```text
https://permit.example1.com/
```

After administrator authentication, navigate to the Permits page.

When creating a permit, the following parameters are available:

| Parameter | Description |
|---|---|
| Comment | Human-readable description of the permit. |
| Embed site | Website where the activation link will be opened. |
| Target site | Website containing the protected resource. |
| Target location | Optional path to redirect the user to after activation. |
| Token lifetime | How long the activation link remains valid (1–1440 minutes). |
| Session lifetime | How long the activated session remains valid (1–1440 minutes). |

After creation, Permit generates a unique activation URL.

The URL is displayed only once, so make sure to copy and securely distribute it.

Deleting a permit also deletes its associated sessions and chained tokens.

## Security model

Permit is designed to minimize the exposure of internal services and reduce the lifetime of temporary access credentials.

### Token security

- Tokens contain 256 bits of cryptographically secure random data.
- Only SHA-256 hashes are stored in PostgreSQL.
- Activation tokens are single-use.
- Session credentials are generated independently of activation tokens.
- Session cookies are marked `HttpOnly` and use `SameSite=Lax`.

### Session expiration

Each permit has two independent expiration periods:

1. Activation token lifetime.
2. Session lifetime, starting when the token is activated.

Expired records are periodically removed from the database.

### Network security

- Keep both Permit listeners bound to localhost or a private network.
- Expose only the required Nginx virtual hosts.
- Use HTTPS for all externally accessible websites.
- Restrict the administration interface to trusted administrators.
- Do not trust client-supplied forwarding headers. Ensure the reverse proxy overwrites them.

**Note:** The current implementation does not yet set the cookie's `Secure` attribute. For HTTPS-only deployments, enabling `Secure: true` in the session cookie configuration is recommended.

## API

Permit exposes two separate HTTP interfaces.

### Public API

The public API is used internally by Nginx and the activation workflow.

Default listener: `127.0.0.1:8090`

| Method | Endpoint | Description |
|---|---|---|
| GET | `/.p/a/{token}` | Validate an activation token and display the activation page. |
| POST | `/.p/a/{token}` | Activate a permit and create a session. |
| GET | `/.p/c/{token}` | Exchange a chained token for a session on the target website. |
| GET | `/.p/v` | Verify the current session. |

### Administration API

Default listener: `127.0.0.1:8091`

API prefix: `/api/v1`

| Method | Endpoint | Description |
|---|---|---|
| GET | `/api/v1/permits` | List permits. |
| POST | `/api/v1/permits` | Create a permit. |
| DELETE | `/api/v1/permits/{id}` | Delete a permit. |

The administration API is protected by OIDC authentication when `AUTH_ENABLED=true`.

OpenAPI documentation can be enabled using `DOCS_ENABLED=true`.

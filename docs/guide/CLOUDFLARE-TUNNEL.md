# Cloudflare Tunnel Deployment

This repository is a standard Go HTTP server, so the safest Cloudflare deployment path is:

1. Run the FHIR engine as a normal container or server.
2. Run `cloudflared` beside it.
3. Publish a hostname from Cloudflare Tunnel to the local service.

## Recommended Hostname

Use:

```text
fhir.zarishsphere.com
```

This keeps the FHIR API separate from the root website and matches the current `zarishsphere.com` Cloudflare zone.

## Why this approach

- No inbound ports need to be opened on the origin server.
- Cloudflare sits in front of the service for TLS, proxying, and security.
- The application can stay as a normal Go server instead of being rewritten for Workers or Pages.

## Prerequisites

- A machine or VM where this repository will run continuously
- Docker and Docker Compose
- A Cloudflare Tunnel created in the `zarishsphere.com` account
- A tunnel token with permission to run the tunnel

## Repo files added for Tunnel

- Base app stack: `docker-compose.yml`
- Tunnel sidecar: `deploy/cloudflare/docker-compose.tunnel.yml`
- Local secret template: `.env.cloudflare.example`

## Dashboard steps

Because the Cloudflare API token available in this session only has DNS and zone permissions, create the tunnel in the Cloudflare dashboard:

1. Go to Cloudflare Dashboard.
2. Open `zarishsphere.com`.
3. Go to `Networks` or `Connectors` and open `Cloudflare Tunnels`.
4. Create a new remotely-managed tunnel.
5. Add a published application route:
   - Hostname: `fhir`
   - Domain: `zarishsphere.com`
   - Service: `http://fhir-engine:8080`
6. Copy the tunnel token.

## Local runtime steps

Create a local secret file:

```bash
cp .env.cloudflare.example .env.cloudflare
```

Edit `.env.cloudflare` and paste the real token:

```bash
TUNNEL_TOKEN=eyJh...
```

Start the full stack:

```bash
docker compose \
  --env-file .env.cloudflare \
  -f docker-compose.yml \
  -f deploy/cloudflare/docker-compose.tunnel.yml \
  up -d
```

## Verification

Local checks:

```bash
docker compose ps
docker compose logs -f fhir-engine
docker compose logs -f cloudflared
```

Public checks:

```bash
curl https://fhir.zarishsphere.com/health
curl https://fhir.zarishsphere.com/fhir/R5/metadata
```

## Notes

- Do not commit `.env.cloudflare`.
- A remotely-managed tunnel stores routing rules in Cloudflare, so the repo only needs the runtime token.
- If you later move this service to a VM or container host, reuse the same tunnel token or rotate it in Cloudflare.

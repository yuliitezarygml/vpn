# vpn

Hiddify-based VPN workspace: client app + server panel (local forks).

## Layout

- `client/` — Flutter app, core, sing-box, geo, helpers
- `server/` — Hiddify-Manager, panel, fronts, relay

## Client (macOS)

```bash
export PATH="$HOME/flutter-3.38.5/bin:$PATH"
cd client/hiddify-app
make build-macos-libs
flutter run -d macos
```

Use Flutter **3.38.5** (not 3.47).

## Server panel (Docker)

```bash
cd server/Hiddify-Manager
# set REDIS_PASSWORD / MYSQL_PASSWORD in docker.env and .env
docker compose up -d --build
```

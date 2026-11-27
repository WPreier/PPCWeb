# Security-Header & Maßnahmen (statische Seiten)

Stand: PPCWeb **0.1.1.8**

Statische HTML-Seiten können Security-Header **nur eingeschränkt** selbst setzen. Wirkungsvoll sind:

1. **HTTP-Header** vom Server (Apache `.htaccess`, SeDat, CDN `_headers`)  
2. **Meta-Tags** in HTML (CSP/Referrer – mit Einschränkungen)  
3. **Inhaltliche Maßnahmen** (kein Inline-Theme mehr, `noopener noreferrer`, …)

---

## In jeder HTML-Seite (`<head>`)

| Maßnahme | Umsetzung |
|---|---|
| Referrer-Policy | `<meta name="referrer" content="strict-origin-when-cross-origin">` |
| Content-Security-Policy | `<meta http-equiv="Content-Security-Policy" content="…">` |
| Theme ohne Inline | `assets/js/theme-init.js` möglichst früh |
| Externe Links | `rel="noopener noreferrer"` |

**Hinweis:** `X-Frame-Options`, `HSTS`, `Permissions-Policy`, `COOP/CORP`, `frame-ancestors` wirken **nur als HTTP-Header**, nicht zuverlässig als Meta.

Aktuelle CSP erlaubt bewusst `'unsafe-inline'` für Scripts (JSON-LD auf der Startseite) und Styles (vorhandene `style="…"`-Attribute) sowie Bunny Fonts.

---

## Server: `theweb/.htaccess` (Apache)

Gesetzt u. a.:

- `X-Content-Type-Options: nosniff`
- `X-Frame-Options: DENY`
- `Referrer-Policy: strict-origin-when-cross-origin`
- `Permissions-Policy: …` (Sensoren/Kamera/Mikrofon/… aus)
- `Cross-Origin-Opener-Policy: same-origin`
- `Cross-Origin-Resource-Policy: same-site`
- `X-Permitted-Cross-Domain-Policies: none`
- `X-DNS-Prefetch-Control: off`
- `X-XSS-Protection: 0` (Legacy-Filter absichtlich aus)
- `Content-Security-Policy: …` inkl. `frame-ancestors 'none'`
- `Strict-Transport-Security` nur bei HTTPS (`env=HTTPS`)
- `Options -Indexes`
- `ErrorDocument 404`
- Block von versehentlichen Config/Secret-Dateien

## CDN/Netlify-ähnlich: `theweb/_headers`

Gleiche Kern-Header für Hosts, die `_headers` unterstützen.

## SeDat / Seda

HTTP-Header idealerweise in **SeDat** setzen (einheitlich für alle Responses). Die Meta-Tags greifen auch, wenn der Proxy noch keine Header injiziert. `.htaccess` greift nur hinter Apache.

---

## Prüfen

```bash
curl -sSI https://ppc.at/ | findstr /I "content-security x-frame referrer permissions strict-transport x-content"
```

Browser: DevTools → Network → Document → Response Headers.

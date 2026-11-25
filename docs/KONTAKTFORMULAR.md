# Kontaktformular PPCWeb – Einrichtung, Tests, Sicherheit

Stand: **2026-11-25** · PPCWeb **0.1.1.5**  
Serverziel: **kmbhelp** `192.168.36.55` (Seda/SeDat + M-Con + ppc-contact)  
Mail Phase 1: **`192.168.36.251:25`** ohne Authentifizierung (über lokalen Forwarder)

---

## 1. Architektur (Überblick)

```
Internet
   │
   ▼
SeDat (TLS, Reverse-Proxy auf kmbhelp)
   ├── /*            → Seda (statisches theweb/ in RAM)
   └── /api/*        → ppc-contact 127.0.0.1:8690
                          │
                          ├─ ALTCHA Challenge / Verify (lokal, kein Tracking)
                          └─ HMAC → M-Con 127.0.0.1:8687
                                       │
                                       └─ SMTP 127.0.0.1:2525
                                            └─ socat → 192.168.36.251:25
```

| Komponente | Rolle |
|---|---|
| **Seda** | Liefert die Website (HTML/CSS/JS) |
| **SeDat** | TLS + Routing `/api` → Bridge |
| **ppc-contact** | Formular-API, Captcha, Rate-Limit, signiert M-Con |
| **M-Con** | Queue, Spam-Check, Versand, Archiv |
| **smtp-forward** | Loopback-SMTP, weil M-Con Klartext-SMTP nur auf Loopback erlaubt |

**Warum eine Bridge?**  
M-Con verlangt HMAC-Signaturen. Das Secret darf **nie** im Browser liegen. Deshalb signiert nur `ppc-contact` auf dem Server.

---

## 2. Captcha-Entscheidung: ALTCHA (selbst gehostet)

| Option | Bewertung für ppc.at |
|---|---|
| Google reCAPTCHA | ❌ Tracking, US-Transfer, Cookie-Banner |
| Cloudflare Turnstile | ⚠️ Besser, aber Drittanbieter |
| hCaptcha | ⚠️ Drittanbieter |
| **ALTCHA (Open Source, self-hosted)** | ✅ DSGVO-freundlich, kein Cookie, WCAG, PoW gegen Bots |

**Wahl:** ALTCHA v1 Proof-of-Work, Challenge von `ppc-contact`, Widget selbst unter `theweb/assets/js/altcha.min.js`.

Zusätzlich (Defense in Depth):

- Honeypot-Feld `website` (für Bots unsichtbar)
- Mindestzeit 3 s zwischen Laden und Absenden
- Rate-Limit pro IP
- Origin/Referer-Prüfung
- M-Con Spam-Modul (Ollama optional)

---

## 3. Dateien in diesem Repo

| Pfad | Zweck |
|---|---|
| `theweb/kontakt.html` | Formular-UI |
| `theweb/assets/js/contact.js` | Submit → `/api/contact` |
| `theweb/assets/js/altcha.min.js` | ALTCHA Widget (vendored) |
| `server/contact-bridge/` | Go-Dienst `ppc-contact` |
| `server/deploy/*.service` | systemd Units |
| `server/deploy/mcon-ppc.yaml.example` | M-Con-Ausschnitt Phase 1 |
| `docs/KONTAKTFORMULAR.md` | Diese Anleitung |

---

## 4. Server einrichten (kmbhelp) – Reihenfolge

> SSH von dieser Entwickler-Maschine auf `kmbhelp@192.168.36.55` war hier **nicht** möglich (kein Public-Key). Bitte die Schritte lokal auf dem Server ausführen.

### 4.1 Voraussetzungen

```bash
# auf kmbhelp
sudo apt update   # bzw. dnf
sudo apt install -y golang-go socat git   # Go 1.23+ empfohlen
```

M-Con und Seda/SeDat sollten bereits laufen oder nach deren eigener Doku installiert werden:

- https://github.com/WPreier/M-Con → `docs/INSTALL.md`
- https://github.com/WPreier/seda · https://github.com/WPreier/sedat

### 4.2 SMTP-Forwarder (Phase 1)

M-Con-Regel: `tls_mode: none` nur zu `127.0.0.1` / `localhost`.

```bash
sudo cp /pfad/zum/PPCWeb/server/deploy/smtp-forward.service /etc/systemd/system/
# Host/Port in der Unit prüfen: 192.168.36.251:25
sudo systemctl daemon-reload
sudo systemctl enable --now smtp-forward.service
sudo systemctl status smtp-forward.service
```

**Einzeln testen:**

```bash
# Von kmbhelp aus – sollte SMTP-Banner vom Relay zeigen
printf 'QUIT\r\n' | nc -v 127.0.0.1 2525
printf 'QUIT\r\n' | nc -v 192.168.36.251 25
```

Firewall: nur kmbhelp darf Port 25 auf `.251` erreichen (Host-Firewall / Router-ACL).

### 4.3 M-Con für PPC konfigurieren

```bash
sudo install -d -o root -g mcon -m 0750 /etc/mcon
# mcon.yaml aus M-Con-Beispiel + server/deploy/mcon-ppc.yaml.example zusammenführen
# secrets.yaml aus server/deploy/mcon-secrets.ppc.example.yaml
sudo chmod 0440 /etc/mcon/secrets.yaml
sudo chown root:mcon /etc/mcon/secrets.yaml

# Secrets erzeugen, z. B.:
openssl rand -hex 24   # admin.token
openssl rand -hex 24   # clients.ppcweb.hmac_secret

sudo systemctl restart mcon
curl -sS http://127.0.0.1:8688/healthz   # Admin-Health laut M-Con-Doku
```

Wichtig in `mcon.yaml`:

- `providers.smtp-local.host: 127.0.0.1`, `port: 2525`, `tls_mode: none`
- Client-ID `ppcweb` mit `allowed_from: [noreply@ppc.at]`
- `server.public_addr` / `admin_addr` nur Loopback

**Einzeln testen (M-Con Test-UI / mcon-ctl):**

```bash
# Admin-Token aus secrets.yaml
curl -sS -H "Authorization: Bearer <admin-token>" \
  http://127.0.0.1:8688/admin/v1/test/send \
  -H 'Content-Type: application/json' \
  -d '{"from":"noreply@ppc.at","to":["preier@ppc.at"],"subject":"M-Con Test","text":"Hallo"}'
```

Oder `mcon-ctl` laut M-Con `docs/TESTING.md`.

### 4.4 ppc-contact bauen und starten

```bash
cd /pfad/zum/PPCWeb/server/contact-bridge
go mod tidy
go build -o ppc-contact .

sudo useradd --system --home-dir /nonexistent --shell /usr/sbin/nologin ppc-contact 2>/dev/null || true
sudo install -d -o root -g ppc-contact -m 0750 /etc/ppc-contact
sudo install -m 0755 ppc-contact /usr/local/bin/ppc-contact
sudo cp config.example.yaml /etc/ppc-contact/config.yaml
sudo chmod 0640 /etc/ppc-contact/config.yaml
sudo chown root:ppc-contact /etc/ppc-contact/config.yaml

# In config.yaml setzen:
# - altcha_secret (openssl rand -hex 24)
# - mcon.hmac_secret = dasselbe wie clients.ppcweb.hmac_secret
# - mcon.key_id = ppcweb
# - mcon.from / to
# - allowed_origins https://ppc.at

sudo cp ../../server/deploy/ppc-contact.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable --now ppc-contact
curl -sS http://127.0.0.1:8690/healthz
curl -sS http://127.0.0.1:8690/api/altcha/challenge
```

### 4.5 Website nach Seda laden

Inhalt von `theweb/` in das von Seda konfigurierte Root legen und Seda reloaden (`seda` Admin-Socket / SIGHUP laut Seda-Doku).

### 4.6 SeDat: `/api` routen

In der SeDat-Konfiguration Path-Prefix `/api/` auf `http://127.0.0.1:8690` legen (siehe `server/deploy/sedat-ppc-snippet.toml` und SeDat `docs/CONFIG.md`). Danach SeDat reload.

**Gesamttest von außen:**

1. https://ppc.at/kontakt.html öffnen  
2. Formular ausfüllen, ALTCHA abwarten, absenden  
3. Mail an `preier@ppc.at` prüfen  
4. M-Con Admin: Message-Status `sent`

---

## 5. Testplan (einzeln → gesamt)

| # | Test | Erwartung |
|---|---|---|
| T1 | `nc 192.168.36.251 25` von kmbhelp | Banner, QUIT ok |
| T2 | `nc 127.0.0.1 2525` | gleiches Banner (Forwarder) |
| T3 | M-Con test/send | 202 + State wird `sent` |
| T4 | `curl …/api/altcha/challenge` | JSON mit algorithm/challenge/signature |
| T5 | POST `/api/contact` ohne Captcha | 400 `captcha_*` |
| T6 | POST mit gefülltem Honeypot | 202, **keine** Mail |
| T7 | Browser-Formular gültig | 202, Mail kommt an |
| T8 | 10× schnell absenden | irgendwann 429 |
| T9 | Absenden in &lt;3 s | 400 `too_fast` |

---

## 6. Versand-Stufen: von lokalem SMTP bis Microsoft 365 OAuth

Alle Stufen ändern **nur M-Con** (`mcon.yaml` + `secrets.yaml`). Frontend und ppc-contact bleiben gleich (`provider` leer = Default).

### Stufe A – Phase 1 (jetzt): Internes SMTP ohne Auth

| Wo | Was |
|---|---|
| `smtp-forward.service` | `127.0.0.1:2525` → `192.168.36.251:25` |
| `mcon.yaml` → `providers.smtp-local` | host `127.0.0.1`, port `2525`, `tls_mode: none` |
| Secrets | `password: ""` |

Lernen: Envelope, Relay-ACL, Open-Relay vermeiden (nur kmbhelp erlauben).

### Stufe B – SMTP mit STARTTLS + Benutzer/Passwort

| Wo | Was |
|---|---|
| `providers.smtp-relay` | host z. B. `smtp.office365.com`, port `587`, `tls_mode: starttls`, `username` |
| `secrets.providers.smtp-relay.password` | App-Passwort / SMTP-Passwort |
| `providers.default` | auf `smtp-relay` stellen **oder** Message `provider: "smtp-relay"` |

Lernen: STARTTLS vs. implicit TLS (465), Authenticated Relay, SPF.

### Stufe C – Transaktionale HTTP-APIs (SendGrid, Postmark, Brevo, Resend, Mailgun, SES)

| Wo | Was |
|---|---|
| `providers.<name>` | `type: sendgrid` / `postmark` / … laut M-Con-Beispiel |
| `secrets.providers.<name>` | API-Key / Token |
| Default oder Fallback | z. B. API primär, `smtp-local` als `fallback` |

Lernen: API-Keys rotieren, Webhooks (Delivery), Kostenlimits.

### Stufe D – Microsoft 365 / Exchange Online mit OAuth2 (Graph)

| Wo | Was |
|---|---|
| Entra ID (Azure AD) | App-Registrierung, Application permission `Mail.Send`, Admin-Consent |
| `providers.graph` | `type: graph`, `send_as: noreply@ppc.at` |
| `secrets.providers.graph` | `tenant_id`, `client_id`, `client_secret` |
| `providers.default` | `graph` |

Schritte zum Lernen:

1. Portal https://entra.microsoft.com → App registrieren  
2. API-Berechtigung Microsoft Graph **Application**: `Mail.Send`  
3. Admin-Zustimmung erteilen  
4. Client Secret erzeugen (Ablauf notieren)  
5. Mailbox/`send_as` muss der App das Senden erlauben (Exchange Application Access Policy empfohlen)  
6. In M-Con Secrets eintragen, `systemctl reload` bzw. Admin-Reload  
7. Testsend über `mcon-ctl` / Admin test/send  

Dokumentation Microsoft: „Get access without a user“ / Client Credentials + Graph `sendMail`.

### Stufe E – Gmail OAuth (optional)

Analog `providers.gmail` + Refresh-Token in Secrets (siehe M-Con `configs/secrets.yaml.example`).

---

## 7. Sicherheit: Einbruch, Ausbruch, Spam

### 7.1 Einbruchsschutz (Angriff von außen)

| Maßnahme | Umsetzung |
|---|---|
| Keine öffentlichen Admin-Ports | M-Con `:8687/:8688`, ppc-contact `:8690` nur `127.0.0.1` |
| TLS nur an SeDat | ACME/Zertifikat dort |
| Path-Routing | Nur `/api/altcha/*` und `/api/contact` freigeben |
| HMAC | Secrets ≥16 Zeichen, nicht im Git |
| Body-Limit | Bridge 64 KiB; SeDat zusätzlich begrenzen |
| OS-Hardening | systemd `ProtectSystem`, `NoNewPrivileges` (Units mitliefern) |
| Updates | Go-Binary neu bauen bei CVEs; OS-Patches |

### 7.2 Ausbruchsschutz (falls ein Dienst kompromittiert)

| Maßnahme | Umsetzung |
|---|---|
| Getrennte User | `mcon`, `ppc-contact`, `nobody` (socat) |
| Minimale Capabilities | Units ohne Caps |
| Kein Shell-Login für Dienstuser | `nologin` |
| Secrets-Rechte | `0440 root:gruppe` |
| Netzwerk | Mail-Relay `.251` nur von kmbhelp; kein Open Relay |
| Archive | M-Con Archive nur lokal, Backup verschlüsselt offsite |

### 7.3 Spam- & Missbrauchsschutz

| Schicht | Was |
|---|---|
| ALTCHA | PoW, ohne Tracking |
| Honeypot | stille Annahme |
| Timing | &lt;3 s verwerfen |
| Rate-Limit | pro IP in ppc-contact + M-Con client rate_limit |
| M-Con Spam | Ollama/OpenRouter optional |
| Dedupe / Idempotenz | M-Con eingebaut |
| Inhalt | keine Attachments vom Formular |

### 7.4 Datenschutz

- Zustimmung Checkbox + Hinweis Impressum  
- Speicherung nur zur Anfragebearbeitung (M-Con Queue/Archiv – Aufbewahrung festlegen)  
- ALTCHA speichert keine Personenprofile  
- Theme weiter nur `localStorage`

---

## 8. Betrieb & Fehlerbilder

| Symptom | Prüfen |
|---|---|
| Captcha „Prüfung fehlgeschlagen“ | **Häufigste Ursache:** `/api/altcha/challenge` kommt nicht bei `ppc-contact` an (SeDat-Route fehlt → Seda liefert HTML-404). Siehe Diagnose unten. |
| Captcha lädt nicht | SeDat `/api` Route, `curl challenge`, Browser-Konsole (Network-Tab) |
| `send_failed` | `journalctl -u mcon -u ppc-contact -u smtp-forward` |
| Mail kommt nicht | Relay `.251` Logs, SPF/From, M-Con Message-State |
| 403 forbidden | `allowed_origins` in ppc-contact config |
| 429 | Rate-Limit – warten oder Limit anpassen |

```bash
journalctl -u ppc-contact -u mcon -u smtp-forward -f
```

### Diagnose: ALTCHA „Prüfung fehlgeschlagen“

Im Browser (F12 → Netzwerk) den Request auf `/api/altcha/challenge` ansehen:

| Antwort | Bedeutung | Fix |
|---|---|---|
| **404** + HTML | SeDat routet `/api` nicht → Seda antwortet | In SeDat Path `/api/` → `http://127.0.0.1:8690` |
| **Connection failed** | ppc-contact down | `systemctl status ppc-contact` / enable --now |
| **403** | Origin-Block (ältere Bridge) | Bridge ≥ 0.1.1.6; `allowed_origins` prüfend |
| **200** + JSON mit `algorithm`, `challenge`, `salt`, `signature`, `maxNumber` | API ok | Widget neu laden; Cache leeren |

Auf dem Server:

```bash
curl -sS -D- http://127.0.0.1:8690/healthz
curl -sS http://127.0.0.1:8690/api/altcha/challenge | head -c 400
# Erwartung: JSON, kein HTML

# Von außen über die Domain (nach SeDat-Route):
curl -sS -D- https://ppc.at/api/altcha/challenge | head -c 400
```

---

## 9. Checkliste Go-Live

- [ ] smtp-forward aktiv, T1/T2 ok  
- [ ] M-Con mit Client `ppcweb`, Testsend ok  
- [ ] ppc-contact health + challenge ok  
- [ ] SeDat routet `/api`  
- [ ] `theweb/` in Seda aktualisiert (0.1.1.5)  
- [ ] Browser-Test T7  
- [ ] Impressum/Datenschutz gelesen und korrekt  
- [ ] Secrets **nicht** in Git, Backup der Config offline  

---

## 10. Was Sie als Nächstes tun (kurz)

1. Auf **kmbhelp** einloggen (SSH-Key einrichten, falls nötig).  
2. Abschnitte **4.2 → 4.6** der Reihe nach umsetzen.  
3. Testplan **Abschnitt 5** abhaken.  
4. Später Versandstufe B–D nach Bedarf (Abschnitt 6).  

Bei Fragen zu M-Con-API/Config: Upstream-Doku in https://github.com/WPreier/M-Con (`docs/API.md`, `CONFIG.md`, `INSTALL.md`, `TESTING.md`).

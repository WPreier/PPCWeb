# PPCWeb auf einem neuen Server installieren und updaten

Stand: **2026-11-27** · Quelle: https://github.com/WPreier/PPCWeb  
Aktuelle Version: Datei [`VERSION`](../VERSION) bzw. Git-Tag `v0.1.1.x`

Diese Anleitung gilt für **jeden neuen Linux-Server** (nicht nur kmbhelp).  
Spezialdokus:

| Thema | Datei |
|---|---|
| Kontaktformular-Architektur & Tests | [`KONTAKTFORMULAR.md`](KONTAKTFORMULAR.md) |
| Nur `ppc-contact` bauen/installieren | [`INSTALL_PPC_CONTACT.md`](INSTALL_PPC_CONTACT.md) |
| Security-Header | [`SECURITY_HEADERS.md`](SECURITY_HEADERS.md) |
| Projektpflege / Versionierung | [`../DOKUMENTATION.md`](../DOKUMENTATION.md) |

---

## 0. Was wird installiert?

PPCWeb hat zwei Schichten. Auf einem neuen Server entscheiden Sie, welche Sie brauchen:

| Variante | Inhalt | Wann |
|---|---|---|
| **A – nur Website** | Ordner `theweb/` (HTML/CSS/JS) | Statischer Host, Apache/Nginx, CDN, Panel |
| **B – Vollstack (ppc.at)** | Seda + SeDat + M-Con + smtp-forward + ppc-contact + `theweb/` | Eigener Server mit Kontaktformular und Mail |

```
Internet
   │
   ▼
SeDat (TLS, Reverse-Proxy)          ← Variante B
   ├── /*      → Seda (theweb/ in RAM)
   └── /api/*  → ppc-contact :8690
                    └─ HMAC → M-Con :8687 → SMTP (lokal/Relay)
```

**Variante A** braucht weder Go noch M-Con – nur den Inhalt von `theweb/`.

---

## 1. Voraussetzungen (neuer Server)

### Allgemein

- Linux (Debian/Ubuntu oder Alma/RHEL)
- SSH-Zugang mit `sudo`
- Öffentliche DNS-Namen (z. B. `ppc.at`, `www.ppc.at`) zeigen auf den Server
- Ausgehendes HTTPS (GitHub, Go-Module, ggf. ACME)

### Variante B zusätzlich

| Paket / Dienst | Zweck |
|---|---|
| `git`, `curl`, Build-Tools | Repo klonen, Go-Binaries bauen |
| Go **1.23+** | `ppc-contact` (und ggf. M-Con) bauen |
| **Seda** | Statische Website ausliefern |
| **SeDat** | TLS + Routing `/api/` |
| **M-Con** | Mail-Queue / Versand |
| `socat` | SMTP-Forwarder (Phase 1) |

Upstream-Install:

- https://github.com/WPreier/seda  
- https://github.com/WPreier/sedat  
- https://github.com/WPreier/M-Con → `docs/INSTALL.md`

---

## 2. Neuinstallation – Variante A (nur Website)

### 2.1 Repo holen (lokal oder auf dem Server)

```bash
# auf dem Server (empfohlen)
sudo mkdir -p /opt/ppcweb
sudo chown "$USER":"$USER" /opt/ppcweb
git clone https://github.com/WPreier/PPCWeb.git /opt/ppcweb
cd /opt/ppcweb
git checkout v0.1.1.8    # feste Version; Tag an VERSION anpassen
```

Oder lokal auf Windows klonen und nur `theweb/` per SFTP hochladen.

### 2.2 Document Root befüllen

**Nur den Inhalt von `theweb/`** gehört ins öffentliche Root (nicht `README.md`, `docs/`, `server/`).

```bash
# Beispiel Apache/Nginx-Root
sudo rsync -a --delete /opt/ppcweb/theweb/ /var/www/ppc.at/

# Rechte je nach Webserver-User anpassen, z. B.:
# sudo chown -R www-data:www-data /var/www/ppc.at
```

Document Root = Ordner, in dem `index.html` liegt.

### 2.3 Webserver

- **Apache:** `theweb/.htaccess` mitliefern (Security-Header, Kompression, 404).  
  `AllowOverride` für dieses Directory muss aktiv sein.
- **Nginx / CDN:** Header aus `theweb/_headers` bzw. [`SECURITY_HEADERS.md`](SECURITY_HEADERS.md) in die Server-Config übernehmen (`.htaccess` gilt dort nicht).
- **Seda:** Inhalt von `theweb/` ins konfigurierte Seda-Root legen und Seda reloaden (siehe Variante B).

### 2.4 Smoke-Test

- Startseite, Mobilmenü, Hell/Dunkel
- `https://IHRE-DOMAIN/.well-known/security.txt`
- Hart neu laden (Strg+F5)

Fertig für Variante A. Updates: Abschnitt **4**.

---

## 3. Neuinstallation – Variante B (Vollstack)

Reihenfolge einhalten. Platzhalter ersetzen:

| Platzhalter | Bedeutung |
|---|---|
| `IHR_USER` | SSH-Account auf dem Server |
| `SERVER` | Hostname oder IP |
| `/PFAD/SEDA/ROOT/` | Von Seda konfiguriertes Website-Root |
| `v0.1.1.8` | Aktueller Release-Tag |

### 3.1 SSH und Basispakete

```bash
ssh IHR_USER@SERVER
```

Debian/Ubuntu:

```bash
sudo apt update
sudo apt install -y git curl build-essential socat

# Go 1.23+ (Distro-Paket oft zu alt)
curl -fsSL https://go.dev/dl/go1.23.6.linux-amd64.tar.gz -o /tmp/go.tgz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf /tmp/go.tgz
echo 'export PATH=/usr/local/go/bin:$PATH' | sudo tee /etc/profile.d/go.sh
source /etc/profile.d/go.sh
go version
```

### 3.2 PPCWeb-Repo nach `/opt/ppcweb`

```bash
sudo mkdir -p /opt/ppcweb
sudo chown "$USER":"$USER" /opt/ppcweb
git clone https://github.com/WPreier/PPCWeb.git /opt/ppcweb
cd /opt/ppcweb
git checkout v0.1.1.8
```

### 3.3 Seda + SeDat installieren

Nach Upstream-Doku installieren und TLS für die Domain einrichten.  
SeDat soll später:

1. `/*` → Seda (statisch)  
2. `/api/` → `http://127.0.0.1:8690` (ppc-contact)  

Hinweis: [`../server/deploy/sedat-ppc-snippet.toml`](../server/deploy/sedat-ppc-snippet.toml) und SeDat `docs/CONFIG.md`.

Ports von M-Con (`8687`/`8688`) und ppc-contact (`8690`) **nicht** öffentlich freigeben – nur Loopback.

### 3.4 SMTP-Forwarder (Phase 1: internes Relay)

Nur nötig, wenn M-Con per Klartext-SMTP auf Loopback sendet und ein internes Relay (z. B. `192.168.36.251:25`) den Versand übernimmt.

```bash
sudo cp /opt/ppcweb/server/deploy/smtp-forward.service /etc/systemd/system/
# Host/Port in der Unit prüfen und anpassen
sudo systemctl daemon-reload
sudo systemctl enable --now smtp-forward.service
sudo systemctl status smtp-forward.service --no-pager
```

Test:

```bash
printf 'QUIT\r\n' | nc -v 127.0.0.1 2525
```

### 3.5 M-Con

M-Con nach Upstream installieren, dann PPC-Ausschnitte mergen:

```bash
sudo install -d -o root -g mcon -m 0750 /etc/mcon
# mcon.yaml: Beispiel aus M-Con + server/deploy/mcon-ppc.yaml.example
# secrets.yaml: aus server/deploy/mcon-secrets.ppc.example.yaml

openssl rand -hex 24   # admin.token
openssl rand -hex 24   # clients.ppcweb.hmac_secret  ← merken für ppc-contact

sudo chmod 0440 /etc/mcon/secrets.yaml
sudo chown root:mcon /etc/mcon/secrets.yaml
sudo systemctl enable --now mcon
```

Wichtig:

- Provider SMTP nur auf `127.0.0.1` (bei `tls_mode: none`)
- Client-ID `ppcweb`, `allowed_from` z. B. `noreply@ppc.at`
- `public_addr` / `admin_addr` nur Loopback

Details: [`KONTAKTFORMULAR.md`](KONTAKTFORMULAR.md) Abschnitt 4.3.

### 3.6 ppc-contact

Kurzfassung – ausführlich: [`INSTALL_PPC_CONTACT.md`](INSTALL_PPC_CONTACT.md).

```bash
cd /opt/ppcweb/server/contact-bridge
go mod tidy
go build -o ppc-contact .

sudo useradd --system --home-dir /nonexistent --shell /usr/sbin/nologin ppc-contact 2>/dev/null || true
sudo install -d -o root -g ppc-contact -m 0750 /etc/ppc-contact
sudo install -m 0755 ppc-contact /usr/local/bin/ppc-contact

sudo cp config.example.yaml /etc/ppc-contact/config.yaml
sudo chown root:ppc-contact /etc/ppc-contact/config.yaml
sudo chmod 0640 /etc/ppc-contact/config.yaml
sudo nano /etc/ppc-contact/config.yaml
# altcha_secret, mcon.hmac_secret (= M-Con clients.ppcweb.hmac_secret),
# listen 127.0.0.1:8690, allowed_origins https://Ihre-Domain

sudo cp /opt/ppcweb/server/deploy/ppc-contact.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable --now ppc-contact
curl -sS http://127.0.0.1:8690/healthz
# Erwartung: ok
```

### 3.7 Website in Seda legen

```bash
sudo rsync -a --delete /opt/ppcweb/theweb/ /PFAD/SEDA/ROOT/
# Seda reload laut Seda-Doku (Admin-Socket / SIGHUP)
```

### 3.8 SeDat: `/api/` routen

In SeDat Path-Prefix `/api/` → `http://127.0.0.1:8690`, dann SeDat reload.

### 3.9 Gesamttest

```bash
curl -sS http://127.0.0.1:8690/healthz
curl -sS http://127.0.0.1:8690/api/altcha/challenge | head -c 200
```

Von außen:

1. `https://Ihre-Domain/` und `/kontakt.html`
2. Formular + ALTCHA absenden → Mail prüfen
3. Optional: M-Con Admin Message-Status `sent`

Checkliste Ende: Abschnitt **6**.

---

## 4. Update (Website und/oder Backend)

### 4.1 Vor dem Update

1. Auf dem Entwicklungsrechner Version bumpfen (`VERSION`, `index.html`, `CHANGELOG.md`) – siehe [`DOKUMENTATION.md`](../DOKUMENTATION.md) Abschnitt 3  
2. Commit + Push nach GitHub  
3. Optional: Tag `v0.1.1.x` setzen und pushen  

### 4.2 Auf dem Server: Repo aktualisieren

```bash
cd /opt/ppcweb
git fetch --tags origin
git checkout v0.1.1.8          # neuer Tag, oder:
# git checkout main && git pull --ff-only
```

### 4.3 Nur Website aktualisieren

```bash
# Apache/Nginx-Root
sudo rsync -a --delete /opt/ppcweb/theweb/ /var/www/ppc.at/

# oder Seda-Root
sudo rsync -a --delete /opt/ppcweb/theweb/ /PFAD/SEDA/ROOT/
# danach Seda reload
```

`--delete` entfernt Dateien, die im Repo nicht mehr existieren. Configs außerhalb von `theweb/` bleiben unberührt.

Ohne Git auf dem Server: lokal `theweb/` per `scp`/`rsync`/SFTP überschreiben.

### 4.4 ppc-contact aktualisieren

```bash
cd /opt/ppcweb/server/contact-bridge
go mod tidy
go build -o ppc-contact .
sudo install -m 0755 ppc-contact /usr/local/bin/ppc-contact
sudo systemctl restart ppc-contact
curl -sS http://127.0.0.1:8690/healthz
```

**Config nicht überschreiben:** `/etc/ppc-contact/config.yaml` bleibt.  
Nur bei neuen Optionen mit `config.example.yaml` vergleichen und manuell ergänzen.

### 4.5 M-Con / smtp-forward / SeDat

Nur nötig, wenn sich Deploy-Dateien oder Upstream geändert haben:

```bash
# Unit-Beispiel erneut kopieren nur wenn gewollt, dann:
sudo systemctl daemon-reload
sudo systemctl restart smtp-forward   # falls geändert
sudo systemctl restart mcon           # nach Config-Merge
# SeDat / Seda: reload laut jeweiliger Doku
```

M-Con Secrets und Produktions-`mcon.yaml` nie blind aus dem Repo überschreiben.

### 4.6 Schnell-Update (alles in einem Rutsch)

```bash
cd /opt/ppcweb
git fetch --tags && git checkout v0.1.1.8

sudo rsync -a --delete theweb/ /PFAD/SEDA/ROOT/
# Seda reload …

cd server/contact-bridge
go mod tidy && go build -o ppc-contact .
sudo install -m 0755 ppc-contact /usr/local/bin/ppc-contact
sudo systemctl restart ppc-contact
curl -sS http://127.0.0.1:8690/healthz
```

---

## 5. Update ohne Git auf dem Server (von Windows)

PowerShell, Projektordner anpassen:

```powershell
$Server = "IHR_USER@SERVER"
$Local  = "C:\Users\Preier.PPC\Meine Ablage\PPCNeu.2026\PPCWeb"

# Website
scp -r "$Local\theweb\*" "${Server}:/tmp/ppcweb-theweb/"
# auf dem Server danach nach Document-Root / Seda-Root rsyncen

# ppc-contact Cross-Compile
cd "$Local\server\contact-bridge"
$env:GOOS="linux"; $env:GOARCH="amd64"
go mod tidy
go build -o ppc-contact .
scp .\ppc-contact "${Server}:/tmp/ppc-contact"
```

Auf dem Server:

```bash
sudo install -m 0755 /tmp/ppc-contact /usr/local/bin/ppc-contact
sudo systemctl restart ppc-contact
```

Empfehlung langfristig: **Git auf dem Server** (Abschnitte 3–4) – Updates sind dann kurz und reproduzierbar.

---

## 6. Checkliste Neuinstallation

### Variante A

- [ ] DNS zeigt auf den Server  
- [ ] Inhalt von `theweb/` im Document Root  
- [ ] HTTPS / Zertifikat ok  
- [ ] Startseite + Dark Mode + Mobilmenü  
- [ ] `/.well-known/security.txt` erreichbar  

### Variante B (zusätzlich)

- [ ] `/opt/ppcweb` ausgecheckt (Tag)  
- [ ] Seda liefert die Site  
- [ ] SeDat: TLS + `/api/` → `127.0.0.1:8690`  
- [ ] `smtp-forward` aktiv (falls Phase-1-SMTP)  
- [ ] M-Con aktiv, Secrets gesetzt, Loopback only  
- [ ] `ppc-contact` aktiv, HMAC = M-Con-Client  
- [ ] `curl …/healthz` und `/api/altcha/challenge` → ok / JSON  
- [ ] Testmail über Kontaktformular  

---

## 7. Häufige Fehler

| Problem | Ursache | Lösung |
|---|---|---|
| Alte Seite nach Deploy | Cache / Seda nicht reloadet | Hart neu laden; Seda reload |
| ALTCHA „Prüfung fehlgeschlagen“ | `/api/` nicht geroutet | SeDat Prefix auf ppc-contact |
| `send_failed` / 502 | HMAC falsch oder M-Con down | Secrets vergleichen; `journalctl -u mcon` |
| `ppc-contact` startet nicht | `Requires=mcon.service` | M-Con starten oder Unit auf `Wants=` |
| Go-Build scheitert | altes Go / kein Netz | Go 1.23+; Proxy für `proxy.golang.org` |
| Leeres Root nach `rsync --delete` | falsches Zielverzeichnis | Zielpfad vor dem Lauf prüfen |

---

## 8. Verwandte Dateien im Repo

| Pfad | Zweck |
|---|---|
| `theweb/` | Öffentliche Website |
| `server/contact-bridge/` | Quellcode `ppc-contact` |
| `server/deploy/ppc-contact.service` | systemd |
| `server/deploy/smtp-forward.service` | Loopback-SMTP → Relay |
| `server/deploy/mcon-ppc.yaml.example` | M-Con PPC-Ausschnitt |
| `server/deploy/mcon-secrets.ppc.example.yaml` | Secrets-Vorlage |
| `server/deploy/sedat-ppc-snippet.toml` | SeDat-Routing-Hinweis |

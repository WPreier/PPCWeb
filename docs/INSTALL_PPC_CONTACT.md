# ppc-contact auf kmbhelp installieren (Schritt für Schritt)

Zielhost: **kmbhelp** `192.168.36.55`  
Quelle: GitHub-Repo **https://github.com/WPreier/PPCWeb** (öffentlich)  
Binary landet unter: **`/usr/local/bin/ppc-contact`**  
Config: **`/etc/ppc-contact/config.yaml`**

**Gesamte PPCWeb-Server-Installation / Updates (Website + Stack):**  
→ [`INSTALL_PPCWEB.md`](INSTALL_PPCWEB.md)

Ja: **`git clone` (bzw. später `git pull`) ist der empfohlene Weg.**  
Sie klonen das ganze PPCWeb-Repo; gebaut wird nur der Unterordner `server/contact-bridge/`.

---

## 0. Kurzüberblick

| Schritt | Was passiert |
|---|---|
| 1 | Per SSH auf kmbhelp einloggen |
| 2 | `git` und Go installieren (einmalig) |
| 3 | Repo nach `/opt/ppcweb` klonen |
| 4 | In `server/contact-bridge` mit Go bauen |
| 5 | Binary + Config nach `/usr/local/bin` und `/etc/ppc-contact` |
| 6 | systemd-Unit aktivieren |
| 7 | Mit `curl` testen |

Voraussetzung für den Mailversand: **M-Con** läuft bereits (oder Sie richten es parallel ein). Die Unit `ppc-contact.service` enthält `Requires=mcon.service` – ohne M-Con startet ppc-contact nicht. Wenn M-Con noch fehlt, erst M-Con installieren **oder** in der Unit vorübergehend `Requires=mcon.service` entfernen / in `Wants=` ändern.

---

## 1. Per SSH auf den Server

Von Ihrem Windows-PC (PowerShell oder Windows-Terminal):

```powershell
ssh IHR_USER@192.168.36.55
```

Ersetzen Sie `IHR_USER` durch den Account, mit dem Sie auf kmbhelp arbeiten (z. B. Ihr Admin-User).  
Wenn der Login mit Passwort geht, reicht das. Mit Key:

```powershell
ssh -i $env:USERPROFILE\.ssh\id_ed25519 IHR_USER@192.168.36.55
```

Nach dem Login sind Sie **auf dem Server** – alle folgenden Befehle dort ausführen (Linux-Bash).

Prüfen:

```bash
hostname
whoami
pwd
```

---

## 2. Werkzeuge installieren (einmalig)

### Debian / Ubuntu

```bash
sudo apt update
sudo apt install -y git curl build-essential

# Go: Distro-Paket oft zu alt. Offizielle Binary empfohlen:
curl -fsSL https://go.dev/dl/go1.23.6.linux-amd64.tar.gz -o /tmp/go.tgz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf /tmp/go.tgz
echo 'export PATH=/usr/local/go/bin:$PATH' | sudo tee /etc/profile.d/go.sh
source /etc/profile.d/go.sh
go version
# Erwartung: go version go1.23.x linux/amd64
```

### AlmaLinux / RHEL

```bash
sudo dnf install -y git curl gcc
# Go analog von https://go.dev/dl/ nach /usr/local/go
```

---

## 3. Repo klonen

Empfohlener Ort: **`/opt/ppcweb`** (nur Quellcode zum Bauen, nicht Document-Root von Seda).

```bash
sudo mkdir -p /opt
sudo chown "$USER":"$USER" /opt   # nur falls /opt noch Ihnen gehören soll – sonst:
# Alternative ohne chown von ganz /opt:
sudo mkdir -p /opt/ppcweb
sudo chown "$USER":"$USER" /opt/ppcweb

cd /opt
git clone https://github.com/WPreier/PPCWeb.git ppcweb
cd /opt/ppcweb

# feste Version (empfohlen für Produktion):
git checkout v0.1.1.8
# oder immer neueste main:
# git checkout main && git pull
```


Wenn `git clone` nach Authentifizierung fragt: Das Repo ist **public** – ohne Login sollte es klappen. Bei Fehlern Netzwerk/Proxy/DNS prüfen.

Struktur prüfen:

```bash
ls -la /opt/ppcweb/server/contact-bridge
# Erwartung: main.go, go.mod, config.example.yaml, build.sh
```

---

## 4. ppc-contact bauen

```bash
cd /opt/ppcweb/server/contact-bridge
go mod tidy
go build -o ppc-contact .

ls -la ppc-contact
./ppc-contact 2>&1 | head -5
# Erwartung: Meldung „config required…“ oder ähnliches – Binary startet also.
```

`go mod tidy` lädt Abhängigkeiten (ALTCHA-Lib, yaml) aus dem Internet. Ohne ausgehenden HTTPS-Zugang schlägt der Build fehl.

Optional mit dem Hilfsskript:

```bash
chmod +x build.sh
./build.sh
```

---

## 5. Systembenutzer, Verzeichnisse, Binary installieren

```bash
# Systemuser (ohne Login-Shell)
sudo useradd --system --home-dir /nonexistent --shell /usr/sbin/nologin ppc-contact 2>/dev/null || true

# Config-Verzeichnis
sudo install -d -o root -g ppc-contact -m 0750 /etc/ppc-contact

# Binary an den Systempfad
sudo install -m 0755 /opt/ppcweb/server/contact-bridge/ppc-contact /usr/local/bin/ppc-contact

# Welche Binary läuft später?
which ppc-contact
ppc-contact 2>&1 | head -3
```

---

## 6. Konfiguration anlegen

```bash
sudo cp /opt/ppcweb/server/contact-bridge/config.example.yaml /etc/ppc-contact/config.yaml
sudo chown root:ppc-contact /etc/ppc-contact/config.yaml
sudo chmod 0640 /etc/ppc-contact/config.yaml
```

Secrets erzeugen:

```bash
openssl rand -hex 24   # → ALTCHA-Secret
openssl rand -hex 24   # → muss GLEICH sein wie clients.ppcweb.hmac_secret in /etc/mcon/secrets.yaml
```

Datei bearbeiten:

```bash
sudo nano /etc/ppc-contact/config.yaml
```

Mindestens setzen:

```yaml
listen: "127.0.0.1:8690"

altcha_secret: "HIER-DAS-ERSTE-OPENSSL-ERGEBNIS"

mcon:
  base_url: "http://127.0.0.1:8687"
  key_id: "ppcweb"
  hmac_secret: "HIER-DASSELBE-WIE-IN-MCON-secrets.yaml"
  from: "noreply@ppc.at"
  to:
    - "preier@ppc.at"
  provider: ""
  timeout_seconds: 15

allowed_origins:
  - "https://ppc.at"
  - "https://www.ppc.at"
```

**Wichtig:** `mcon.hmac_secret` und in M-Con `clients.ppcweb.hmac_secret` müssen **bytegleich** sein, sonst lehnt M-Con die Signatur ab.

---

## 7. systemd-Dienst einrichten

```bash
sudo cp /opt/ppcweb/server/deploy/ppc-contact.service /etc/systemd/system/ppc-contact.service
sudo systemctl daemon-reload
sudo systemctl enable --now ppc-contact.service
sudo systemctl status ppc-contact.service --no-pager
```

Wenn der Start wegen fehlendem M-Con scheitert:

```bash
systemctl status mcon --no-pager
# M-Con zuerst starten, oder in der Unit Requires=mcon.service → Wants=mcon.service ändern:
sudo systemctl edit --full ppc-contact.service
# danach:
sudo systemctl daemon-reload
sudo systemctl restart ppc-contact
```

Logs:

```bash
journalctl -u ppc-contact -e --no-pager
```

---

## 8. Funktionstest (nur Loopback)

```bash
curl -sS http://127.0.0.1:8690/healthz
# Erwartung: ok

curl -sS http://127.0.0.1:8690/api/altcha/challenge | head -c 400
# Erwartung: JSON mit algorithm, challenge, salt, signature, maxNumber
# NICHT: HTML oder Connection refused
```

Wenn das klappt, fehlt für die Website nur noch:

1. SeDat: `/api/` → `http://127.0.0.1:8690`  
2. `theweb/` in Seda aktualisieren  

Details: [`KONTAKTFORMULAR.md`](KONTAKTFORMULAR.md) Abschnitte 4.5–4.6 und 8 (Diagnose ALTCHA).

---

## 9. Später aktualisieren (neue Version)

```bash
cd /opt/ppcweb
git fetch --tags
git checkout v0.1.1.8          # oder neuerer Tag / main
cd server/contact-bridge
go mod tidy
go build -o ppc-contact .
sudo install -m 0755 ppc-contact /usr/local/bin/ppc-contact
sudo systemctl restart ppc-contact
curl -sS http://127.0.0.1:8690/healthz
```


Config unter `/etc/ppc-contact/config.yaml` bleibt erhalten (nicht überschreiben, außer neue Optionen nötig).

Website-Dateien separat:

```bash
# Inhalt von /opt/ppcweb/theweb/ in Ihr Seda-Root kopieren, dann Seda reload
# (Pfad zum Seda-Root Ihrer Installation einsetzen)
sudo rsync -a --delete /opt/ppcweb/theweb/ /PFAD/ZUM/SEDA/ROOT/
# danach Seda reload laut Seda-Doku
```

---

## 10. Alternative ohne Clone auf dem Server: von Windows kopieren

Wenn auf kmbhelp kein Git/Go möglich ist:

### A) Auf Windows bauen (Cross-Compile) und per SCP legen

PowerShell (Go muss auf dem PC installiert sein):

```powershell
cd "C:\Users\Preier.PPC\Meine Ablage\PPCNeu.2026\PPCWeb\server\contact-bridge"
$env:GOOS="linux"; $env:GOARCH="amd64"; go mod tidy; go build -o ppc-contact .
scp .\ppc-contact IHR_USER@192.168.36.55:/tmp/ppc-contact
scp .\config.example.yaml IHR_USER@192.168.36.55:/tmp/config.example.yaml
scp ..\deploy\ppc-contact.service IHR_USER@192.168.36.55:/tmp/ppc-contact.service
```

Auf dem Server:

```bash
sudo install -m 0755 /tmp/ppc-contact /usr/local/bin/ppc-contact
# dann Config + systemd wie in Abschnitt 5–7
```

### B) Nur Quellcode per SCP, auf dem Server bauen

```powershell
# Von Windows: Ordner contact-bridge hochladen
scp -r "C:\Users\Preier.PPC\Meine Ablage\PPCNeu.2026\PPCWeb\server\contact-bridge" IHR_USER@192.168.36.55:/tmp/
```

```bash
cd /tmp/contact-bridge
go mod tidy && go build -o ppc-contact .
sudo install -m 0755 ppc-contact /usr/local/bin/ppc-contact
```

**Empfehlung:** Variante **git clone auf dem Server** (Abschnitte 3–8) – Updates sind dann ein `git pull` + Rebuild.

---

## 11. Häufige Fehler

| Problem | Ursache | Lösung |
|---|---|---|
| `git: command not found` | git fehlt | `apt install git` |
| `go: command not found` | PATH / Go nicht installiert | Abschnitt 2, `source /etc/profile.d/go.sh` |
| `go mod tidy` Timeout | kein Internet | Proxy/DNS/Firewall für `proxy.golang.org` |
| `Connection refused :8690` | Dienst nicht aktiv | `systemctl status ppc-contact` |
| Unit fails, M-Con missing | `Requires=mcon.service` | M-Con starten oder Unit anpassen |
| Challenge JSON ok, Website rot | SeDat `/api` fehlt | SeDat-Route setzen |
| Absenden `send_failed` / 502 | HMAC falsch oder M-Con down | Secrets vergleichen, `journalctl -u mcon` |

---

## 12. Checkliste „ppc-contact ist drauf“

- [ ] `ssh` auf kmbhelp klappt  
- [ ] `/opt/ppcweb` existiert, Tag ausgecheckt  
- [ ] `/usr/local/bin/ppc-contact` vorhanden  
- [ ] `/etc/ppc-contact/config.yaml` mit echten Secrets  
- [ ] `systemctl is-active ppc-contact` → `active`  
- [ ] `curl http://127.0.0.1:8690/healthz` → `ok`  
- [ ] `curl http://127.0.0.1:8690/api/altcha/challenge` → JSON  

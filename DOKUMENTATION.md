# Dokumentation – PPCWeb (ppc.at)

Stand: **2026-11-25** · Aktuelle Version: siehe Datei [`VERSION`](VERSION) und Footer auf der Startseite (`theweb/index.html`).

**Kontaktformular / M-Con / Captcha:** ausführlich in [`docs/KONTAKTFORMULAR.md`](docs/KONTAKTFORMULAR.md).

Diese Dokumentation beschreibt Aufbau, Versionierung, Pflege und Deployment der Website von **Penguin Power Communication**.

---

## 0. Lokal ↔ GitHub (Verbindung)

| | |
|---|---|
| **Lokaler Pfad** | `…\Meine Ablage\PPCNeu.2026\PPCWeb` (Google Drive / Meine Ablage) |
| **Remote** | `origin` → https://github.com/WPreier/PPCWeb.git |
| **Branch** | `main` (tracked: `origin/main`) |
| **Web-Inhalt** | Ordner `theweb/` |
| **Versionierung** | aktiv (Schema `0.1.1.<Patch>`, Dateien unten) |

Prüfung (PowerShell im Projektordner):

```bash
git remote -v
git status -sb
Get-Content VERSION
```

Wenn `origin` auf `WPreier/PPCWeb` zeigt und `main...origin/main` ohne „ahead/behind“ ist, sind lokal und GitHub synchron.

**Hinweis:** Ein verschachtelter Ordner `PPCWeb/PPCWeb/` (Duplikat durch Sync) gehört **nicht** zum Repo und steht in `.gitignore`. Arbeiten Sie immer im äußeren Projektordner mit `theweb/`, `VERSION` und `.git`.

---

## 1. Was ist dieses Projekt?

Statische Website (HTML / CSS / JS) für **https://ppc.at** – persönlicher IT-Dienstleister (Administration, Security) in Ostösterreich.

- Kein Build-Schritt, kein Framework
- Deutsch als Hauptsprache
- Hell-/Dunkelmodus
- Kontakt primär per **Telefon** und **E-Mail** (kein Formular)

Inhaber / Ansprechpartner: **Ing. Wolfgang Preier** · `preier@ppc.at` · `+43 664 8585826`

---

## 2. Verzeichnisstruktur

```
PPCWeb/                          ← Projekt / Git-Repo
├── VERSION                      Aktuelle Versionsnummer (eine Zeile)
├── CHANGELOG.md                 Alle Änderungen mit Datum & Version
├── DOKUMENTATION.md             Diese Datei
├── README.md                    Kurzüberblick
├── .gitignore
└── theweb/                      ← Document Root: nur das auf den Webserver
    ├── index.html               Startseite (enthält Versionsnummer)
    ├── security.html / administration.html / …
    ├── kontakt.html / impressum.html / 404.html
    ├── robots.txt / sitemap.xml / humans.txt
    ├── site.webmanifest / .htaccess
    ├── .well-known/security.txt
    └── assets/
        ├── css/style.css
        ├── js/main.js
        └── images/ (logo, penguin, favicon.svg, …)
```

| Datei / Ordner | Aufgabe |
|---|---|
| `theweb/` | Alles, was öffentlich ausgeliefert wird |
| `theweb/index.html` | Einstieg; Hero; Version im `<meta>` und Footer |
| `theweb/assets/css/style.css` | Gesamtes Styling (CSS-Variablen für Theme) |
| `theweb/assets/js/main.js` | Mobiles Menü, Light/Dark, `localStorage` |
| `theweb/.well-known/security.txt` | Offizieller Sicherheitskontakt (siehe unten) |
| `VERSION` | Maschinenlesbare Versionsnummer (nur im Repo) |
| `CHANGELOG.md` | Historie jeder Patch-Version (nur im Repo) |

### Was ist `.well-known` / `security.txt`?

Das ist **kein** geheimes Systemverzeichnis, sondern ein **Internet-Standard**:

- URL: `https://ppc.at/.well-known/security.txt`
- Inhalt: Wohin Sicherheitsforscher eine Lücke melden sollen (E-Mail, Telefon, Sprachen, Policy)
- Gehört zu den **Web-Dateien** und muss im Document Root liegen (bei uns also in `theweb/`)
- Für Besucher der normalen Website unsichtbar; Tools und Suchende erwarten genau diesen Pfad

### Weitere Profi-Dateien in `theweb/`

| Datei | Zweck |
|---|---|
| `robots.txt` | Suchmaschinen: was crawlen darf + Link zur Sitemap |
| `sitemap.xml` | Liste aller öffentlichen Seiten für Google & Co. |
| `humans.txt` | Menschliche Projektinfos (Autor, Kontakt) |
| `site.webmanifest` | Name, Farben, Icon (Browser / „Zum Home-Bildschirm“) |
| `assets/images/favicon.svg` | Browser-Tab-Icon |
| `404.html` | Freundliche Fehlerseite |
| `.htaccess` | Apache: Security-Header, Kompression, 404 (nur Apache) |

---

## 3. Versionierung

### Schema

```
0.1.1.<Patchnummer>
 │ │ │ └─ Patch: jede konkrete Änderung / Veröffentlichung (1, 2, 3, …)
 │ │ └─── festgehaltene Basis-Linie (Minor-Fix-Ebene)
 │ └───── Produktlinie 0.1
 └─────── Major 0 (Aufbau / Redesign-Phase)
```

Aktuell: **0.1.1.6**

### Wo steht die Version?

| Ort | Inhalt |
|---|---|
| `VERSION` | Nur die Nummer, z. B. `0.1.1.4` |
| `theweb/index.html` | `<meta name="version" content="…">` und Footer |
| `CHANGELOG.md` | Abschnitt `## [0.1.1.x] – JJJJ-MM-TT` |
| Git-Tag (optional) | `v0.1.1.4` |

### Neue Version veröffentlichen (Checkliste)

1. Änderungen am Code/Content in `theweb/` vornehmen
2. In `VERSION` die Patchnummer erhöhen (z. B. `0.1.1.3` → `0.1.1.4`)
3. In `theweb/index.html` beide Stellen anpassen (`meta` + Footer + HTML-Kommentar)
4. In `CHANGELOG.md` neuen Abschnitt mit **Datum** und **Nummer** oben einfügen
5. Optional: diese Doku aktualisieren; bei neuen Seiten `sitemap.xml` ergänzen
6. Git committen (und pushen)
7. Optional: `git tag v0.1.1.4` und Tag pushen
8. **Inhalt von `theweb/`** auf den Webserver deployen (siehe Abschnitt 5)

---

## 4. Design & Technik (Kurz)

### Farben (CSS-Variablen in `theweb/assets/css/style.css`)

- Hintergrund hell: `--bg` / warm: `--bg-warm`
- Text: `--text`, gedämpft: `--text-muted`
- Marken-Ink: `--ink` (Nachtblau)
- Akzent / CTA: `--accent` (Pinguin-Gelb)
- Dark Mode: `[data-theme="dark"]` setzt dieselben Variablen neu

### Interaktion

- Theme: Button im Header; Speicherung in `localStorage` (`theme`)
- Beim ersten Besuch: Systemeinstellung `prefers-color-scheme` berücksichtigen
- Mobiles Menü: Hamburger, schließt bei Link-Klick und Escape

### Schriften

Nunito über [Bunny Fonts](https://fonts.bunny.net) (Datenschutzhinweis im Impressum).

---

## 5. Deployment auf den eigenen Webserver

Die Seite ist **fertig zum Hochladen** – keine Installation, kein `npm build`.

1. Lokal prüfen:
   ```bash
   cd theweb
   python -m http.server 8000
   ```
2. Per FTP/SFTP/Panel **nur den Inhalt von `theweb/`** hochladen:
   - alle `*.html`
   - Ordner `assets/`
   - Ordner `.well-known/`
3. Document Root = Ordner, in dem `index.html` liegt (also der hochgeladene `theweb`-Inhalt)
4. Prüfen: Startseite, Telefon-Links, Dark Mode, Mobilmenü  
   und optional: `https://ppc.at/.well-known/security.txt`

**Tipp:** Nach dem Upload hart neu laden (Strg+F5).

Markdown-Dateien (`README`, `CHANGELOG`, …) gehören **nicht** zwingend auf den öffentlichen Server.

---

## 6. Inhalte pflegen – was wo ändern?

| Ziel | Datei |
|---|---|
| Startseite / Hero / CTA | `theweb/index.html` |
| Security-Leistungen | `theweb/security.html` |
| Administration | `theweb/administration.html` |
| Kundenliste | `theweb/referenzen.html` |
| Partner / Malama | `theweb/partner.html` |
| Telefon, E-Mail, Standort | `theweb/kontakt.html` (+ Footer auf allen Seiten) |
| Rechtliches | `theweb/impressum.html` |
| Farben, Abstände, Layout | `theweb/assets/css/style.css` |
| Menü / Theme-Verhalten | `theweb/assets/js/main.js` |
| Logo / Pinguin / Partnerlogo | `theweb/assets/images/` |
| Security-Kontakt | `theweb/.well-known/security.txt` |

**Telefonnummer ändern:** Suche in `theweb/` nach `+436648585826` und `+43 664 8585826`.

**Header/Footer:** In jeder HTML-Datei kopiert – Nav-Änderungen auf allen Seiten in `theweb/`.

---

## 7. Rechtliches & Sicherheit

- Impressum: Adresse Pilgersdorf, UID `ATU44653709`, Kontaktdaten
- Datenschutz: kein Formular, kein Marketing-Tracking; Theme nur lokal
- `.well-known/security.txt`: Kontakt für Sicherheitsmeldungen (Ablaufdatum beachten)

---

## 8. Git / Sicherung

- Änderungen versionieren: Commit mit aussagekräftiger Message
- Remote (GitHub): `https://github.com/WPreier/PPCWeb.git`
- Empfohlen: nach jeder Patch-Version committen **und** deployen
- ZIP-Sicherungen außerhalb des Repo-Ordners, z. B.  
  `..\PPCWeb_Backup_0.1.1.x_JJJJ-MM-TT.zip`

Siehe auch `README.md` und `CHANGELOG.md`.

---

## 9. Offene optionale Erweiterungen

- Porträtfoto „Ihr Ansprechpartner“ (noch kein Bild im Projekt)
- Logo als SVG (hell/dunkel) statt JPG
- Bilder weiter komprimieren / WebP
- Gemeinsame Header-/Footer-Includes

---

## 10. Schnellreferenz Versionen

| Version | Datum | Kurzbeschreibung |
|---|---|---|
| 0.1.1.6 | 2026-11-25 | ALTCHA Challenge-Diagnose / API-Routing-Hinweis |
| 0.1.1.5 | 2026-11-25 | Kontaktformular + ALTCHA + M-Con-Bridge |
| 0.1.1.4 | 2026-11-22 | GitHub-Verbindung dokumentiert, Sync-Duplikat ignoriert |
| 0.1.1.3 | 2026-11-20 | Profi-Erweiterungen (SEO, security.txt, 404, …) |
| 0.1.1.2 | 2026-11-20 | Web-Dateien nach `theweb/` |
| 0.1.1.1 | 2026-11-20 | Redesign warm, Phone-first, Doku & Versionierung |
| 1.0.0 | 2026-11-20 | Erster Entwurf (altes Schema, abgelöst) |

Vollständige Details: [`CHANGELOG.md`](CHANGELOG.md)

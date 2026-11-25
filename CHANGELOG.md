# Changelog – PPCWeb

Alle Änderungen an diesem Projekt werden in dieser Datei dokumentiert.

**Versionsschema:** `0.1.1.<Patchnummer>`  
Beispiel: `0.1.1.1` = Basis 0.1.1, Patch 1.

Bei jeder Veröffentlichung:
1. Patchnummer erhöhen (Datei `VERSION`)
2. Eintrag hier mit **Datum** und **Versionsnummer** ergänzen
3. Versionsnummer in `theweb/index.html` aktualisieren (`meta name="version"` und Footer)
4. Git-Commit (und Tag optional: `v0.1.1.5`)

---

## [0.1.1.5] – 2026-11-25

### Added
- Kontaktformular mit ALTCHA (selbst gehostet), Honeypot, Timing- und Rate-Limits
- Server-Bridge `server/contact-bridge` (ppc-contact) für HMAC-sichere Anbindung an M-Con
- Deploy-Vorlagen: systemd Units, M-Con-Ausschnitt Phase-1-SMTP, SeDat-Hinweis
- Dokumentation `docs/KONTAKTFORMULAR.md` (Setup, Tests, Sicherheit, SMTP→MS365 OAuth)

### Changed
- `kontakt.html` wieder mit Formular; Datenschutztext im Impressum angepasst

### Notes
- Phase-1-Versand: SMTP über lokalen Forwarder `127.0.0.1:2525` → `192.168.36.251:25` (M-Con erlaubt Klartext nur auf Loopback)

---

## [0.1.1.4] – 2026-11-22

### Changed
- Dokumentation: Abschnitt „Lokal ↔ GitHub“ (Pfad Meine Ablage, Remote, Sync-Prüfung)
- `.gitignore`: `.cursor/` und verschachteltes `/PPCWeb/` ausgeschlossen

### Notes
- Versionierung war bereits aktiv (seit 0.1.1.1); dieser Patch bestätigt Verbindung und Absicherung gegen Sync-Duplikate

---

## [0.1.1.3] – 2026-11-20

### Added
- `robots.txt`, `sitemap.xml`, `humans.txt`, `site.webmanifest`, Favicon (SVG)
- Open-Graph- und Twitter-Meta-Tags, Canonical-URLs, `theme-color` auf allen Seiten
- JSON-LD (`ProfessionalService`) auf der Startseite
- Skip-Link für Barrierefreiheit
- Professionelle `404.html`
- Apache-`.htaccess` (Security-Header, Kompression, 404)
- `security.txt` um `Policy` ergänzt

### Notes
- Erwartete Profi-Erweiterungen für SEO, Sicherheitshinweise und Social Sharing

---

## [0.1.1.2] – 2026-11-20

### Changed
- Alle auslieferbaren Web-Dateien nach `theweb/` verschoben (HTML, `assets/`, `.well-known/`)
- Dokumentation und README an die neue Struktur angepasst

### Notes
- Projektroot enthält Doku/Versionierung; Document Root auf dem Server = Inhalt von `theweb/`

---

## [0.1.1.1] – 2026-11-20

### Added
- Warmes Marken-Farbsystem (Nachtblau + Pinguin-Gelb), Dark Mode mit Charakter
- Phone-first Kontakt: Telefon im Header, Hero, CTAs und Footer (`tel:`-Links)
- Neue Startseite nach Leitidee (Hero, Leistungen, Vergleich, 3 Schritte, CTA)
- SVG-Icons statt Emojis auf allen Seiten
- Partnerseite mit Malama IT (Logo aus dem Header dorthin verschoben)
- Vollständiges Impressum (Adresse, UID, WKO, Gewerbebehörde) und Datenschutz
- Projektdokumentation (`DOKUMENTATION.md`), Versionsdatei (`VERSION`), Versionshistorie

### Changed
- Header entschlackt (kein Malama-Logo mehr in der Navigation)
- Schrift: Nunito über Bunny Fonts (statt Inter / Google Fonts)
- Kontaktseite ohne Formular – nur Telefon, E-Mail und Standort
- Referenzen-Listenmarker als Häkchen-Icons

### Removed
- Hosting als Leistungsangebot (Texte/Meta)
- Kontaktformular
- Hinweis „made with Grok“ im Footer

### Notes
- Erste versionierte Veröffentlichung unter dem Schema `0.1.1.<Patch>`
- Vorläufer-Stand war intern als Konzept `1.0.0` (2026-11-20) dokumentiert; siehe unten

---

## [1.0.0] – 2026-11-20 (Vorgänger / unveröffentlicht unter neuem Schema)

### Added
- Initial modern website for ppc.at
- Light/dark mode toggle
- Security, Administration, Referenzen, Partner, Kontakt, Impressum
- security.txt
- Responsive design with penguin branding

### Notes
- Based on previous Pflichtenheft and content documents
- Abgelöst durch **0.1.1.1** (neues Versionsschema und Redesign)

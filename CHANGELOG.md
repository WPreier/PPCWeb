# Changelog – PPCWeb

Alle Änderungen an diesem Projekt werden in dieser Datei dokumentiert.

**Versionsschema:** `0.1.1.<Patchnummer>`  
Beispiel: `0.1.1.1` = Basis 0.1.1, Patch 1.

Bei jeder Veröffentlichung:
1. Patchnummer erhöhen (Datei `VERSION`)
2. Eintrag hier mit **Datum** und **Versionsnummer** ergänzen
3. Versionsnummer in `index.html` aktualisieren (`meta name="version"` und Footer)
4. Git-Commit (und Tag optional: `v0.1.1.1`)

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

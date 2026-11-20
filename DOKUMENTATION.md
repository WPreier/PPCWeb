# Dokumentation – PPCWeb (ppc.at)

Stand: **2026-11-20** · Aktuelle Version: siehe Datei [`VERSION`](VERSION) und Footer auf der Startseite.

Diese Dokumentation beschreibt Aufbau, Versionierung, Pflege und Deployment der Website von **Penguin Power Communication**.

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
PPCWeb/
├── index.html              Startseite (enthält Versionsnummer)
├── security.html           IT-Security
├── administration.html     IT-Administration
├── referenzen.html         Kundenreferenzen
├── partner.html            Partner (u. a. Malama IT)
├── kontakt.html            Kontakt (Telefon / E-Mail / Adresse)
├── impressum.html          Impressum & Datenschutz
├── VERSION                 Aktuelle Versionsnummer (eine Zeile)
├── CHANGELOG.md            Alle Änderungen mit Datum & Version
├── DOKUMENTATION.md        Diese Datei
├── README.md               Kurzüberblick (auch auf Englisch lesbar)
├── .gitignore
├── .well-known/
│   └── security.txt        Security-Kontakt für Forscher
└── assets/
    ├── css/style.css       Design, Farben, Layout, Dark Mode
    ├── js/main.js          Menü + Theme-Umschalter
    └── images/
        ├── logo.jpg        PPC-Logo (Wortmarke + Pinguin)
        ├── penguin.jpg     Pinguin-Motiv (Hero)
        └── Malama_IT_Logo.jpg
```

| Datei / Ordner | Aufgabe |
|---|---|
| `index.html` | Einstieg; Hero; Version im `<meta>` und Footer |
| `assets/css/style.css` | Gesamtes Styling (CSS-Variablen für Theme) |
| `assets/js/main.js` | Mobiles Menü, Light/Dark, `localStorage` |
| `VERSION` | Maschinenlesbare Versionsnummer |
| `CHANGELOG.md` | Historie jeder Patch-Version |

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

Aktuell: **0.1.1.1** (erste Patch-Version der Basis 0.1.1).

### Wo steht die Version?

| Ort | Inhalt |
|---|---|
| `VERSION` | Nur die Nummer, z. B. `0.1.1.1` |
| `index.html` | `<meta name="version" content="…">` und Footer `v0.1.1.1` |
| `CHANGELOG.md` | Abschnitt `## [0.1.1.x] – JJJJ-MM-TT` |
| Git-Tag (optional) | `v0.1.1.1` |

### Neue Version veröffentlichen (Checkliste)

1. Änderungen am Code/Content vornehmen
2. In `VERSION` die Patchnummer erhöhen (z. B. `0.1.1.1` → `0.1.1.2`)
3. In `index.html` beide Stellen anpassen (`meta` + Footer + HTML-Kommentar)
4. In `CHANGELOG.md` neuen Abschnitt mit **Datum** und **Nummer** oben einfügen (`Added` / `Changed` / `Removed` / `Fixed`)
5. Optional: diese Doku aktualisieren, wenn Struktur oder Abläufe sich ändern
6. Git committen (und pushen)
7. Optional: `git tag v0.1.1.2` und Tag pushen
8. Dateien auf den Webserver deployen (siehe Abschnitt 5)

---

## 4. Design & Technik (Kurz)

### Farben (CSS-Variablen in `style.css`)

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

1. Lokal prüfen: `index.html` im Browser öffnen oder  
   `python -m http.server 8000` im Projektordner
2. Per FTP/SFTP/Panel alle Dateien und Ordner hochladen:
   - alle `*.html`
   - Ordner `assets/`
   - Ordner `.well-known/`
   - optional die Markdown-Dateien (müssen nicht öffentlich sein; schaden aber nicht)
3. Document Root so setzen, dass `index.html` die Startseite ist
4. `https://ihre-domain/` aufrufen und Header, Telefon-Links, Dark Mode, Mobilmenü testen

**Tipp:** Nach dem Upload den Browser-Cache leeren oder hart neu laden (Strg+F5), damit CSS/JS sicher aktualisiert werden.

---

## 6. Inhalte pflegen – was wo ändern?

| Ziel | Datei |
|---|---|
| Startseite / Hero / CTA | `index.html` |
| Security-Leistungen | `security.html` |
| Administration | `administration.html` |
| Kundenliste | `referenzen.html` |
| Partner / Malama | `partner.html` |
| Telefon, E-Mail, Standort | `kontakt.html` (+ Footer auf allen Seiten) |
| Rechtliches | `impressum.html` |
| Farben, Abstände, Layout | `assets/css/style.css` |
| Menü / Theme-Verhalten | `assets/js/main.js` |
| Logo / Pinguin / Partnerlogo | `assets/images/` |

**Telefonnummer ändern:** Suche projektweit nach `+436648585826` und `+43 664 8585826` und ersetze konsistent (Header, Hero, Kontakt, Footer, Impressum, `security.txt`).

**Header/Footer:** Aktuell in jeder HTML-Datei kopiert. Bei Nav-Änderungen alle Seiten anpassen (oder später auf Includes/SSG umstellen).

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
- Zusätzliche ZIP-Sicherung (außerhalb des Repo-Ordners):  
  `..\PPCWeb_Backup_0.1.1.1_2026-11-20.zip` (bei neuen Versionen analog benennen)

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
| 0.1.1.1 | 2026-11-20 | Redesign warm, Phone-first, Doku & Versionierung |
| 1.0.0 | 2026-11-20 | Erster Entwurf (altes Schema, abgelöst) |

Vollständige Details: [`CHANGELOG.md`](CHANGELOG.md)

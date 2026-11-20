# PPCWeb – Penguin Power Communication

Statische Website für **Penguin Power Communication** ([ppc.at](https://ppc.at)).

Persönlicher IT-Dienstleister in Österreich: Administration und Security. Warmes Pinguin-Branding, Hell-/Dunkelmodus, Kontakt per Telefon und E-Mail.

**Aktuelle Version:** siehe Datei [`VERSION`](VERSION) (Schema `0.1.1.<Patchnummer>`)

## Dokumentation

| Datei | Inhalt |
|---|---|
| [`DOKUMENTATION.md`](DOKUMENTATION.md) | Aufbau, Pflege, Deploy, Checklisten |
| [`CHANGELOG.md`](CHANGELOG.md) | Alle Änderungen mit Datum und Versionsnummer |
| [`VERSION`](VERSION) | Aktuelle Versionsnummer (eine Zeile) |

## Features

- Responsive Design
- Light / Dark Theme
- Seiten: Start, Security, Administration, Referenzen, Partner, Kontakt, Impressum
- `security.txt` unter `.well-known/`
- Phone-first CTAs (kein Kontaktformular)

## Tech

Statisches HTML / CSS / JS – kein Build nötig.

- CSS-Variablen für Theming (`assets/css/style.css`)
- JS: Theme + mobiles Menü (`assets/js/main.js`)
- Schrift: Nunito (Bunny Fonts)

## Struktur

```
.
├── index.html                 ← Versionsnummer hier (meta + Footer)
├── administration.html
├── security.html
├── referenzen.html
├── partner.html
├── kontakt.html
├── impressum.html
├── VERSION
├── CHANGELOG.md
├── DOKUMENTATION.md
├── assets/
│   ├── css/style.css
│   ├── js/main.js
│   └── images/
└── .well-known/security.txt
```

## Lokal testen

```bash
python -m http.server 8000
```

Dann `http://localhost:8000` öffnen – oder `index.html` direkt im Browser.

## Neue Patch-Version

1. Code ändern  
2. `VERSION` erhöhen (z. B. `0.1.1.2`)  
3. `index.html` Version anpassen  
4. Eintrag in `CHANGELOG.md` mit Datum  
5. Commit / Deploy  

Details: [`DOKUMENTATION.md`](DOKUMENTATION.md)

## License

Proprietary – Penguin Power Communication / Ing. Wolfgang Preier

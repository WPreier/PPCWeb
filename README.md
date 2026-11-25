# PPCWeb – Penguin Power Communication

Statische Website für **Penguin Power Communication** ([ppc.at](https://ppc.at)).

Persönlicher IT-Dienstleister in Österreich: Administration und Security. Warmes Pinguin-Branding, Hell-/Dunkelmodus, Kontakt per Telefon und E-Mail.

**Aktuelle Version:** siehe Datei [`VERSION`](VERSION) (Schema `0.1.1.<Patchnummer>`)

**GitHub:** https://github.com/WPreier/PPCWeb · lokal unter `Meine Ablage\PPCNeu.2026\PPCWeb` (Remote `origin`)

**Kontaktformular:** [`docs/KONTAKTFORMULAR.md`](docs/KONTAKTFORMULAR.md) (Seda + M-Con + ALTCHA auf kmbhelp)

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

- CSS-Variablen für Theming (`theweb/assets/css/style.css`)
- JS: Theme + mobiles Menü (`theweb/assets/js/main.js`)
- Schrift: Nunito (Bunny Fonts)

## Struktur

```
.
├── VERSION / CHANGELOG.md / DOKUMENTATION.md / README.md
└── theweb/                      ← auf den Webserver hochladen
    ├── index.html               ← Versionsnummer (meta + Footer)
    ├── *.html / 404.html
    ├── robots.txt / sitemap.xml / humans.txt
    ├── site.webmanifest
    ├── .htaccess                ← Apache
    ├── assets/…
    └── .well-known/security.txt
```

## Lokal testen

```bash
cd theweb
python -m http.server 8000
```

Dann `http://localhost:8000` öffnen.

## Neue Patch-Version

1. Code in `theweb/` ändern  
2. `VERSION` erhöhen (z. B. `0.1.1.3`)  
3. `theweb/index.html` Version anpassen  
4. Eintrag in `CHANGELOG.md` mit Datum  
5. Commit / Deploy  

Details: [`DOKUMENTATION.md`](DOKUMENTATION.md)

## License

Proprietary – Penguin Power Communication / Ing. Wolfgang Preier

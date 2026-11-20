# PPCWeb – Penguin Power Communication

Modern, responsive website for **Penguin Power Communication** (ppc.at).

Personal IT service provider in Austria, focusing on Administration, Security, Hosting and related services. Clean design with light/dark mode, based on the company penguin branding.

## Features

- Responsive design (mobile-friendly)
- Light / Dark theme switch (light is default)
- German primary language
- Pages: Start, Security, Administration, Referenzen, Partner, Kontakt, Impressum
- Security.txt included (`.well-known/security.txt`)
- Contact form and clear call-to-actions

## Tech

Static HTML / CSS / JS. No build step required.

- CSS: Custom properties for theming
- JS: Theme toggle + mobile menu
- Fonts: Inter via Google Fonts

## Structure

```
.
├── index.html
├── administration.html
├── security.html
├── referenzen.html
├── partner.html
├── kontakt.html
├── impressum.html
├── assets/
│   ├── css/style.css
│   ├── js/main.js
│   └── images/
│       ├── logo.jpg
│       ├── penguin.jpg
│       └── Malama_IT_Logo.jpg
└── .well-known/security.txt
```

## Local preview

Open `index.html` in a browser or serve with any static server:

```bash
python3 -m http.server 8000
```

## Versioning

Semantic versioning via Git tags. Current initial version: v1.0.0

## License

Proprietary – Penguin Power Communication / Ing. Wolfgang Preier

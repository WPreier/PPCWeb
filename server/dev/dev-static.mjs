#!/usr/bin/env node
/**
 * Lokaler Dev-Server für PPCWeb:
 * - static: ../theweb auf :8001
 * - proxy:  /api/* → http://127.0.0.1:8690 (ppc-contact)
 */
import http from 'node:http';
import fs from 'node:fs';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const ROOT = path.resolve(__dirname, '../../theweb');
const PORT = Number(process.env.PPC_DEV_PORT || 8001);
const API = process.env.PPC_CONTACT_URL || 'http://127.0.0.1:8690';

const TYPES = {
  '.html': 'text/html; charset=utf-8',
  '.css': 'text/css; charset=utf-8',
  '.js': 'text/javascript; charset=utf-8',
  '.mjs': 'text/javascript; charset=utf-8',
  '.json': 'application/json',
  '.svg': 'image/svg+xml',
  '.jpg': 'image/jpeg',
  '.jpeg': 'image/jpeg',
  '.png': 'image/png',
  '.webp': 'image/webp',
  '.ico': 'image/x-icon',
  '.txt': 'text/plain; charset=utf-8',
  '.xml': 'application/xml',
  '.webmanifest': 'application/manifest+json',
};

function send(res, code, body, type = 'text/plain; charset=utf-8') {
  res.writeHead(code, { 'Content-Type': type, 'Cache-Control': 'no-store' });
  res.end(body);
}

function proxy(req, res) {
  const target = new URL(req.url, API);
  const headers = { ...req.headers, host: target.host };
  delete headers['content-length'];
  const upstream = http.request(
    {
      protocol: target.protocol,
      hostname: target.hostname,
      port: target.port,
      path: target.pathname + target.search,
      method: req.method,
      headers,
    },
    (up) => {
      res.writeHead(up.statusCode || 502, up.headers);
      up.pipe(res);
    }
  );
  upstream.on('error', (err) => {
    send(
      res,
      502,
      JSON.stringify({
        error: 'ppc_contact_unreachable',
        detail: String(err.message),
        hint: 'Starten Sie ppc-contact auf 127.0.0.1:8690',
      }),
      'application/json'
    );
  });
  req.pipe(upstream);
}

const server = http.createServer((req, res) => {
  const url = new URL(req.url || '/', `http://127.0.0.1:${PORT}`);
  if (url.pathname.startsWith('/api/')) {
    proxy(req, res);
    return;
  }
  let rel = decodeURIComponent(url.pathname);
  if (rel === '/') rel = '/index.html';
  const file = path.normalize(path.join(ROOT, rel));
  if (!file.startsWith(ROOT)) {
    send(res, 403, 'forbidden');
    return;
  }
  fs.readFile(file, (err, data) => {
    if (err) {
      send(res, 404, 'not found');
      return;
    }
    const ext = path.extname(file).toLowerCase();
    send(res, 200, data, TYPES[ext] || 'application/octet-stream');
  });
});

server.listen(PORT, '127.0.0.1', () => {
  console.log(`PPCWeb static+proxy http://127.0.0.1:${PORT}/`);
  console.log(`  root: ${ROOT}`);
  console.log(`  /api → ${API}`);
});

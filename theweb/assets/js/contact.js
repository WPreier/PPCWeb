(() => {
  const form = document.getElementById('contact-form');
  if (!form) return;

  const statusEl = document.getElementById('contact-status');
  const submitBtn = document.getElementById('contact-submit');
  const issuedAt = document.getElementById('issued_at');
  issuedAt.value = String(Math.floor(Date.now() / 1000));

  const setStatus = (msg, kind) => {
    statusEl.textContent = msg || '';
    statusEl.dataset.kind = kind || '';
    statusEl.hidden = !msg;
  };

  form.addEventListener('submit', async (ev) => {
    ev.preventDefault();
    setStatus('', '');

    if (!form.reportValidity()) return;

    const privacy = form.querySelector('#privacy');
    if (!privacy.checked) {
      setStatus('Bitte der Datenschutzhinweis-Zustimmung zustimmen.', 'error');
      return;
    }

    const altchaInput = form.querySelector('input[name="altcha"]');
    if (!altchaInput || !altchaInput.value) {
      setStatus('Bitte warten Sie, bis die Bot-Prüfung abgeschlossen ist.', 'error');
      return;
    }

    const payload = {
      name: form.name.value.trim(),
      email: form.email.value.trim(),
      phone: form.phone.value.trim(),
      topic: form.topic.value,
      message: form.message.value.trim(),
      privacy: true,
      website: form.website.value,
      altcha: altchaInput.value,
      issued_at: Number(issuedAt.value || 0)
    };

    submitBtn.disabled = true;
    setStatus('Wird gesendet …', 'info');

    try {
      const res = await fetch('/api/contact', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json', Accept: 'application/json' },
        body: JSON.stringify(payload),
        credentials: 'same-origin'
      });
      const data = await res.json().catch(() => ({}));
      if (!res.ok) {
        const map = {
          captcha_failed: 'Bot-Prüfung fehlgeschlagen. Bitte Seite neu laden und erneut versuchen.',
          rate_limited: 'Zu viele Versuche. Bitte kurz warten.',
          too_fast: 'Bitte etwas langsamer absenden.',
          form_expired: 'Das Formular ist abgelaufen. Bitte Seite neu laden.',
          privacy_required: 'Datenzustimmung erforderlich.',
          send_failed: 'Versand vorübergehend nicht möglich. Bitte anrufen oder E-Mail schreiben.',
          invalid_email: 'Bitte eine gültige E-Mail-Adresse angeben.',
          invalid_message: 'Bitte eine Nachricht eingeben.',
          invalid_name: 'Bitte Ihren Namen angeben.'
        };
        setStatus(map[data.error] || 'Senden fehlgeschlagen. Bitte später erneut versuchen.', 'error');
        submitBtn.disabled = false;
        return;
      }
      form.reset();
      issuedAt.value = String(Math.floor(Date.now() / 1000));
      setStatus('Danke! Ihre Nachricht ist eingegangen. Wir melden uns persönlich.', 'ok');
      submitBtn.disabled = false;
    } catch (err) {
      setStatus('Netzwerkfehler. Bitte Verbindung prüfen oder anrufen.', 'error');
      submitBtn.disabled = false;
    }
  });
})();

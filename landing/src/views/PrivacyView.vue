<script setup lang="ts">
import { onBeforeUnmount, onMounted } from 'vue'

const PAGE_TITLE = 'Privacy · Takeout Metadata Fixer'

onMounted(() => {
  document.title = PAGE_TITLE
})

onBeforeUnmount(() => {
  document.title = 'Takeout Metadata Fixer'
})
</script>

<template>
  <main class="content-column">
    <article class="article">
      <h1 class="article-title">Privacy policy</h1>
      <p class="article-meta">Last updated: April 12, 2026</p>

      <p>
        This page describes how this website (“the site”) handles information. It does not replace
        legal advice.
      </p>

      <h2>What this site is</h2>
      <p>
        The site is a marketing page for Takeout Metadata Fixer, a desktop application. Using the
        desktop app happens on your computer; this policy focuses on visiting this website.
      </p>

      <h2>Information we collect on this site</h2>
      <p>
        We do not ask you to create an account or submit personal data to use the marketing page.
        Standard web server or hosting logs (for example IP address, user agent, and requested URL)
        may be collected by your hosting provider as part of delivering the site. Refer to your
        host’s documentation for how long those logs are kept.
      </p>

      <h2>Third-party content</h2>
      <p>
        The site may load third-party scripts or widgets (for example a “Buy me a coffee” widget).
        Those services have their own privacy policies and may set cookies or collect usage data
        when their script runs in your browser. Links to GitHub, documentation, or downloads leave
        this site and are governed by those destinations’ policies.
      </p>

      <h2>Cookies</h2>
      <p>
        We do not operate first-party cookies for tracking on this marketing page. Third-party
        embeds may use cookies or similar storage; you can control those through your browser
        settings.
      </p>

      <h2>Children</h2>
      <p>
        The site is not directed at children under 13, and we do not knowingly collect their data.
      </p>

      <h2>Changes</h2>
      <p>
        We may update this policy from time to time. The “Last updated” date at the top will change.
      </p>

      <h2>Contact</h2>
      <p>
        For privacy questions about this site, you can reach out via the contact options linked from
        <a class="text-link" href="https://mrashidi.me" target="_blank" rel="noopener noreferrer"
          >mrashidi.me</a
        >.
      </p>
    </article>
  </main>
</template>

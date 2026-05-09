<script setup lang="ts">
const REPO_URL = 'https://github.com/MRdevX/takeout-md-fixer'
const SITE_URL = 'https://mrashidi.me'
const COPYRIGHT_HOLDER = 'Mahdi Rashidi'

const year = new Date().getFullYear()
</script>

<template>
  <div class="landing-bg" aria-hidden="true" />
  <div class="page">
    <header class="site-header content-column">
      <RouterLink to="/" class="text-link header-link header-brand"
        >Takeout Metadata Fixer</RouterLink
      >
      <a
        :href="REPO_URL"
        class="btn btn--ghost btn--sm header-github"
        target="_blank"
        rel="noopener noreferrer"
        >GitHub</a
      >
    </header>

    <RouterView />

    <footer class="site-footer">
      <div class="footer-inner content-column">
        <nav class="footer-nav" aria-label="Footer">
          <ul class="footer-links">
            <li>
              <RouterLink class="footer-link" to="/privacy">Privacy</RouterLink>
            </li>
            <li>
              <a class="footer-link" :href="SITE_URL" target="_blank" rel="noopener noreferrer"
                >mrashidi.me</a
              >
            </li>
            <li>
              <a class="footer-link" :href="REPO_URL" target="_blank" rel="noopener noreferrer">GitHub</a>
            </li>
          </ul>
        </nav>
        <p class="footer-legal">© {{ year }} {{ COPYRIGHT_HOLDER }}. All rights reserved.</p>
      </div>
    </footer>
  </div>
</template>

<style scoped>
.header-brand {
  margin-right: auto;
  font-weight: 600;
  font-family: var(--font-display);
  font-size: var(--text-base);
  letter-spacing: 0.02em;
}

.header-github {
  flex-shrink: 0;
}
</style>

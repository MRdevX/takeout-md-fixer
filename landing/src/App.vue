<script setup lang="ts">
const REPO_URL = 'https://github.com/MRdevX/takeout-md-fixer'
const SITE_URL = 'https://mrashidi.me'
const COPYRIGHT_HOLDER = 'Mahdi Rashidi'

const year = new Date().getFullYear()
</script>

<template>
  <div class="page">
    <header class="site-header content-column">
      <RouterLink to="/" class="header-brand">Takeout Metadata Fixer</RouterLink>
      <a :href="REPO_URL" class="btn btn--outline" target="_blank" rel="noopener noreferrer"
        >GitHub</a
      >
    </header>

    <RouterView />

    <footer class="site-footer">
      <div class="footer-inner content-column">
        <nav aria-label="Footer">
          <ul class="footer-links">
            <li><RouterLink class="footer-link" to="/privacy">Privacy</RouterLink></li>
            <li>
              <a class="footer-link" :href="SITE_URL" target="_blank" rel="noopener noreferrer"
                >mrashidi.me</a
              >
            </li>
            <li>
              <a class="footer-link" :href="REPO_URL" target="_blank" rel="noopener noreferrer"
                >GitHub</a
              >
            </li>
          </ul>
        </nav>
        <p class="footer-legal">© {{ year }} {{ COPYRIGHT_HOLDER }}. Free and open source.</p>
      </div>
    </footer>
  </div>
</template>

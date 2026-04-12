<script setup lang="ts">
const RELEASES_URL = 'https://github.com/MRdevX/takeout-md-fixer/releases'
const REPO_URL = 'https://github.com/MRdevX/takeout-md-fixer'
const EXIFTOOL_URL = 'https://exiftool.org/'
const BUY_ME_A_COFFEE_URL = 'https://www.buymeacoffee.com/mrdevx'
</script>

<template>
  <div class="page">
    <header class="site-header content-column">
      <span class="site-title">Takeout Metadata Fixer</span>
      <a
        :href="REPO_URL"
        class="about-link toolbar-btn"
        target="_blank"
        rel="noopener noreferrer"
      >
        GitHub
      </a>
    </header>

    <main class="content-column">
      <div class="hero">
        <h1>Takeout Metadata Fixer</h1>
        <p class="subtitle">
          Copies date and location from Google Takeout <code class="inline-code">.json</code> sidecars into your photos and
          videos so other apps show them correctly.
        </p>
        <p class="subtitle-trust">
          <strong>Writes to your files.</strong> Use a copy of your library if you want to be safe.
        </p>

        <div class="landing-note" role="note">
          <strong>ExifTool</strong> must be installed on your system; the desktop app uses it to update metadata. See
          <a class="text-link" :href="EXIFTOOL_URL" target="_blank" rel="noopener noreferrer">exiftool.org</a> for install
          steps.
        </div>

        <div class="landing-actions">
          <a :href="RELEASES_URL" class="btn btn-primary" target="_blank" rel="noopener noreferrer">Download</a>
          <a :href="REPO_URL" class="about-link toolbar-btn" target="_blank" rel="noopener noreferrer">View source</a>
        </div>

        <p class="landing-tertiary">
          Binaries and release notes:
          <a class="text-link" :href="RELEASES_URL" target="_blank" rel="noopener noreferrer">GitHub Releases</a>
        </p>
        <p class="landing-tertiary">
          <a class="text-link" :href="BUY_ME_A_COFFEE_URL" target="_blank" rel="noopener noreferrer">Buy me a coffee</a>
          if you find this useful.
        </p>
      </div>
    </main>
  </div>
</template>

<style scoped>
.inline-code {
  font-size: 0.88em;
  background: var(--bg-subtle);
  padding: 1px 6px;
  border-radius: 4px;
  font-family: ui-monospace, SFMono-Regular, Menlo, monospace;
}
</style>

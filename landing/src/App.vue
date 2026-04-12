<script setup lang="ts">
const RELEASES_URL = 'https://github.com/MRdevX/takeout-md-fixer/releases'
const REPO_URL = 'https://github.com/MRdevX/takeout-md-fixer'
const EXIFTOOL_URL = 'https://exiftool.org/'
const BUY_ME_A_COFFEE_URL = 'https://www.buymeacoffee.com/mrdevx'
</script>

<template>
  <div class="page">
    <header class="site-header content-column">
      <a :href="REPO_URL" class="text-link header-link" target="_blank" rel="noopener noreferrer">GitHub</a>
    </header>

    <main class="content-column">
      <div class="hero">
        <h1>Takeout Metadata Fixer</h1>
        <p class="lead">
          Desktop app for Google Takeout: reads <code class="inline-code">.json</code> sidecars next to your photos and
          videos and writes dates, GPS, and related metadata into the files with ExifTool. That way imports into iCloud, a
          NAS, or other tools see sensible dates and locations.
        </p>
        <p class="lead">
          <strong>Writes to your files.</strong> Use a copy of your library if you want to be safe.
        </p>

        <aside class="landing-note" aria-label="ExifTool requirement">
          <strong>ExifTool</strong> must be installed on your system; the app uses it to update metadata. See
          <a class="text-link" :href="EXIFTOOL_URL" target="_blank" rel="noopener noreferrer">exiftool.org</a> for install
          steps.
        </aside>

        <div class="landing-actions">
          <a :href="RELEASES_URL" class="btn btn-primary" target="_blank" rel="noopener noreferrer">Download</a>
        </div>
        <p class="landing-hint">Prebuilt macOS (DMG) and Windows builds on GitHub Releases.</p>

        <p class="landing-tertiary">
          <a class="text-link" :href="BUY_ME_A_COFFEE_URL" target="_blank" rel="noopener noreferrer">Buy me a coffee</a>
          if you find this useful.
        </p>
      </div>
    </main>
  </div>
</template>

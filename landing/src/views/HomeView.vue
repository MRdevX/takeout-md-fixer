<script setup lang="ts">
import { computed, ref } from 'vue'

const RELEASES_URL = 'https://github.com/MRdevX/takeout-md-fixer/releases'
const GITHUB_ISSUES_URL = 'https://github.com/MRdevX/takeout-md-fixer/issues'
const REPO_URL = 'https://github.com/MRdevX/takeout-md-fixer'
const SITE_URL = 'https://mrashidi.me'

const SCREENSHOTS = [
  {
    src: 'https://ecuc8c19st5krt52.public.blob.vercel-storage.com/step1-main.png',
    alt: 'Choose a Takeout folder before starting',
  },
  {
    src: 'https://ecuc8c19st5krt52.public.blob.vercel-storage.com/step2-review.png',
    alt: 'Review scan summary before starting; optional delete JSON sidecars after success',
  },
  {
    src: 'https://ecuc8c19st5krt52.public.blob.vercel-storage.com/step3-progress.png',
    alt: 'Progress while metadata is written to files; pause or stop available',
  },
  {
    src: 'https://ecuc8c19st5krt52.public.blob.vercel-storage.com/step4-done.png',
    alt: 'Finished: succeeded, skipped, and failed counts',
  },
] as const

const slideIndex = ref(0)

const touchStartX = ref(0)
const touchStartY = ref(0)

const slideCount = SCREENSHOTS.length
const currentSlide = computed(() => slideIndex.value + 1)

function goPrev() {
  slideIndex.value = (slideIndex.value - 1 + slideCount) % slideCount
}

function goNext() {
  slideIndex.value = (slideIndex.value + 1) % slideCount
}

function goToSlide(i: number) {
  slideIndex.value = i
}

function onTouchStart(e: TouchEvent) {
  const t = e.touches[0]
  if (!t) return
  touchStartX.value = t.clientX
  touchStartY.value = t.clientY
}

function onTouchEnd(e: TouchEvent) {
  const t = e.changedTouches[0]
  if (!t) return
  const dx = t.clientX - touchStartX.value
  const dy = t.clientY - touchStartY.value
  const adx = Math.abs(dx)
  const ady = Math.abs(dy)
  if (adx < 56 || adx < ady * 1.15) return
  if (dx < 0) goNext()
  else goPrev()
}

function onCarouselKeydown(e: KeyboardEvent) {
  if (e.key === 'ArrowLeft') {
    e.preventDefault()
    goPrev()
  } else if (e.key === 'ArrowRight') {
    e.preventDefault()
    goNext()
  } else if (e.key === 'Home') {
    e.preventDefault()
    goToSlide(0)
  } else if (e.key === 'End') {
    e.preventDefault()
    goToSlide(slideCount - 1)
  }
}
</script>

<template>
  <main class="content-column">
    <div class="hero-wrap">
      <section class="hero-section">
        <div class="hero-background" aria-hidden="true" />
        <div class="hero-gradient-bg" aria-hidden="true" />
        <div class="hero-inner">
          <div class="hero-heading">
            <h1 class="cyberpunk-h1">Takeout Metadata Fixer</h1>
            <p class="hero-tagline">
              Copies dates and GPS from Google Takeout JSON files into your photos and videos.
            </p>
          </div>
          <div class="hero-prose">
            <p class="lead">
              Takeout gives you JSON files next to your photos and videos with the taken time, GPS, and a few
              other fields. A lot of programs only read the media file, so dates and locations look wrong even
              though the export is fine. This app copies that data into the files. After that, things like
              iCloud, a NAS import, or your usual photo app can show the dates and places you expect.
            </p>
            <p class="lead">
              The app edits your originals on disk, so copy the folder first if you want to keep an untouched
              backup.
            </p>
            <p class="lead">
              You can pause or stop mid-run and resume later from a checkpoint file in that folder. After a
              successful run you can optionally remove the JSON sidecars. The last screen shows how many files
              succeeded, were skipped, or failed.
            </p>
          </div>
        </div>
      </section>

      <section class="screens-section" aria-labelledby="carousel-heading">
        <h2 id="carousel-heading" class="section-label">
          <svg
            class="section-label-icon"
            xmlns="http://www.w3.org/2000/svg"
            width="20"
            height="20"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
            aria-hidden="true"
          >
            <rect width="18" height="18" x="3" y="3" rx="2" ry="2" />
            <circle cx="9" cy="9" r="2" />
            <path d="m21 15-3.086-3.086a2 2 0 0 0-2.828 0L6 21" />
          </svg>
          Screenshots
        </h2>
        <div class="terminal-window">
          <div class="terminal-window-inner">
            <section
              class="screens-carousel"
              aria-roledescription="carousel"
              tabindex="0"
              @keydown="onCarouselKeydown"
            >
              <div class="carousel-viewport" @touchstart.passive="onTouchStart" @touchend="onTouchEnd">
                <button
                  type="button"
                  class="btn btn--icon btn--media carousel-nav carousel-nav-prev"
                  aria-label="Previous screenshot"
                  @click="goPrev"
                >
                  <span aria-hidden="true">‹</span>
                </button>
                <div
                  class="carousel-track"
                  :style="{ transform: `translate3d(-${slideIndex * 100}%, 0, 0)` }"
                  role="group"
                  :aria-label="`Slide ${currentSlide} of ${slideCount}`"
                >
                  <div
                    v-for="(slide, i) in SCREENSHOTS"
                    :key="slide.src"
                    class="carousel-slide"
                    :aria-hidden="i !== slideIndex"
                  >
                    <img
                      :src="slide.src"
                      :alt="slide.alt"
                      :loading="i === 0 ? 'eager' : 'lazy'"
                      decoding="async"
                      class="carousel-img"
                    />
                  </div>
                </div>
                <button
                  type="button"
                  class="btn btn--icon btn--media carousel-nav carousel-nav-next"
                  aria-label="Next screenshot"
                  @click="goNext"
                >
                  <span aria-hidden="true">›</span>
                </button>
              </div>
              <p class="carousel-status" aria-live="polite">{{ currentSlide }} of {{ slideCount }}</p>
              <div class="carousel-dots" role="tablist" aria-label="Screenshot thumbnails">
                <button
                  v-for="(slide, i) in SCREENSHOTS"
                  :key="'dot-' + slide.src"
                  type="button"
                  role="tab"
                  class="carousel-dot"
                  :class="{ 'is-active': i === slideIndex }"
                  :aria-selected="i === slideIndex"
                  :tabindex="i === slideIndex ? 0 : -1"
                  :aria-label="`Screenshot ${i + 1}: ${slide.alt}`"
                  @click="goToSlide(i)"
                />
              </div>
            </section>
          </div>
        </div>
      </section>

      <section class="content-section cta-panel">
        <div class="cta-section-bg" aria-hidden="true" />
        <div class="cta-content">
          <h2 class="cta-title">Download</h2>
          <p class="cta-description">
            Signed builds for macOS (DMG) and Windows (.exe) are attached to each release on GitHub. Open the
            latest one and download the installer for your system.
          </p>
          <div class="landing-actions">
            <a
              :href="RELEASES_URL"
              class="btn btn--primary"
              target="_blank"
              rel="noopener noreferrer"
            >
              Download from GitHub
            </a>
          </div>
          <p class="landing-hint">
            Every release lists checksums so you can confirm the installer before you run it.
            <a
              class="text-link text-link--prose"
              :href="GITHUB_ISSUES_URL"
              target="_blank"
              rel="noopener noreferrer"
              >Open an issue on GitHub</a
            >
            if anything misbehaves.
          </p>
          <p class="landing-hint">
            If the app saved you a headache:
            <a class="text-link text-link--prose" :href="REPO_URL" target="_blank" rel="noopener noreferrer"
              >Star the repo on GitHub</a
            >
            so others can find it, or say thanks with a coffee below—stars, caffeine, both, or neither. Thank you
            either way.
          </p>
          <p class="landing-hint">
            For other projects and more about me, visit
            <a class="text-link text-link--prose" :href="SITE_URL" target="_blank" rel="noopener noreferrer"
              >mrashidi.me</a
            >.
          </p>
          <div class="landing-actions cta-bmc-wrap">
            <iframe
              class="cta-bmc-iframe"
              src="/buy-me-a-coffee-embed.html"
              title="Buy me a coffee"
              loading="lazy"
              referrerpolicy="no-referrer-when-downgrade"
            />
          </div>
        </div>
      </section>
    </div>
  </main>
</template>

<style scoped>
.cta-bmc-wrap {
    align-items: center;
}

.cta-bmc-iframe {
    display: block;
    margin-inline: auto;
    width: 216px;
    height: 62px;
    max-width: 100%;
    border: none;
    outline: none;
    background: transparent;
    box-shadow: none;
    overflow: hidden;
    color-scheme: dark;
}

.screens-carousel {
    position: relative;
    width: 100%;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: var(--stack-md);
    outline: none;
}

.carousel-viewport {
    position: relative;
    width: 100%;
    border-radius: var(--radius-inner);
    overflow: hidden;
    touch-action: pan-y pinch-zoom;
}

.carousel-track {
    display: flex;
    backface-visibility: hidden;
    transform: translateZ(0);
    transition: transform 0.48s var(--ease-smooth);
}

@media (prefers-reduced-motion: reduce) {
    .carousel-track {
        transition: none;
    }
}

.carousel-slide {
    flex: 0 0 100%;
    min-width: 0;
}

.carousel-img {
    display: block;
    width: 100%;
    height: auto;
    max-height: var(--screenshot-max-height);
    object-fit: contain;
    object-position: top center;
    border-radius: var(--radius-inner);
    border: 1px solid rgb(51 65 85 / 0.45);
    background: rgb(10 12 16 / 0.5);
}

.carousel-nav {
    position: absolute;
    top: 50%;
    z-index: 2;
    translate: 0 -50%;
    width: var(--carousel-nav-size);
    height: var(--carousel-nav-size);
    font-size: 1.25rem;
    line-height: 1;
}

.carousel-nav:focus-visible {
    z-index: 3;
}

.carousel-nav-prev {
    left: max(var(--space-2), env(safe-area-inset-left));
}

.carousel-nav-next {
    right: max(var(--space-2), env(safe-area-inset-right));
}

.carousel-status {
    position: absolute;
    width: 1px;
    height: 1px;
    padding: 0;
    margin: -1px;
    overflow: hidden;
    clip: rect(0, 0, 0, 0);
    clip-path: inset(50%);
    white-space: nowrap;
    border: 0;
}

.carousel-dots {
    display: flex;
    flex-wrap: wrap;
    align-items: center;
    justify-content: center;
    gap: var(--space-2);
    padding-block: var(--space-2);
}

.carousel-dot {
    min-width: var(--touch-target);
    min-height: var(--touch-target);
    padding: 0;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    border: none;
    background: transparent;
    cursor: pointer;
    -webkit-tap-highlight-color: transparent;
}

.carousel-dot::after {
    content: "";
    width: var(--space-2);
    height: var(--space-2);
    border-radius: 999px;
    background: #64748b;
    transition:
        background var(--duration-fast) var(--ease-out),
        transform var(--duration-fast) var(--ease-out),
        box-shadow var(--duration-fast) var(--ease-out);
}

.carousel-dot:hover::after {
    background: #94a3b8;
}

.carousel-dot.is-active::after {
    background: var(--neon-orange);
    transform: scale(1.25);
    box-shadow: 0 0 10px rgb(255 95 31 / 0.35);
}

@media (prefers-reduced-motion: reduce) {
    .carousel-dot.is-active::after {
        transform: none;
    }
}

.carousel-dot:focus-visible {
    outline: var(--focus-ring);
    outline-offset: var(--focus-offset);
    border-radius: var(--radius-control);
}
</style>

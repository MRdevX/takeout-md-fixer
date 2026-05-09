<script setup lang="ts">
import { computed, ref } from 'vue'

const RELEASES_URL = 'https://github.com/MRdevX/takeout-md-fixer/releases'
const EXIFTOOL_URL = 'https://exiftool.org/'

const SCREENSHOTS = [
  {
    src: 'https://ecuc8c19st5krt52.public.blob.vercel-storage.com/step1-main.png',
    alt: 'Takeout Metadata Fixer — welcome screen with folder selection',
  },
  {
    src: 'https://ecuc8c19st5krt52.public.blob.vercel-storage.com/step2-review.png',
    alt: 'Review step showing folder summary and counts',
  },
  {
    src: 'https://ecuc8c19st5krt52.public.blob.vercel-storage.com/step3-progress.png',
    alt: 'Update in progress with progress bar',
  },
  {
    src: 'https://ecuc8c19st5krt52.public.blob.vercel-storage.com/step4-done.png',
    alt: 'Done screen with succeeded, skipped, and failed counts',
  },
] as const

const slideIndex = ref(0)

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
          <h1 class="cyberpunk-h1">Takeout Metadata Fixer</h1>
          <p class="hero-tagline">Desktop · Google Takeout · ExifTool</p>
          <p class="lead text-body">
            Desktop app for Google Takeout: reads <code class="inline-code">.json</code> sidecars next to
            your photos and videos and writes dates, GPS, and related metadata into the files with
            ExifTool. That way imports into iCloud, a NAS, or other tools see sensible dates and
            locations.
          </p>
          <p class="lead">
            <strong>Writes to your files.</strong> Use a copy of your library if you want to be safe.
          </p>
        </div>
      </section>

      <aside class="landing-note content-section" aria-label="ExifTool requirement">
        <strong>ExifTool</strong> must be installed on your system; the app uses it to update metadata.
        See
        <a class="text-link" :href="EXIFTOOL_URL" target="_blank" rel="noopener noreferrer"
          >exiftool.org</a
        >
        for install steps.
      </aside>

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
              <div class="carousel-viewport">
                <button
                  type="button"
                  class="carousel-nav carousel-nav-prev"
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
                  class="carousel-nav carousel-nav-next"
                  aria-label="Next screenshot"
                  @click="goNext"
                >
                  <span aria-hidden="true">›</span>
                </button>
              </div>
              <p class="carousel-status" aria-live="polite">{{ currentSlide }} of {{ slideCount }}</p>
              <div class="carousel-dots" role="tablist" aria-label="Choose screenshot">
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
          <h2 class="cta-title">Ready to fix metadata</h2>
          <p class="cta-description">
            Prebuilt macOS (DMG) and Windows builds are published on GitHub Releases.
          </p>
          <div class="landing-actions">
            <a :href="RELEASES_URL" class="btn" target="_blank" rel="noopener noreferrer">
              Download
            </a>
          </div>
          <p class="landing-hint">Verify checksums on the release page when you download.</p>
        </div>
      </section>
    </div>
  </main>
</template>

<style scoped>
.screens-carousel {
  position: relative;
  width: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--space-3);
}

.carousel-viewport {
  position: relative;
  width: 100%;
  border-radius: var(--radius-sm);
  overflow: hidden;
}

.carousel-track {
  display: flex;
  transition: transform 0.38s cubic-bezier(0.22, 1, 0.36, 1);
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
  border-radius: var(--radius-sm);
  border: 1px solid rgb(255 95 31 / 0.2);
}

.carousel-nav {
  position: absolute;
  top: 50%;
  z-index: 1;
  translate: 0 -50%;
  width: 2.5rem;
  height: 2.5rem;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 0;
  border: 1px solid rgb(255 95 31 / 0.35);
  border-radius: 999px;
  background: rgb(17 19 24 / 0.82);
  color: rgb(var(--foreground-rgb));
  font-size: 1.35rem;
  line-height: 1;
  cursor: pointer;
  backdrop-filter: blur(8px);
  transition:
    background 0.15s,
    color 0.15s,
    border-color 0.15s,
    box-shadow 0.15s;
}

.carousel-nav:hover {
  background: rgb(255 95 31 / 0.12);
  color: var(--neon-orange);
  border-color: rgb(255 95 31 / 0.55);
  box-shadow: 0 0 16px rgb(255 95 31 / 0.2);
}

.carousel-nav:focus-visible {
  outline: var(--focus-ring);
  outline-offset: var(--focus-offset);
}

.carousel-nav-prev {
  left: var(--space-2);
}

.carousel-nav-next {
  right: var(--space-2);
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
}

.carousel-dot {
  width: 0.5rem;
  height: 0.5rem;
  padding: 0;
  border: none;
  border-radius: 999px;
  background: #64748b;
  cursor: pointer;
  transition:
    background 0.15s,
    transform 0.15s;
}

.carousel-dot:hover {
  background: #94a3b8;
}

.carousel-dot.is-active {
  background: var(--neon-orange);
  transform: scale(1.15);
}

.carousel-dot:focus-visible {
  outline: var(--focus-ring);
  outline-offset: var(--focus-offset);
}

@media (max-width: 480px) {
  .carousel-nav {
    width: 2.25rem;
    height: 2.25rem;
    font-size: 1.2rem;
  }

  .carousel-nav-prev {
    left: var(--space-1);
  }

  .carousel-nav-next {
    right: var(--space-1);
  }
}
</style>

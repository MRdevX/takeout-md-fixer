<script setup lang="ts">
import { computed, ref } from 'vue'

const RELEASES_URL = 'https://github.com/MRdevX/takeout-md-fixer/releases'
const GITHUB_ISSUES_URL = 'https://github.com/MRdevX/takeout-md-fixer/issues'
const REPO_URL = 'https://github.com/MRdevX/takeout-md-fixer'
const SITE_URL = 'https://mrashidi.me'

const SCREENSHOTS = [
  {
    src: 'https://ecuc8c19st5krt52.public.blob.vercel-storage.com/step1-main.png',
    alt: 'The welcome screen, before choosing a Takeout folder',
  },
  {
    src: 'https://ecuc8c19st5krt52.public.blob.vercel-storage.com/step2-review.png',
    alt: 'The review screen: how many files are ready, and the option to clean up .json files afterwards',
  },
  {
    src: 'https://ecuc8c19st5krt52.public.blob.vercel-storage.com/step3-progress.png',
    alt: 'Progress while the app writes dates and locations, with pause and stop',
  },
  {
    src: 'https://ecuc8c19st5krt52.public.blob.vercel-storage.com/step4-done.png',
    alt: 'The finished screen, with how many photos were fixed, skipped, and failed',
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
  <main>
    <div class="content-column">
      <section class="hero">
        <div class="polaroids" aria-hidden="true">
          <div class="polaroid">
            <span class="polaroid-film" /><span class="polaroid-cap">NO DATE</span>
          </div>
          <div class="polaroid">
            <span class="polaroid-film" /><span class="polaroid-cap polaroid-cap--fixed"
              >JUL 14 2019<span class="polaroid-place"> · ROME</span></span
            >
          </div>
          <div class="polaroid">
            <span class="polaroid-film" /><span class="polaroid-cap">NO DATE</span>
          </div>
        </div>

        <h1>Give your photos their dates back.</h1>
        <p class="hero-tagline">
          Google Takeout saves dates and locations in separate <code>.json</code> files. This app
          puts that info back inside your photos and videos, so they show up correctly in any app.
        </p>

        <div class="hero-actions">
          <a
            :href="RELEASES_URL"
            class="btn btn--primary"
            target="_blank"
            rel="noopener noreferrer"
          >
            Download for macOS or Windows
          </a>
        </div>

        <p class="hero-note">
          <span>Free and open source.</span>
          <span class="hero-note-dot" aria-hidden="true" />
          <span>It changes your files, so keep a backup to be extra safe.</span>
        </p>
      </section>

      <section class="section" aria-labelledby="why-heading">
        <p class="section-label">WHY YOUR DATES LOOK WRONG</p>
        <h2 id="why-heading" class="section-title">
          The dates are not missing. They are in the wrong place.
        </h2>
        <p class="section-lead">
          A Takeout export puts each photo's taken time and location in a small
          <code>.json</code> file next to it, instead of inside the photo. Most apps only read the
          photo, so everything lands on the day you downloaded it. This app moves that info back
          where it belongs, using the free
          <a
            class="text-link"
            href="https://exiftool.org/"
            target="_blank"
            rel="noopener noreferrer"
            >ExifTool</a
          >
          program. After that, Photos, iCloud, a NAS import, or whatever you use shows the right day
          and place.
        </p>
      </section>

      <section class="section" aria-labelledby="how-heading">
        <p class="section-label">HOW IT GOES</p>
        <h2 id="how-heading" class="section-title">Three steps, then you are done.</h2>
        <ol class="steps">
          <li class="step">
            <span class="step-num" aria-hidden="true">1</span>
            <p class="step-title">Choose your folder</p>
            <p class="step-text">
              Point the app at your unzipped Takeout folder. It reads everything inside.
            </p>
          </li>
          <li class="step">
            <span class="step-num" aria-hidden="true">2</span>
            <p class="step-title">See what it found</p>
            <p class="step-text">
              How many files are ready, how many have nothing to restore, and whether to delete the
              <code>.json</code> files afterwards.
            </p>
          </li>
          <li class="step">
            <span class="step-num" aria-hidden="true">3</span>
            <p class="step-title">Let it run</p>
            <p class="step-text">
              Pause or stop whenever you like. The app remembers where it left off, so you can
              finish later.
            </p>
          </li>
        </ol>
        <ul class="chips" aria-label="What the last screen tells you">
          <li class="chip">
            <span class="dot dot-ok" aria-hidden="true" />Fixed, with dates and places back
          </li>
          <li class="chip">
            <span class="dot dot-warn" aria-hidden="true" />Skipped, nothing to change
          </li>
          <li class="chip">
            <span class="dot dot-err" aria-hidden="true" />Failed, run again to retry
          </li>
        </ul>
      </section>

      <section class="section" aria-labelledby="shots-heading">
        <p class="section-label">SCREENSHOTS</p>
        <h2 id="shots-heading" class="section-title">A look inside.</h2>
        <div class="shots-frame">
          <section
            class="carousel"
            aria-roledescription="carousel"
            tabindex="0"
            @keydown="onCarouselKeydown"
          >
            <div
              class="carousel-viewport"
              @touchstart.passive="onTouchStart"
              @touchend="onTouchEnd"
            >
              <button
                type="button"
                class="btn btn--icon carousel-nav carousel-nav-prev"
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
                class="btn btn--icon carousel-nav carousel-nav-next"
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
      </section>

      <section class="download" aria-labelledby="download-heading">
        <h2 id="download-heading" class="section-title">Get it.</h2>
        <p class="section-lead">
          Signed builds for macOS and Windows are attached to every release on GitHub. Open the
          latest one and take the installer for your system.
        </p>
        <a :href="RELEASES_URL" class="btn btn--primary" target="_blank" rel="noopener noreferrer">
          Download from GitHub
        </a>
        <p class="download-hint">
          Every release lists checksums, so you can check the installer before you run it. If
          something misbehaves,
          <a class="text-link" :href="GITHUB_ISSUES_URL" target="_blank" rel="noopener noreferrer"
            >open an issue</a
          >.
        </p>
        <p class="download-hint">
          If it saved you a headache, you can
          <a class="text-link" :href="REPO_URL" target="_blank" rel="noopener noreferrer"
            >star the repo</a
          >
          so others find it, or say thanks with a coffee. Either is lovely, neither is expected. For
          my other work, there is
          <a class="text-link" :href="SITE_URL" target="_blank" rel="noopener noreferrer"
            >mrashidi.me</a
          >.
        </p>
        <iframe
          class="bmc-frame"
          src="/buy-me-a-coffee-embed.html"
          title="Buy me a coffee"
          loading="lazy"
          referrerpolicy="no-referrer-when-downgrade"
        />
      </section>
    </div>
  </main>
</template>

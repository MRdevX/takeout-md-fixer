import { Browser, Events } from "@wailsio/runtime";
import { MetadataService } from "../bindings/takeout-md-fixer/internal/service";

// macOS: keep custom top bar below traffic lights (see main.go InvisibleTitleBarHeight).
if (typeof navigator !== "undefined" && navigator.userAgent.includes("Macintosh")) {
  document.documentElement.classList.add("platform-macos");
}

const views = {
  welcome: document.getElementById("view-welcome"),
  scan: document.getElementById("view-scan"),
  processing: document.getElementById("view-processing"),
  done: document.getElementById("view-done"),
};

const STEP_LABELS = {
  welcome: "",
  scan: "STEP 2 OF 3 · REVIEW",
  processing: "STEP 3 OF 3 · FIXING",
  done: "",
};

const stepLabel = document.getElementById("step-label");
const btnBack = document.getElementById("btn-back");
const btnFix = document.getElementById("btn-fix");

const processingTitle = document.getElementById("processing-title");
const processingSubtitle = document.getElementById("processing-subtitle");
const progressPercentEl = document.getElementById("progress-percent");
const progressBarEl = document.getElementById("progress-bar");
const progressSegments = progressBarEl ? Array.from(progressBarEl.querySelectorAll(".segment")) : [];
const progressTextEl = document.getElementById("progress-text");
const progressFileEl = document.getElementById("progress-file");
const progressFileNameEl = document.getElementById("progress-file-name");
const progressFilePathEl = document.getElementById("progress-file-path");
const btnFixPause = document.getElementById("btn-fix-pause");
const btnFixResume = document.getElementById("btn-fix-resume");
const btnFixStop = document.getElementById("btn-fix-stop");
const checkpointHint = document.getElementById("checkpoint-hint");
const scanEmptyMessage = document.getElementById("scan-empty-message");

const PROCESSING_SUBTITLE_DEFAULT = "You can stop anytime. The app remembers where it left off.";
const PROCESSING_SUBTITLE_PAUSED = "Paused. Resume picks up where it stopped.";

/** @param {number} n */
function num(n) {
  return Number(n || 0).toLocaleString();
}

/** @param {number} n @param {string} one @param {string} many */
function plural(n, one, many) {
  return Number(n) === 1 ? one : many;
}

/** @param {string} filePath */
function basename(filePath) {
  if (!filePath) return "";
  const normalized = filePath.replace(/\\/g, "/");
  const parts = normalized.split("/");
  const last = parts.pop();
  return last || filePath;
}

function showView(name) {
  Object.values(views).forEach((v) => v.classList.remove("active"));
  views[name].classList.add("active");
  if (stepLabel) stepLabel.textContent = STEP_LABELS[name] ?? "";
}

function setProcessingPaused(paused) {
  if (!processingTitle || !btnFixPause || !btnFixResume) return;
  processingTitle.textContent = paused ? "Paused" : "Processing…";
  if (processingSubtitle) {
    processingSubtitle.textContent = paused ? PROCESSING_SUBTITLE_PAUSED : PROCESSING_SUBTITLE_DEFAULT;
  }
  btnFixPause.classList.toggle("hidden", paused);
  btnFixResume.classList.toggle("hidden", !paused);
  if (stepLabel) stepLabel.textContent = paused ? "STEP 3 OF 3 · PAUSED" : STEP_LABELS.processing;
}

function resetProcessingControls() {
  setProcessingPaused(false);
}

/** Light up whole segments, then part-fill the one the progress is inside. */
function paintSegments(pct) {
  const filled = (pct / 100) * progressSegments.length;
  progressSegments.forEach((seg, i) => {
    seg.classList.toggle("is-full", filled >= i + 1);
    seg.classList.toggle("is-partial", filled > i && filled < i + 1);
  });
}

function setProgressUi(current, total, file = "") {
  const safeTotal = Number(total) > 0 ? Number(total) : 0;
  const safeCurrent = Math.max(0, Number(current) || 0);
  const boundedCurrent = safeTotal > 0 ? Math.min(safeCurrent, safeTotal) : safeCurrent;
  const pct = safeTotal > 0 ? Math.round((boundedCurrent / safeTotal) * 100) : 0;

  if (progressPercentEl) progressPercentEl.textContent = String(pct);
  paintSegments(pct);

  if (progressBarEl) {
    progressBarEl.setAttribute("aria-valuenow", String(pct));
    const valText = safeTotal > 0 ? `${boundedCurrent} of ${safeTotal} files, ${pct} percent` : `${pct} percent`;
    progressBarEl.setAttribute("aria-valuetext", valText);
  }
  if (progressTextEl) {
    progressTextEl.textContent = `${num(boundedCurrent)} of ${num(safeTotal)} done`;
  }
  const base = basename(file);
  if (progressFileEl) progressFileEl.classList.toggle("hidden", !file);
  if (progressFileNameEl) progressFileNameEl.textContent = base;
  if (progressFilePathEl) progressFilePathEl.textContent = base === file ? "" : file;
}

let currentPath = "";
let scanData = null;
/** Whether ExifTool was found (PATH plus common install locations). */
let exiftoolOk = true;

const aboutModal = document.getElementById("about-modal");
const aboutModalCloseBtn = document.getElementById("about-modal-close");
const aboutBtn = document.getElementById("btn-about");
const exiftoolWarningEl = document.getElementById("exiftool-warning");
const scanErrorBanner = document.getElementById("scan-error-banner");
const scanErrorMsg = document.getElementById("scan-error-msg");
let aboutLastFocused = null;

function hideScanErrorBanner() {
  scanErrorBanner?.classList.add("hidden");
}

function showScanErrorBanner(message) {
  if (!scanErrorBanner || !scanErrorMsg) return;
  scanErrorMsg.textContent = message;
  scanErrorBanner.classList.remove("hidden");
}

/** @param {unknown} err */
function formatErrorMessage(err) {
  if (err == null) return "Something went wrong.";
  if (typeof err === "string") return err;
  if (typeof err === "object" && err !== null && "message" in err && typeof err.message === "string" && err.message) {
    return err.message;
  }
  try {
    return JSON.stringify(err);
  } catch {
    return String(err);
  }
}

async function updateCheckpointHint() {
  if (!checkpointHint) return;
  if (!currentPath) {
    checkpointHint.classList.add("hidden");
    return;
  }
  try {
    const ok = await MetadataService.FixCheckpointAvailable(currentPath);
    if (ok) {
      checkpointHint.textContent = "Picking up where you left off. Files already fixed are skipped.";
      checkpointHint.classList.remove("hidden");
    } else {
      checkpointHint.classList.add("hidden");
    }
  } catch (e) {
    console.error("FixCheckpointAvailable error:", e);
    checkpointHint.classList.add("hidden");
  }
}

async function refreshExiftoolStatus() {
  try {
    const st = await MetadataService.ExiftoolCheck();
    exiftoolOk = st.ok;
    if (!st.ok) {
      exiftoolWarningEl.classList.remove("hidden");
      const msgEl = exiftoolWarningEl.querySelector(".exiftool-banner-msg");
      if (msgEl) {
        msgEl.textContent = st.message || "One thing is missing: ExifTool";
      }
    } else {
      exiftoolWarningEl.classList.add("hidden");
    }
    if (scanData) {
      renderScanResults(scanData);
    }
    await updateCheckpointHint();
  } catch (e) {
    console.error("ExiftoolCheck error:", e);
  }
}

document.getElementById("exiftool-doc-link")?.addEventListener("click", (e) => {
  e.preventDefault();
  Browser.OpenURL("https://exiftool.org/");
});

document.getElementById("btn-exiftool-recheck")?.addEventListener("click", () => {
  refreshExiftoolStatus();
});

function openAbout() {
  aboutModal.classList.add("open");
  aboutModal.setAttribute("aria-hidden", "false");
  aboutLastFocused = document.activeElement;
  if (aboutModalCloseBtn) {
    aboutModalCloseBtn.focus();
  }
}

function closeAbout() {
  aboutModal.classList.remove("open");
  aboutModal.setAttribute("aria-hidden", "true");
  if (aboutLastFocused instanceof HTMLElement) {
    aboutLastFocused.focus();
  } else {
    aboutBtn?.focus();
  }
  aboutLastFocused = null;
}

function getAboutFocusableElements() {
  if (!aboutModal) return [];
  return Array.from(
    aboutModal.querySelectorAll(
      'button:not([disabled]), a[href], input:not([disabled]), select:not([disabled]), textarea:not([disabled]), [tabindex]:not([tabindex="-1"])',
    ),
  ).filter((el) => el instanceof HTMLElement && el.offsetParent !== null);
}

function trapAboutFocus(event) {
  if (event.key !== "Tab" || !aboutModal.classList.contains("open")) return;
  const focusable = getAboutFocusableElements();
  if (focusable.length === 0) return;
  const first = focusable[0];
  const last = focusable[focusable.length - 1];
  const active = document.activeElement;

  if (event.shiftKey && active === first) {
    event.preventDefault();
    last.focus();
  } else if (!event.shiftKey && active === last) {
    event.preventDefault();
    first.focus();
  }
}

aboutBtn.addEventListener("click", openAbout);
aboutModalCloseBtn.addEventListener("click", closeAbout);
aboutModal.querySelectorAll("[data-close-modal]").forEach((el) => {
  el.addEventListener("click", closeAbout);
});
document.addEventListener("keydown", (e) => {
  if (e.key === "Escape" && aboutModal.classList.contains("open")) {
    closeAbout();
  }
  trapAboutFocus(e);
});

// WebView often blocks default link navigation; open via OS browser / mail client.
aboutModal.addEventListener("click", (e) => {
  const a = e.target.closest("a");
  if (!a || !aboutModal.contains(a)) return;
  const href = a.getAttribute("href");
  if (!href || href.startsWith("#")) return;
  e.preventDefault();
  e.stopPropagation();
  Browser.OpenURL(href);
});

document.getElementById("btn-select").addEventListener("click", async () => {
  try {
    hideScanErrorBanner();
    const path = await MetadataService.SelectFolder();
    if (!path) return;
    currentPath = path;

    showView("scan");
    views.scan?.setAttribute("aria-busy", "true");
    btnFix.disabled = true;
    btnFix.textContent = "Fix files";
    if (scanEmptyMessage) {
      scanEmptyMessage.classList.add("hidden");
    }
    const scanPathEl = document.getElementById("scan-path");
    if (scanPathEl) {
      scanPathEl.textContent = basename(path);
      scanPathEl.setAttribute("title", path);
    }
    document.getElementById("stat-total").textContent = "…";
    document.getElementById("stat-matched").textContent = "…";

    scanData = await MetadataService.ScanFolder(path);
    renderScanResults(scanData);
    await updateCheckpointHint();
  } catch (err) {
    console.error("SelectFolder/Scan error:", err);
    scanData = null;
    currentPath = "";
    showScanErrorBanner("Could not open or scan that folder. " + formatErrorMessage(err));
    showView("welcome");
  } finally {
    views.scan?.setAttribute("aria-busy", "false");
  }
});

btnBack.addEventListener("click", () => {
  showView("welcome");
});

btnFixPause?.addEventListener("click", () => {
  MetadataService.FixPause()
    .then(() => setProcessingPaused(true))
    .catch((e) => console.error("FixPause error:", e));
});

btnFixResume?.addEventListener("click", () => {
  MetadataService.FixResume()
    .then(() => setProcessingPaused(false))
    .catch((e) => console.error("FixResume error:", e));
});

btnFixStop?.addEventListener("click", () => {
  MetadataService.FixAbort().catch((e) => console.error("FixAbort error:", e));
});

btnFix.addEventListener("click", async () => {
  if (!currentPath || !exiftoolOk) return;
  hideScanErrorBanner();
  showView("processing");
  resetProcessingControls();

  setProgressUi(0, 0);

  const offComplete = Events.Once("fix-complete", (event) => {
    const payload = event.data;
    if (payload.error) {
      console.error("Fix run failed:", payload.error);
      showScanErrorBanner("The run did not finish. " + payload.error);
      resetProcessingControls();
      showView("scan");
      void (async () => {
        await updateCheckpointHint();
        const path = currentPath;
        if (!path) return;
        try {
          scanData = await MetadataService.ScanFolder(path);
          renderScanResults(scanData);
        } catch (e) {
          console.error("Rescan after fix error:", e);
        }
      })();
      return;
    }
    if (payload.result) {
      renderDoneResults(payload.result);
      showView("done");
    }
  });

  try {
    const deleteJson = document.getElementById("chk-delete-json").checked;
    let resume = false;
    try {
      resume = await MetadataService.FixCheckpointAvailable(currentPath);
    } catch (e) {
      console.error("FixCheckpointAvailable error:", e);
    }
    await MetadataService.FixMetadata(currentPath, deleteJson, resume);
  } catch (err) {
    offComplete();
    console.error("FixMetadata error:", err);
    showScanErrorBanner("Could not start the run. " + formatErrorMessage(err));
    resetProcessingControls();
    showView("scan");
    await updateCheckpointHint();
    try {
      scanData = await MetadataService.ScanFolder(currentPath);
      renderScanResults(scanData);
    } catch (e) {
      console.error("Rescan after FixMetadata error:", e);
    }
  }
});

document.getElementById("btn-restart").addEventListener("click", () => {
  currentPath = "";
  scanData = null;
  const note = document.getElementById("done-resume-note");
  if (note) {
    note.textContent = "";
    note.classList.add("hidden");
  }
  showView("welcome");
});

Events.On("fix-progress", (event) => {
  const p = event.data;
  setProgressUi(p.current, p.total, p.file);
});

function renderScanResults(data) {
  const ready = Number(data.withJson) || 0;
  const total = Number(data.totalMedia) || 0;
  const noData = Number(data.withoutJson) || 0;
  const orphans = Number(data.orphanJson) || 0;

  document.getElementById("stat-matched").textContent = num(ready);
  document.getElementById("stat-total").textContent = num(total);
  document.getElementById("chip-ready").textContent = `${num(ready)} ready to fix`;
  document.getElementById("chip-nodata").textContent = `${num(noData)} ${plural(noData, "has", "have")} no info to restore`;
  document.getElementById("chip-orphan").textContent =
    `${num(orphans)} extra .json ${plural(orphans, "file", "files")}, safe to ignore`;

  const scanPathEl = document.getElementById("scan-path");
  if (scanPathEl && data.folderPath) {
    scanPathEl.textContent = basename(data.folderPath);
    scanPathEl.setAttribute("title", data.folderPath);
  }

  if (scanEmptyMessage) {
    scanEmptyMessage.classList.add("hidden");
    scanEmptyMessage.textContent = "";
  }

  if (!data.files || data.files.length === 0) {
    if (scanEmptyMessage) {
      scanEmptyMessage.textContent = "No photos or videos this app can read in that folder.";
      scanEmptyMessage.classList.remove("hidden");
    }
    btnFix.disabled = true;
    btnFix.textContent = "Fix files";
    return;
  }

  btnFix.disabled = ready === 0 || !exiftoolOk;
  btnFix.textContent = ready > 0 ? `Fix ${num(ready)} ${plural(ready, "file", "files")}` : "Nothing to fix here";
}

function renderDoneResults(result) {
  const success = Number(result.success) || 0;
  const doneTitle = document.getElementById("done-title");
  if (doneTitle) {
    if (result.aborted) {
      doneTitle.textContent = `Stopped after ${num(success)} ${plural(success, "photo", "photos")}.`;
    } else if (success > 0) {
      doneTitle.textContent = `${num(success)} ${plural(success, "photo", "photos")} got ${plural(success, "its", "their")} ${plural(success, "date", "dates")} back.`;
    } else {
      doneTitle.textContent = "Nothing needed changing.";
    }
  }

  const doneSummary = document.getElementById("done-summary");
  if (doneSummary) {
    doneSummary.textContent = result.aborted
      ? "Stopped early. Your progress is saved: open the same folder and start again to carry on."
      : "This run.";
  }

  const resumeNote = document.getElementById("done-resume-note");
  if (resumeNote) {
    const notes = [];
    if (result.resumed) notes.push("Carried on from last time.");
    if (result.aborted) notes.push("Open this folder again to finish the rest.");
    if (notes.length > 0) {
      resumeNote.textContent = notes.join(" ");
      resumeNote.classList.remove("hidden");
    } else {
      resumeNote.textContent = "";
      resumeNote.classList.add("hidden");
    }
  }

  document.getElementById("result-success").textContent = num(success);
  document.getElementById("result-skipped").textContent = num(result.skipped);
  document.getElementById("result-failed").textContent = num(result.failed);

  const extra = document.getElementById("result-json-delete");
  const parts = [];
  if (result.jsonDeleted > 0) {
    parts.push(
      `Deleted ${num(result.jsonDeleted)} .json ${plural(result.jsonDeleted, "file", "files")}. Their info is now saved inside your photos.`,
    );
  }
  if (result.jsonDeleteFailed > 0) {
    parts.push(
      `Could not delete ${num(result.jsonDeleteFailed)} .json ${plural(result.jsonDeleteFailed, "file", "files")}.`,
    );
  }
  if (parts.length > 0) {
    extra.textContent = parts.join(" ");
    extra.classList.remove("hidden");
  } else {
    extra.textContent = "";
    extra.classList.add("hidden");
  }
}

refreshExiftoolStatus();
showView("welcome");

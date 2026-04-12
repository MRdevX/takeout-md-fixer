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

const btnBack = document.getElementById("btn-back");

const processingTitle = document.getElementById("processing-title");
const processingSubtitle = document.getElementById("processing-subtitle");
const progressBarEl = document.getElementById("progress-bar");
const progressTextEl = document.getElementById("progress-text");
const progressFileNameEl = document.getElementById("progress-file-name");
const progressFilePathEl = document.getElementById("progress-file-path");
const btnFixPause = document.getElementById("btn-fix-pause");
const btnFixResume = document.getElementById("btn-fix-resume");
const btnFixStop = document.getElementById("btn-fix-stop");
const checkpointHint = document.getElementById("checkpoint-hint");
const scanEmptyMessage = document.getElementById("scan-empty-message");

const PROCESSING_SUBTITLE_DEFAULT =
  "Copying date and location into your files. This may take a few minutes. You can pause or stop; progress is saved.";
const PROCESSING_SUBTITLE_PAUSED = "Paused. Click Resume to continue, or Stop to finish after this file. Progress is saved.";

const stepperSteps = document.querySelectorAll("#app-stepper .stepper-step");
const stepperRails = document.querySelectorAll("#app-stepper .stepper-rail");

/** @param {string} filePath */
function basename(filePath) {
  if (!filePath) return "";
  const normalized = filePath.replace(/\\/g, "/");
  const parts = normalized.split("/");
  const last = parts.pop();
  return last || filePath;
}

/**
 * @param {"welcome" | "scan" | "processing" | "done"} name
 */
function updateStepper(name) {
  if (!stepperSteps.length) return;
  stepperSteps.forEach((el) => {
    el.classList.remove("stepper-step--current", "stepper-step--complete");
    el.removeAttribute("aria-current");
  });
  stepperRails.forEach((el) => el.classList.remove("stepper-rail--complete"));

  if (name === "welcome") {
    if (stepperSteps[0]) {
      stepperSteps[0].classList.add("stepper-step--current");
      stepperSteps[0].setAttribute("aria-current", "step");
    }
  } else if (name === "scan") {
    if (stepperSteps[0]) stepperSteps[0].classList.add("stepper-step--complete");
    if (stepperRails[0]) stepperRails[0].classList.add("stepper-rail--complete");
    if (stepperSteps[1]) {
      stepperSteps[1].classList.add("stepper-step--current");
      stepperSteps[1].setAttribute("aria-current", "step");
    }
  } else if (name === "processing") {
    if (stepperSteps[0]) stepperSteps[0].classList.add("stepper-step--complete");
    if (stepperSteps[1]) stepperSteps[1].classList.add("stepper-step--complete");
    if (stepperRails[0]) stepperRails[0].classList.add("stepper-rail--complete");
    if (stepperRails[1]) stepperRails[1].classList.add("stepper-rail--complete");
    if (stepperSteps[2]) {
      stepperSteps[2].classList.add("stepper-step--current");
      stepperSteps[2].setAttribute("aria-current", "step");
    }
  } else if (name === "done") {
    stepperSteps.forEach((el) => el.classList.add("stepper-step--complete"));
    stepperRails.forEach((el) => el.classList.add("stepper-rail--complete"));
  }
}

function showView(name) {
  Object.values(views).forEach((v) => v.classList.remove("active"));
  views[name].classList.add("active");
  updateStepper(name);
  const onScan = name === "scan";
  btnBack.hidden = !onScan;
  btnBack.setAttribute("aria-hidden", onScan ? "false" : "true");
}

function setProcessingPaused(paused) {
  if (!processingTitle || !btnFixPause || !btnFixResume) return;
  if (paused) {
    processingTitle.textContent = "Paused";
    if (processingSubtitle) processingSubtitle.textContent = PROCESSING_SUBTITLE_PAUSED;
    btnFixPause.classList.add("hidden");
    btnFixResume.classList.remove("hidden");
  } else {
    processingTitle.textContent = "Processing…";
    if (processingSubtitle) processingSubtitle.textContent = PROCESSING_SUBTITLE_DEFAULT;
    btnFixPause.classList.remove("hidden");
    btnFixResume.classList.add("hidden");
  }
}

function resetProcessingControls() {
  setProcessingPaused(false);
}

function setProgressUi(current, total, file = "") {
  const safeTotal = Number(total) > 0 ? Number(total) : 0;
  const safeCurrent = Math.max(0, Number(current) || 0);
  const boundedCurrent = safeTotal > 0 ? Math.min(safeCurrent, safeTotal) : safeCurrent;
  const pct = safeTotal > 0 ? Math.round((boundedCurrent / safeTotal) * 100) : 0;

  if (progressBarEl) {
    progressBarEl.style.width = `${pct}%`;
    progressBarEl.setAttribute("aria-valuenow", String(pct));
    const valText = safeTotal > 0 ? `${boundedCurrent} of ${safeTotal} files, ${pct} percent` : `${pct} percent`;
    progressBarEl.setAttribute("aria-valuetext", valText);
  }
  if (progressTextEl) {
    progressTextEl.textContent =
      safeTotal > 0 ? `${boundedCurrent} of ${safeTotal} (${pct}%)` : `${boundedCurrent} of ${safeTotal} (${pct}%)`;
  }
  const base = basename(file);
  if (progressFileNameEl && progressFilePathEl) {
    if (!file) {
      progressFileNameEl.textContent = "";
      progressFilePathEl.textContent = "";
    } else {
      progressFileNameEl.textContent = base;
      progressFilePathEl.textContent = base === file ? "" : file;
    }
  }
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
      checkpointHint.textContent =
        "Continuing from your last session — files already updated will be skipped automatically.";
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
        msgEl.textContent = st.message || "ExifTool not found.";
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
    document.getElementById("btn-fix").disabled = true;
    if (scanEmptyMessage) {
      scanEmptyMessage.classList.add("hidden");
    }
    const scanPathEl = document.getElementById("scan-path");
    if (scanPathEl) {
      scanPathEl.textContent = `Folder: ${basename(path)}`;
      scanPathEl.setAttribute("title", path);
    }
    document.getElementById("stat-total").textContent = "...";
    document.getElementById("stat-matched").textContent = "...";
    document.getElementById("stat-unmatched").textContent = "...";
    document.getElementById("stat-orphan-json").textContent = "...";

    scanData = await MetadataService.ScanFolder(path);
    renderScanResults(scanData);
    await updateCheckpointHint();
  } catch (err) {
    console.error("SelectFolder/Scan error:", err);
    scanData = null;
    currentPath = "";
    showScanErrorBanner("We couldn’t open or scan that folder. " + formatErrorMessage(err));
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

document.getElementById("btn-fix").addEventListener("click", async () => {
  if (!currentPath || !exiftoolOk) return;
  hideScanErrorBanner();
  showView("processing");
  resetProcessingControls();

  setProgressUi(0, 0);

  const offComplete = Events.Once("fix-complete", (event) => {
    const payload = event.data;
    if (payload.error) {
      console.error("Fix run failed:", payload.error);
      showScanErrorBanner("The update didn’t finish. " + payload.error);
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
    showScanErrorBanner("We couldn’t start the update. " + formatErrorMessage(err));
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
  document.getElementById("stat-total").textContent = data.totalMedia;
  document.getElementById("stat-matched").textContent = data.withJson;
  document.getElementById("stat-unmatched").textContent = data.withoutJson;
  document.getElementById("stat-orphan-json").textContent = data.orphanJson ?? 0;

  const scanPathEl = document.getElementById("scan-path");
  if (scanPathEl && data.folderPath) {
    scanPathEl.textContent = `Folder: ${basename(data.folderPath)}`;
    scanPathEl.setAttribute("title", data.folderPath);
  }

  if (scanEmptyMessage) {
    scanEmptyMessage.classList.add("hidden");
    scanEmptyMessage.textContent = "";
  }

  if (!data.files || data.files.length === 0) {
    if (scanEmptyMessage) {
      scanEmptyMessage.textContent = "No supported photos or videos found in this folder.";
      scanEmptyMessage.classList.remove("hidden");
    }
    document.getElementById("btn-fix").disabled = true;
    return;
  }

  document.getElementById("btn-fix").disabled = data.withJson === 0 || !exiftoolOk;
}

function renderDoneResults(result) {
  const doneTitle = document.getElementById("done-title");
  if (doneTitle) {
    doneTitle.textContent = result.aborted ? "Stopped" : "Done";
  }

  const doneSummary = document.getElementById("done-summary");
  if (doneSummary) {
    doneSummary.textContent = result.aborted
      ? "Stopped before all files were processed. Progress is saved. Click Start over, select the same folder, and press Start to finish."
      : "Results for this run.";
  }

  const resumeNote = document.getElementById("done-resume-note");
  if (resumeNote) {
    if (result.resumed) {
      resumeNote.textContent = "Continued from a previous run.";
      resumeNote.classList.remove("hidden");
    } else {
      resumeNote.textContent = "";
      resumeNote.classList.add("hidden");
    }
  }

  document.getElementById("result-success").textContent = result.success;
  document.getElementById("result-skipped").textContent = result.skipped;
  document.getElementById("result-failed").textContent = result.failed;

  const extra = document.getElementById("result-json-delete");
  const parts = [];
  if (result.jsonDeleted > 0) {
    parts.push(
      `Deleted ${result.jsonDeleted} extra info file${result.jsonDeleted === 1 ? "" : "s"} (dates and locations stay in your photos and videos).`,
    );
  }
  if (result.jsonDeleteFailed > 0) {
    parts.push(`Could not delete ${result.jsonDeleteFailed} extra info file${result.jsonDeleteFailed === 1 ? "" : "s"}.`);
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
updateStepper("welcome");

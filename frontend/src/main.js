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
const progressBarEl = document.getElementById("progress-bar");
const progressTextEl = document.getElementById("progress-text");
const progressFileEl = document.getElementById("progress-file");
const btnFixPause = document.getElementById("btn-fix-pause");
const btnFixResume = document.getElementById("btn-fix-resume");
const btnFixStop = document.getElementById("btn-fix-stop");
const chkResumeFix = document.getElementById("chk-resume-fix");
const checkpointHint = document.getElementById("checkpoint-hint");

function showView(name) {
    Object.values(views).forEach((v) => v.classList.remove("active"));
    views[name].classList.add("active");
    const onScan = name === "scan";
    btnBack.hidden = !onScan;
    btnBack.setAttribute("aria-hidden", onScan ? "false" : "true");
}

function setProcessingPaused(paused) {
    if (!processingTitle || !btnFixPause || !btnFixResume) return;
    if (paused) {
        processingTitle.textContent = "Paused";
        btnFixPause.classList.add("hidden");
        btnFixResume.classList.remove("hidden");
    } else {
        processingTitle.textContent = "Working…";
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
    }
    if (progressTextEl) {
        progressTextEl.textContent = `${boundedCurrent} / ${safeTotal} (${pct}%)`;
    }
    if (progressFileEl) {
        progressFileEl.textContent = file;
    }
}

function renderFileListMessage(message, loading = false) {
    const tbody = document.getElementById("file-list-body");
    if (!tbody) return;
    const rowClass = loading ? "file-list-message file-list-message--loading" : "file-list-message";
    const messageClass = loading ? "file-list-message-text file-list-message-text--loading" : "file-list-message-text";
    const spinner = loading ? '<span class="inline-spinner" aria-hidden="true"></span>' : "";
    tbody.innerHTML = `<tr><td colspan="3" class="${rowClass}"><span class="${messageClass}">${spinner}${message}</span></td></tr>`;
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
    if (!checkpointHint || !chkResumeFix) return;
    if (!currentPath) {
        checkpointHint.classList.add("hidden");
        chkResumeFix.checked = false;
        chkResumeFix.disabled = true;
        return;
    }
    try {
        const ok = await MetadataService.FixCheckpointAvailable(currentPath);
        if (ok) {
            checkpointHint.textContent =
                "A previous run saved progress in this folder. Check “Continue where you left off” below to skip files already fixed, or delete the hidden file .takeout-md-fixer-checkpoint.json in this folder to start from scratch.";
            checkpointHint.classList.remove("hidden");
            chkResumeFix.disabled = false;
        } else {
            checkpointHint.classList.add("hidden");
            chkResumeFix.checked = false;
            chkResumeFix.disabled = true;
        }
    } catch (e) {
        console.error("FixCheckpointAvailable error:", e);
        checkpointHint.classList.add("hidden");
        chkResumeFix.checked = false;
        chkResumeFix.disabled = true;
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
                msgEl.textContent = st.message || "ExifTool was not found.";
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
            'button:not([disabled]), a[href], input:not([disabled]), select:not([disabled]), textarea:not([disabled]), [tabindex]:not([tabindex="-1"])'
        )
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
        renderFileListMessage("Scanning...", true);
        document.getElementById("scan-path").textContent = path;
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
        showScanErrorBanner(
            "Could not open or scan that folder. " + formatErrorMessage(err)
        );
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
            showScanErrorBanner("Fix did not finish. " + payload.error);
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
        const resume = chkResumeFix?.checked === true;
        await MetadataService.FixMetadata(currentPath, deleteJson, resume);
    } catch (err) {
        offComplete();
        console.error("FixMetadata error:", err);
        showScanErrorBanner("Could not start the fix. " + formatErrorMessage(err));
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

    const tbody = document.getElementById("file-list-body");
    if (!data.files || data.files.length === 0) {
        renderFileListMessage("No media files found");
        document.getElementById("btn-fix").disabled = true;
        return;
    }

    tbody.innerHTML = data.files
        .map(
            (f) => `<tr>
            <td title="${escapeHtml(f.path)}">${escapeHtml(f.name)}</td>
            <td>${f.hasJson ? '<span class="badge badge-yes">Yes</span>' : '<span class="badge badge-no">No</span>'}</td>
            <td><span class="badge badge-${f.status}">${f.status}</span></td>
        </tr>`
        )
        .join("");

    document.getElementById("btn-fix").disabled = data.withJson === 0 || !exiftoolOk;
}

function renderDoneResults(result) {
    const doneTitle = document.getElementById("done-title");
    if (doneTitle) {
        doneTitle.textContent = result.aborted ? "Stopped early" : "Done";
    }

    const resumeNote = document.getElementById("done-resume-note");
    if (resumeNote) {
        if (result.resumed) {
            resumeNote.textContent = "Continued from a previous interrupted run.";
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
        parts.push(`Sidecars removed: ${result.jsonDeleted}`);
    }
    if (result.jsonDeleteFailed > 0) {
        parts.push(`Could not remove: ${result.jsonDeleteFailed}`);
    }
    if (parts.length > 0) {
        extra.textContent = parts.join(" · ");
        extra.classList.remove("hidden");
    } else {
        extra.textContent = "";
        extra.classList.add("hidden");
    }
}

function escapeHtml(str) {
    const div = document.createElement("div");
    div.textContent = str;
    return div.innerHTML;
}

refreshExiftoolStatus();

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

function showView(name) {
    Object.values(views).forEach((v) => v.classList.remove("active"));
    views[name].classList.add("active");
    const onScan = name === "scan";
    btnBack.hidden = !onScan;
    btnBack.setAttribute("aria-hidden", onScan ? "false" : "true");
}

let currentPath = "";
let scanData = null;
/** Whether ExifTool was found (PATH plus common install locations). */
let exiftoolOk = true;

const aboutModal = document.getElementById("about-modal");
const exiftoolWarningEl = document.getElementById("exiftool-warning");

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
}

function closeAbout() {
    aboutModal.classList.remove("open");
    aboutModal.setAttribute("aria-hidden", "true");
}

document.getElementById("btn-about").addEventListener("click", openAbout);
document.getElementById("about-modal-close").addEventListener("click", closeAbout);
aboutModal.querySelectorAll("[data-close-modal]").forEach((el) => {
    el.addEventListener("click", closeAbout);
});
document.addEventListener("keydown", (e) => {
    if (e.key === "Escape" && aboutModal.classList.contains("open")) {
        closeAbout();
    }
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
        const path = await MetadataService.SelectFolder();
        if (!path) return;
        currentPath = path;

        showView("scan");
        document.getElementById("btn-fix").disabled = true;
        document.getElementById("file-list-body").innerHTML =
            '<tr><td colspan="3" style="text-align:center;padding:24px;color:var(--text-muted)">Scanning...</td></tr>';
        document.getElementById("scan-path").textContent = path;
        document.getElementById("stat-total").textContent = "...";
        document.getElementById("stat-matched").textContent = "...";
        document.getElementById("stat-unmatched").textContent = "...";

        scanData = await MetadataService.ScanFolder(path);
        renderScanResults(scanData);
    } catch (err) {
        console.error("SelectFolder/Scan error:", err);
    }
});

btnBack.addEventListener("click", () => {
    showView("welcome");
});

document.getElementById("btn-fix").addEventListener("click", async () => {
    if (!currentPath || !exiftoolOk) return;
    showView("processing");

    document.getElementById("progress-bar").style.width = "0%";
    document.getElementById("progress-text").textContent = "0 / 0";
    document.getElementById("progress-file").textContent = "";

    try {
        const deleteJson = document.getElementById("chk-delete-json").checked;
        const result = await MetadataService.FixMetadata(currentPath, deleteJson);
        renderDoneResults(result);
        showView("done");
    } catch (err) {
        console.error("FixMetadata error:", err);
        showView("scan");
    }
});

document.getElementById("btn-restart").addEventListener("click", () => {
    currentPath = "";
    scanData = null;
    showView("welcome");
});

Events.On("fix-progress", (event) => {
    const p = event.data;
    const pct = Math.round((p.current / p.total) * 100);
    document.getElementById("progress-bar").style.width = pct + "%";
    document.getElementById("progress-text").textContent = `${p.current} / ${p.total}`;
    document.getElementById("progress-file").textContent = p.file;
});

function renderScanResults(data) {
    document.getElementById("stat-total").textContent = data.totalMedia;
    document.getElementById("stat-matched").textContent = data.withJson;
    document.getElementById("stat-unmatched").textContent = data.withoutJson;

    const tbody = document.getElementById("file-list-body");
    if (!data.files || data.files.length === 0) {
        tbody.innerHTML =
            '<tr><td colspan="3" style="text-align:center;padding:24px;color:var(--text-muted)">No media files found</td></tr>';
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

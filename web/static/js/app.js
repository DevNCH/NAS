(function () {
    "use strict";

    var viewLogin = document.getElementById("view-login");
    var viewDashboard = document.getElementById("view-dashboard");

    var loginForm = document.getElementById("login-form");
    var loginError = document.getElementById("login-error");

    var userName = document.getElementById("user-name");
    var userRole = document.getElementById("user-role");
    var btnLogout = document.getElementById("btn-logout");

    var fileInput = document.getElementById("file-input");
    var filePickerLabel = document.getElementById("file-picker-label");
    var btnUpload = document.getElementById("btn-upload");
    var uploadMsg = document.getElementById("upload-msg");

    var btnRefresh = document.getElementById("btn-refresh");
    var filesTable = document.getElementById("files-table");
    var filesBody = document.getElementById("files-body");
    var filesEmpty = document.getElementById("files-empty");

    var adminPanel = document.getElementById("admin-panel");
    var btnAdminPing = document.getElementById("btn-admin-ping");
    var adminPingResult = document.getElementById("admin-ping-result");

    var statusDot = document.getElementById("status-dot");
    var statusText = document.getElementById("status-text");

    function setStatus(online) {
        statusDot.classList.remove("online", "offline");
        statusDot.classList.add(online ? "online" : "offline");
        statusText.textContent = online
            ? "servidor online"
            : "não foi possível falar com o servidor";
    }

    function showMsg(el, text, kind) {
        el.textContent = text;
        el.className = "msg " + (kind === "error" ? "msg-error" : "msg-success");
        el.hidden = false;
    }

    function hideMsg(el) {
        el.hidden = true;
    }

    function formatSize(bytes) {
        if (bytes < 1024) return bytes + " B";
        if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + " KB";
        return (bytes / (1024 * 1024)).toFixed(1) + " MB";
    }

    function formatDate(iso) {
        var d = new Date(iso);
        if (isNaN(d.getTime())) return iso;
        return d.toLocaleString("pt-BR", {
            day: "2-digit",
            month: "2-digit",
            year: "numeric",
            hour: "2-digit",
            minute: "2-digit",
        });
    }

    function showLogin() {
        viewLogin.hidden = false;
        viewDashboard.hidden = true;
    }

    function showDashboard(user) {
        viewLogin.hidden = true;
        viewDashboard.hidden = false;

        userName.textContent = user.username;
        userRole.textContent = user.role;

        adminPanel.hidden = user.role !== "admin";

        loadFiles();
    }

    // ---- sessão ----

    function checkSession() {
        fetch("/me", { credentials: "same-origin" })
            .then(function (res) {
                setStatus(true);
                if (res.status === 200) {
                    return res.json().then(showDashboard);
                }
                showLogin();
            })
            .catch(function () {
                setStatus(false);
                showLogin();
            });
    }

    loginForm.addEventListener("submit", function (evt) {
        evt.preventDefault();
        hideMsg(loginError);

        var username = document.getElementById("login-username").value.trim();
        var password = document.getElementById("login-password").value;

        fetch("/login", {
            method: "POST",
            credentials: "same-origin",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ username: username, password: password }),
        })
            .then(function (res) {
                setStatus(true);
                if (!res.ok) {
                    return res.json().then(function (data) {
                        throw new Error(data.error || "não foi possível entrar");
                    });
                }
                return res.json();
            })
            .then(function (user) {
                loginForm.reset();
                showDashboard(user);
            })
            .catch(function (err) {
                setStatus(false);
                showMsg(loginError, err.message, "error");
            });
    });

    btnLogout.addEventListener("click", function () {
        fetch("/logout", { method: "POST", credentials: "same-origin" }).finally(
            function () {
                showLogin();
            }
        );
    });

    // ---- upload ----

    fileInput.addEventListener("change", function () {
        var file = fileInput.files[0];
        filePickerLabel.textContent = file ? file.name : "Escolher arquivo";
        btnUpload.disabled = !file;
        hideMsg(uploadMsg);
    });

    btnUpload.addEventListener("click", function () {
        var file = fileInput.files[0];
        if (!file) return;

        btnUpload.disabled = true;
        btnUpload.textContent = "Enviando…";
        hideMsg(uploadMsg);

        var formData = new FormData();
        formData.append("file", file);

        fetch("/upload", {
            method: "POST",
            credentials: "same-origin",
            body: formData,
        })
            .then(function (res) {
                if (!res.ok) {
                    return res.json().then(function (data) {
                        throw new Error(data.error || "falha no envio");
                    });
                }
                return res.json();
            })
            .then(function () {
                showMsg(uploadMsg, "Arquivo enviado com sucesso.", "success");
                fileInput.value = "";
                filePickerLabel.textContent = "Escolher arquivo";
                loadFiles();
            })
            .catch(function (err) {
                showMsg(uploadMsg, err.message, "error");
            })
            .finally(function () {
                btnUpload.disabled = true;
                btnUpload.textContent = "Enviar";
            });
    });

    // ---- lista de arquivos ----

    function loadFiles() {
        fetch("/files", { credentials: "same-origin" })
            .then(function (res) {
                if (!res.ok) throw new Error("falha ao carregar arquivos");
                return res.json();
            })
            .then(renderFiles)
            .catch(function () {
                filesTable.hidden = true;
                filesEmpty.hidden = false;
                filesEmpty.textContent = "Não foi possível carregar os arquivos.";
            });
    }

    function renderFiles(files) {
        filesBody.innerHTML = "";

        if (!files || files.length === 0) {
            filesTable.hidden = true;
            filesEmpty.hidden = false;
            filesEmpty.textContent = "Nenhum arquivo enviado ainda.";
            return;
        }

        filesEmpty.hidden = true;
        filesTable.hidden = false;

        files.forEach(function (file) {
            var tr = document.createElement("tr");

            var tdName = document.createElement("td");
            var nameSpan = document.createElement("span");
            nameSpan.className = "file-name";
            nameSpan.textContent = file.filename;
            tdName.appendChild(nameSpan);

            var tdSize = document.createElement("td");
            tdSize.textContent = formatSize(file.size);

            var tdDate = document.createElement("td");
            tdDate.textContent = formatDate(file.created_at);

            var tdAction = document.createElement("td");
            tdAction.className = "col-action";
            var link = document.createElement("a");
            link.href = "/download/" + file.id;
            link.className = "btn btn-ghost btn-small";
            link.textContent = "Baixar";
            tdAction.appendChild(link);

            tr.appendChild(tdName);
            tr.appendChild(tdSize);
            tr.appendChild(tdDate);
            tr.appendChild(tdAction);

            filesBody.appendChild(tr);
        });
    }

    btnRefresh.addEventListener("click", loadFiles);

    // ---- painel admin ----

    btnAdminPing.addEventListener("click", function () {
        adminPingResult.textContent = "testando…";

        fetch("/admin/ping", { credentials: "same-origin" })
            .then(function (res) {
                if (!res.ok) throw new Error("acesso negado (" + res.status + ")");
                return res.json();
            })
            .then(function (data) {
                adminPingResult.textContent = "✅ " + data.message;
            })
            .catch(function (err) {
                adminPingResult.textContent = "❌ " + err.message;
            });
    });

    checkSession();
})();

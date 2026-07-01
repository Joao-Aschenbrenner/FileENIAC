// SPDX-License-Identifier: MIT
use serde::Serialize;
use std::sync::{Arc, Condvar, Mutex};
use tauri::Manager;
use tauri_plugin_shell::process::CommandChild;
use tauri_plugin_shell::ShellExt;

struct BackendState {
    port: String,
    token: String,
    child: Option<CommandChild>,
    ready: bool,
}

struct AppState {
    backend: Mutex<BackendState>,
    port_ready: Arc<Condvar>,
    log_path: String,
    bootstrap_log_path: String,
    startup_error: Mutex<Option<StartupError>>,
}

#[derive(Serialize, Clone)]
struct StartupError {
    message: String,
    exit_code: Option<i32>,
}

#[derive(Serialize)]
struct BackendInfo {
    base_url: String,
    token: String,
    ready: bool,
}

#[derive(Serialize)]
struct DiagnosticsInfo {
    log_path: String,
    bootstrap_log_path: String,
    startup_error: Option<StartupError>,
}

fn generate_token() -> String {
    use std::time::{SystemTime, UNIX_EPOCH};
    let t = SystemTime::now()
        .duration_since(UNIX_EPOCH)
        .unwrap_or_default()
        .as_nanos();
    format!("{:032x}", t)
}

fn strip_token(line: &str) -> String {
    if let Some(idx) = line.find("token=") {
        let before = &line[..idx + 6];
        let after = &line[idx + 6..];
        if let Some(space) = after.find(' ') {
            format!("{}***{}", before, &after[space..])
        } else {
            format!("{}***", before)
        }
    } else {
        line.to_string()
    }
}

fn write_log(log_path: &str, line: &str) {
    use std::io::Write;
    if let Ok(mut f) = std::fs::OpenOptions::new()
        .create(true)
        .append(true)
        .open(log_path)
    {
        let _ = writeln!(f, "{}", line);
    }
}

#[tauri::command]
fn get_backend_info(state: tauri::State<AppState>) -> BackendInfo {
    let cvar = state.port_ready.clone();
    let mut b = state.backend.lock().unwrap();
    let deadline = std::time::Instant::now() + std::time::Duration::from_secs(20);
    while b.port.is_empty() && !b.ready {
        let now = std::time::Instant::now();
        if now >= deadline {
            write_log(
                &state.bootstrap_log_path,
                "[bootstrap] get_backend_info_timeout waiting_for_port=true",
            );
            let mut err = state.startup_error.lock().unwrap();
            if err.is_none() {
                *err = Some(StartupError {
                    message: "Tempo esgotado ao preparar o ambiente local.".to_string(),
                    exit_code: None,
                });
            }
            break;
        }
        let remaining = deadline.saturating_duration_since(now);
        let guard = cvar.wait_timeout(b, remaining).unwrap();
        b = guard.0;
    }
    BackendInfo {
        base_url: if b.port.is_empty() {
            String::new()
        } else {
            format!("http://127.0.0.1:{}/api", b.port)
        },
        token: b.token.clone(),
        ready: b.ready && !b.port.is_empty(),
    }
}

#[tauri::command]
fn get_diagnostics(state: tauri::State<AppState>) -> DiagnosticsInfo {
    let err = state.startup_error.lock().unwrap().clone();
    DiagnosticsInfo {
        log_path: state.log_path.clone(),
        bootstrap_log_path: state.bootstrap_log_path.clone(),
        startup_error: err,
    }
}

#[tauri::command]
fn get_sidecar_log_path(state: tauri::State<AppState>) -> String {
    state.log_path.clone()
}

#[cfg_attr(mobile, tauri::mobile_entry_point)]
pub fn run() {
    tauri::Builder::default()
        .plugin(tauri_plugin_opener::init())
        .plugin(tauri_plugin_dialog::init())
        .plugin(tauri_plugin_shell::init())
        .setup(|app| {
            let token = generate_token();
            let log_dir = app
                .path()
                .app_log_dir()
                .unwrap_or_else(|_| std::path::PathBuf::from("."));
            let _ = std::fs::create_dir_all(&log_dir);
            let log_path = log_dir
                .join("backend.log")
                .to_string_lossy()
                .to_string();
            let bootstrap_log_path = log_dir
                .join("fileeniac-bootstrap.log")
                .to_string_lossy()
                .to_string();

            let _ = std::fs::remove_file(&bootstrap_log_path);
            write_log(&bootstrap_log_path, "[bootstrap] starting FileENIAC desktop");
            write_log(&bootstrap_log_path, &format!("[bootstrap] app_version=0.1.9"));

            let port_ready = Arc::new(Condvar::new());

            let state = AppState {
                backend: Mutex::new(BackendState {
                    port: String::new(),
                    token: token.clone(),
                    child: None,
                    ready: false,
                }),
                port_ready: port_ready.clone(),
                log_path: log_path.clone(),
                bootstrap_log_path: bootstrap_log_path.clone(),
                startup_error: Mutex::new(None),
            };

            write_log(&bootstrap_log_path, &format!("[bootstrap] token_generated len={}", token.len()));

            let sidecar_result = app.shell().sidecar("fileeniac-backend");

            match sidecar_result {
                Ok(sidecar_cmd) => {
                    write_log(&bootstrap_log_path, "[bootstrap] sidecar_binary found");

                    let spawn_result = sidecar_cmd
                        .args(["serve", "--host", "127.0.0.1", "--port", "0"])
                        .env("ENIAC_API_TOKEN", &token)
                        .spawn();

                    match spawn_result {
                        Ok((mut rx, child)) => {
                            write_log(&bootstrap_log_path, "[bootstrap] sidecar_spawned pid_ok=true");

                            {
                                let mut b = state.backend.lock().unwrap();
                                b.child = Some(child);
                            }

                            app.manage(state);

                            let app_handle = app.handle().clone();
                            let log_path_clone = log_path.clone();
                            let bootstrap_log_clone = bootstrap_log_path.clone();
                            let port_ready_clone = port_ready.clone();
                            std::thread::spawn(move || {
                                use tauri_plugin_shell::process::CommandEvent;
                                let rt = tokio::runtime::Builder::new_current_thread()
                                    .enable_all()
                                    .build()
                                    .expect("failed to create tokio runtime");
                                rt.block_on(async {
                                    while let Some(event) = rx.recv().await {
                                        match event {
                                            CommandEvent::Stdout(line) => {
                                                let s = String::from_utf8_lossy(&line);
                                                let stripped = strip_token(&s);
                                                write_log(&log_path_clone, &format!("[stdout] {}", stripped));
                                                if s.starts_with("FILEENIAC_READY") {
                                                    for part in s.split_whitespace() {
                                                        if let Some(port_val) = part.strip_prefix("port=") {
                                                            let state = app_handle.state::<AppState>();
                                                            let mut b = state.backend.lock().unwrap();
                                                             b.port = port_val.to_string();
                                                             b.ready = true;
                                                            {
                                                                let mut err = state.startup_error.lock().unwrap();
                                                                *err = None;
                                                            }
                                                             write_log(&bootstrap_log_clone, &format!("[bootstrap] backend_ready port={}", port_val));
                                                            drop(b);
                                                            port_ready_clone.notify_all();
                                                        }
                                                    }
                                                }
                                            }
                                            CommandEvent::Stderr(line) => {
                                                let s = String::from_utf8_lossy(&line);
                                                let stripped = strip_token(&s);
                                                write_log(&log_path_clone, &format!("[stderr] {}", stripped));
                                            }
                                            CommandEvent::Terminated(status) => {
                                                let code = status.code;
                                                write_log(&log_path_clone, &format!("[exit] code={:?}", code));
                                                write_log(&bootstrap_log_clone, &format!("[bootstrap] sidecar_exit code={:?}", code));
                                                let state = app_handle.state::<AppState>();
                                                let mut b = state.backend.lock().unwrap();
                                                b.ready = true;
                                                if b.port.is_empty() {
                                                    let mut err = state.startup_error.lock().unwrap();
                                                    *err = Some(StartupError {
                                                        message: "O servico auxiliar encerrou inesperadamente.".to_string(),
                                                        exit_code: code,
                                                    });
                                                }
                                                drop(b);
                                                port_ready_clone.notify_all();
                                                break;
                                            }
                                            CommandEvent::Error(err) => {
                                                write_log(&log_path_clone, &format!("[error] {}", err));
                                                write_log(&bootstrap_log_clone, &format!("[bootstrap] sidecar_error msg={}", err));
                                                let state = app_handle.state::<AppState>();
                                                let mut b = state.backend.lock().unwrap();
                                                b.ready = true;
                                                {
                                                    let mut startup_err = state.startup_error.lock().unwrap();
                                                    *startup_err = Some(StartupError {
                                                        message: format!("Falha ao iniciar servico auxiliar: {}", err),
                                                        exit_code: None,
                                                    });
                                                }
                                                drop(b);
                                                port_ready_clone.notify_all();
                                                break;
                                            }
                                            _ => {}
                                        }
                                    }
                                });
                            });
                        }
                        Err(e) => {
                            write_log(&bootstrap_log_path, &format!("[bootstrap] sidecar_spawn_failed msg={}", e));
                            {
                                let mut b = state.backend.lock().unwrap();
                                b.ready = true;
                            }
                            {
                                let mut err = state.startup_error.lock().unwrap();
                                *err = Some(StartupError {
                                    message: format!("Nao foi possivel iniciar o servico auxiliar: {}", e),
                                    exit_code: None,
                                });
                            }
                            port_ready.notify_all();
                            app.manage(state);
                        }
                    }
                }
                Err(e) => {
                    write_log(&bootstrap_log_path, &format!("[bootstrap] sidecar_binary_not_found msg={}", e));
                    {
                        let mut b = state.backend.lock().unwrap();
                        b.ready = true;
                    }
                    {
                        let mut err = state.startup_error.lock().unwrap();
                        *err = Some(StartupError {
                            message: "Componente auxiliar nao encontrado no aplicativo.".to_string(),
                            exit_code: None,
                        });
                    }
                    port_ready.notify_all();
                    app.manage(state);
                }
            }

            Ok(())
        })
        .invoke_handler(tauri::generate_handler![get_backend_info, get_sidecar_log_path, get_diagnostics])
        .on_window_event(|window, event| {
            if let tauri::WindowEvent::CloseRequested { .. } = event {
                let state = window.state::<AppState>();
                let mut b = state.backend.lock().unwrap();
                if let Some(child) = b.child.take() {
                    let _ = child.kill();
                }
            }
        })
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}

// SPDX-License-Identifier: MIT
use serde::Serialize;
use std::sync::Mutex;
use tauri::Manager;
use tauri_plugin_shell::ShellExt;
use tauri_plugin_shell::process::CommandChild;

struct BackendState {
    port: String,
    token: String,
    child: Option<CommandChild>,
}

struct AppState {
    backend: Mutex<BackendState>,
    log_path: String,
}

#[derive(Serialize)]
struct BackendInfo {
    base_url: String,
    token: String,
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

#[tauri::command]
fn get_backend_info(state: tauri::State<AppState>) -> BackendInfo {
    let b = state.backend.lock().unwrap();
    BackendInfo {
        base_url: format!("http://127.0.0.1:{}/api", b.port),
        token: b.token.clone(),
    }
}

#[tauri::command]
fn get_sidecar_log_path(state: tauri::State<AppState>) -> String {
    state.log_path.clone()
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

            let state = AppState {
                backend: Mutex::new(BackendState {
                    port: String::new(),
                    token: token.clone(),
                    child: None,
                }),
                log_path: log_path.clone(),
            };

            let sidecar = app
                .shell()
                .sidecar("fileeniac-backend")
                .expect("failed to find fileeniac-backend sidecar binary");

            let (mut rx, child) = sidecar
                .args([
                    "serve",
                    "--host",
                    "127.0.0.1",
                    "--port",
                    "0",
                ])
                .env("ENIAC_API_TOKEN", &token)
                .spawn()
                .expect("failed to spawn fileeniac sidecar");

            {
                let mut b = state.backend.lock().unwrap();
                b.child = Some(child);
            }

            app.manage(state);

            let app_handle = app.handle().clone();
            let log_path_clone = log_path.clone();
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
                                write_log(
                                    &log_path_clone,
                                    &format!("[exit] code={:?}", status.code),
                                );
                                break;
                            }
                            CommandEvent::Error(err) => {
                                write_log(&log_path_clone, &format!("[error] {}", err));
                                break;
                            }
                            _ => {}
                        }
                    }
                });
            });

            Ok(())
        })
        .invoke_handler(tauri::generate_handler![get_backend_info, get_sidecar_log_path])
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

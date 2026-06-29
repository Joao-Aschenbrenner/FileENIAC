// SPDX-License-Identifier: MIT
use std::env;

#[tauri::command]
fn get_api_port() -> String {
    env::var("ENIAC_API_PORT").unwrap_or_else(|_| "8080".to_string())
}

#[tauri::command]
fn get_api_token() -> String {
    env::var("ENIAC_API_TOKEN").unwrap_or_default()
}

#[cfg_attr(mobile, tauri::mobile_entry_point)]
pub fn run() {
    tauri::Builder::default()
        .plugin(tauri_plugin_opener::init())
        .plugin(tauri_plugin_dialog::init())
        .invoke_handler(tauri::generate_handler![get_api_port, get_api_token])
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}

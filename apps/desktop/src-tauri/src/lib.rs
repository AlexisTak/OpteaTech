mod commands;

#[cfg_attr(mobile, tauri::mobile_entry_point)]
pub fn run() {
    tauri::Builder::default()
        .plugin(tauri_plugin_store::Builder::default().build())
        .plugin(tauri_plugin_notification::init())
        .plugin(tauri_plugin_shell::init())
        .plugin(tauri_plugin_dialog::init())
        .plugin(tauri_plugin_fs::init())
        .invoke_handler(tauri::generate_handler![
            commands::auth::store_tokens,
            commands::auth::get_tokens,
            commands::auth::clear_tokens,
            commands::files::upload_file,
            commands::notifications::send_notification,
            commands::system::get_platform,
            commands::system::open_url,
        ])
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}

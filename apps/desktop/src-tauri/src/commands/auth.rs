use serde::{Deserialize, Serialize};
use tauri_plugin_store::StoreExt;

#[derive(Debug, Serialize, Deserialize)]
pub struct TokenPair {
    pub access_token: String,
    pub refresh_token: String,
}

#[tauri::command]
pub async fn store_tokens(
    app: tauri::AppHandle,
    tokens: TokenPair,
) -> Result<(), String> {
    let store = app.store("auth.bin").map_err(|e| e.to_string())?;
    store.set("access_token", tokens.access_token);
    store.set("refresh_token", tokens.refresh_token);
    store.save().map_err(|e| e.to_string())
}

#[tauri::command]
pub async fn get_tokens(
    app: tauri::AppHandle,
) -> Result<Option<TokenPair>, String> {
    let store = app.store("auth.bin").map_err(|e| e.to_string())?;
    let access = store
        .get("access_token")
        .and_then(|v| v.as_str().map(String::from));
    let refresh = store
        .get("refresh_token")
        .and_then(|v| v.as_str().map(String::from));

    match (access, refresh) {
        (Some(a), Some(r)) => Ok(Some(TokenPair {
            access_token: a,
            refresh_token: r,
        })),
        _ => Ok(None),
    }
}

#[tauri::command]
pub async fn clear_tokens(
    app: tauri::AppHandle,
) -> Result<(), String> {
    let store = app.store("auth.bin").map_err(|e| e.to_string())?;
    store.delete("access_token");
    store.delete("refresh_token");
    store.save().map_err(|e| e.to_string())
}

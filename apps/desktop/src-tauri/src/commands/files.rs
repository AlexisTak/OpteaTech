use serde::Serialize;
use std::path::PathBuf;

#[derive(Serialize)]
pub struct UploadedFile {
    pub name: String,
    pub path: String,
    pub size: u64,
    pub mime_type: String,
}

#[tauri::command]
pub async fn upload_file(path: String) -> Result<UploadedFile, String> {
    let path = PathBuf::from(&path);

    let metadata = std::fs::metadata(&path).map_err(|e| format!("Cannot read file: {e}"))?;

    let name = path
        .file_name()
        .and_then(|n| n.to_str())
        .unwrap_or("unknown")
        .to_string();

    let mime = match path.extension().and_then(|e| e.to_str()) {
        Some("pdf") => "application/pdf",
        Some("png") => "image/png",
        Some("jpg") | Some("jpeg") => "image/jpeg",
        Some("zip") => "application/zip",
        Some("svg") => "image/svg+xml",
        _ => "application/octet-stream",
    }
    .to_string();

    Ok(UploadedFile {
        name,
        path: path.to_string_lossy().to_string(),
        size: metadata.len(),
        mime_type: mime,
    })
}

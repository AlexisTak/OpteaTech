import { getCurrentWindow } from '@tauri-apps/api/window';

export function Titlebar() {
  const win = getCurrentWindow();

  return (
    <div className="titlebar" data-tauri-drag-region>
      <strong style={{ fontSize: 13 }}>optea.tech Admin</strong>
      <div className="row" onMouseDown={(e) => e.stopPropagation()}>
        <button className="button secondary" onClick={() => win.minimize()}>_</button>
        <button className="button secondary" onClick={() => win.toggleMaximize()}>[]</button>
        <button className="button secondary" onClick={() => win.close()}>X</button>
      </div>
    </div>
  );
}

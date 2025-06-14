import useSWR, { mutate } from 'swr';
import { useState } from 'react';
const fetcher = (u: string) => fetch(u).then(r => r.json());

export default function Home() {
  const { data: imgs = [] } = useSWR('/api/images', fetcher);
  const [sel, setSel]       = useState<File | null>(null);
  const [busy, setBusy]     = useState(false);

  const choose = (e: React.ChangeEvent<HTMLInputElement>) =>
    setSel(e.target.files?.[0] ?? null);

  async function upload() {
    if (!sel) return;
    setBusy(true);
    const fd = new FormData();
    fd.append('image', sel);           /* ← name="image" を統一 */ 
    await fetch('/api/upload', { method: 'POST', body: fd });
    await mutate('/api/images');       /* ↺ 自動再描画 */  
    setSel(null); setBusy(false);
  }

  return (
    <main style={{ maxWidth: 840, margin: '0 auto', fontFamily: 'sans-serif' }}>
      <h1>画像アップロード</h1>

      <input type="file" accept="image/*" onChange={choose} />
      <button onClick={upload} disabled={!sel || busy}>アップロード</button>

      <h2>画像一覧</h2>
      {imgs.length === 0 && <p>画像がありません</p>}

      <div style={{ display: 'flex', flexWrap: 'wrap' }}>
        {imgs.map((i: any) => (
          <div key={i.name} style={{ margin: 8, textAlign: 'center' }}>
            <img src={i.url} alt={i.name}
                 style={{ maxWidth: 200, border: '1px solid #ccc' }} />
            <p style={{ fontSize: 12 }}>{i.name}</p>
          </div>
        ))}
      </div>
    </main>
  );
}

import { minio } from '@/lib/minio';

export default async function handler(_req, res) {
  /* バケット名を環境変数から取得 */
  const bucket = process.env.MINIO_BUCKET!;

  /* MinIO が返すオブジェクトストリームを変数に受け取る */
  const stream = minio.listObjectsV2(bucket, '', true);

  const imgs: { name: string; url: string }[] = [];
  const base = `${process.env.NEXT_PUBLIC_MINIO_PUBLIC_URL}/${bucket}/`;

  /* data イベントで 1 オブジェクトずつ収集 */
  stream.on('data', obj => {
    imgs.push({
      name: obj.name,
      url: base + obj.name,        // プリサインしない場合は固定 URL
    });
  });

  /* 完全に読み終えたら 200 を返す */
  stream.on('end', () => {
    res.status(200).json(imgs);
  });

  /* 途中でエラーがあれば 500 を返す */
  stream.on('error', err => {
    console.error('listObjectsV2 error:', err);
    res.status(500).json({ error: 'MinIO list error' });
  });
}

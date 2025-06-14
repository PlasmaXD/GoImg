import formidable, { File } from "formidable";
import { minio } from "@/lib/minio";

export const config = { api: { bodyParser: false } };

export default async function handler(req, res) {
  if (req.method !== "POST") return res.status(405).end();

  /* ① multipart 解析 */
  const form = formidable({ keepExtensions: true, multiples: false });
  const { files } = await new Promise<{ files: formidable.Files }>(
    (ok, ng) => form.parse(req, (err, _f, files) => (err ? ng(err) : ok({ files })))
  );

  /* ② name="image" を **必ず** 付ける */
  const fileObj = files.image as File | File[] | undefined;
  if (!fileObj) return res.status(400).json({ error: "no file" });

  const file = Array.isArray(fileObj) ? fileObj[0] : fileObj;
  const tempPath = file.filepath ?? file.path;   // v1/v2 どちらでも可
  if (typeof tempPath !== "string")
    return res.status(400).json({ error: "invalid file path" });

  /* ③ MinIO へ PUT */
  const key  = `${Date.now()}_${file.originalFilename}`;
  await minio.fPutObject(                       // filePath は string 必須
    process.env.MINIO_BUCKET!,
    key,
    tempPath,
    {}
  );

  const url = await minio.presignedGetObject(   // 公開 1h の署名付 URL﻿:contentReference[oaicite:4]{index=4}
    process.env.MINIO_BUCKET!,
    key,
    60 * 60
  );
  res.status(200).json({ key, url });
}

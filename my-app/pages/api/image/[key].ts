import { minio } from "@/lib/minio";

export default async function handler(req, res) {
  if (req.method !== "DELETE") return res.status(405).end();
  const { key } = req.query;

  await minio.removeObject("processed-images", key as string); // DELETE API:contentReference[oaicite:7]{index=7}
  res.status(204).end();
}

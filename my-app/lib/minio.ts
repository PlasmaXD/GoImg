import { Client } from "minio";

export const minio = new Client({
  endPoint: process.env.MINIO_ENDPOINT!,   // localhost
  port:     Number(process.env.MINIO_PORT),// 9000
  useSSL:   false,
  accessKey: process.env.MINIO_ACCESS_KEY!,
  secretKey: process.env.MINIO_SECRET_KEY!,
});

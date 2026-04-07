import { createDecipheriv } from 'node:crypto';

const key = Buffer.from('Example');
const iv = Buffer.from('Example');
const encryptedText = Buffer.from('Example', 'utf8');
export async function decrypt() {
  const decipher = createDecipheriv('aes-256-ctr', key, iv);

  const decryptedText = Buffer.concat([
    decipher.update(encryptedText),
    decipher.final(),
  ]);

  return decryptedText;
}

import { createCipheriv, randomBytes, scrypt } from 'node:crypto';
import { promisify } from 'node:util';

export async function encrypt() {
  const iv = randomBytes(16);
  const password = 'Password used to generate key';

  const key = (await promisify(scrypt)(password, 'salt', 32)) as Buffer;
  const cipher = createCipheriv('aes-256-ctr', key, iv);

  const textToEncrypt = 'Example';

  const encryptedText = Buffer.concat([
    cipher.update(textToEncrypt),
    cipher.final(),
  ]);

  return encryptedText;
}

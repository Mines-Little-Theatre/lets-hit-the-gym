import { decodeHex } from "../../util/hex.js";

const encoder = new TextEncoder();

export async function validateInteraction(
  body: ArrayBuffer,
  signature: string | null,
  timestamp: string | null,
  publicKey: string,
): Promise<boolean> {
  try {
    if (
      !signature ||
      !timestamp ||
      // fail if request is more than 3 seconds old
      Date.now() / 1000 > Number.parseInt(timestamp) + 3
    ) {
      return false;
    }

    const signatureBytes = decodeHex(signature);
    const timestampBytes = encoder.encode(timestamp);
    const publicKeyKey = await crypto.subtle.importKey(
      "raw",
      decodeHex(publicKey),
      "Ed25519",
      false,
      ["verify"],
    );

    const data = new Uint8Array(timestampBytes.length + body.byteLength);
    data.set(timestampBytes);
    data.set(new Uint8Array(body), timestampBytes.length);

    return crypto.subtle.verify("Ed25519", publicKeyKey, signatureBytes, data);
  } catch (e) {
    console.warn("validateInteraction", e);
    return false;
  }
}

const ASCII_0 = 48;
const ASCII_9 = 57;
const ASCII_A = 65;
const ASCII_Z = 90;
const ASCII_a = 97;
const ASCII_z = 122;

function decodeHexit(hexit: number): number {
  if (ASCII_0 <= hexit && hexit <= ASCII_9) {
    return hexit - ASCII_0;
  } else if (ASCII_A <= hexit && hexit <= ASCII_Z) {
    return hexit - ASCII_A + 10;
  } else if (ASCII_a <= hexit && hexit <= ASCII_z) {
    return hexit - ASCII_a + 10;
  } else {
    throw new Error("invalid hex character");
  }
}

export function decodeHex(hex: string): Uint8Array {
  if (hex.length % 2) {
    throw new Error("decodeHex called on string with odd length");
  }

  const buf = new Uint8Array(hex.length / 2);
  for (let i = 0; i < buf.length; i++) {
    buf[i] =
      (decodeHexit(hex.charCodeAt(2 * i)) << 4) |
      decodeHexit(hex.charCodeAt(2 * i + 1));
  }
  return buf;
}

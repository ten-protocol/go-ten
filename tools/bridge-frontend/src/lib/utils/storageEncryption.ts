import * as CryptoJS from "crypto-js";

const _key = "TEN";

function encrypt(txt: any) {
  return CryptoJS.AES.encrypt(txt, _key).toString();
}

function decrypt(txtToDecrypt: any) {
  console.log("🚀 ~ decrypt ~ txtToDecrypt:", txtToDecrypt);
  return CryptoJS.AES.decrypt(txtToDecrypt, _key).toString(CryptoJS.enc.Utf8);
}

export { encrypt, decrypt };

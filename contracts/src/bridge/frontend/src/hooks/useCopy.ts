import { showToast } from "@/src/components/ui/use-toast";
import { ToastType } from "@/types";
import React from "react";

export const useCopy = () => {
  const [copied, setCopied] = React.useState(false);

  const copyToClipboard = (text: string, name?: string) => {
    copyText(text)
      .catch(() => fallbackCopyTextToClipboard(text))
      .then(() => {
        showToast(ToastType.DESTRUCTIVE, `${name ? name : "Copied!"}`);
        setCopied(true);
      })
      .catch(() => {
        showToast(
          ToastType.DESTRUCTIVE,
          `Couldn't copy ${name ? name : "Text"}!!!`
        );
      })
      .finally(() => {
        setTimeout(() => setCopied(false), 1000);
      });
  };

  return {
    copyToClipboard,
    copied,
  };
};

const copyText = async (text: string) => {
  return navigator.clipboard.writeText(text);
};

const fallbackCopyTextToClipboard = (text: string) => {
  return new Promise((res, rej) => {
    var textArea = document.createElement("textarea");
    textArea.value = text;

    // Avoid scrolling to bottom
    textArea.style.top = "0";
    textArea.style.left = "0";
    textArea.style.position = "fixed";

    document.body.appendChild(textArea);
    textArea.focus();
    textArea.select();

    try {
      document.execCommand("copy");
    } catch (err) {
      rej(err);
    }

    document.body.removeChild(textArea);
    res(text);
  });
};

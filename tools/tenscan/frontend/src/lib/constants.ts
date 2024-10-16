export const pollingInterval = 5000;
export const maxRetries = 3;
export const pricePollingInterval = 60 * 1000;

export const RESET_COPIED_TIMEOUT = 2000;

const calculateOffset = (page: number, size: number) => {
  if (page <= 0) return 0;
  return (page - 1) * size;
};

export const getOptions = (query: { page?: number; size?: number }) => {
  const defaultSize = 20;
  const size = query.size
    ? +(query.size > 100 ? 100 : query.size)
    : defaultSize;
  const page = query.page || 1;
  const offset = calculateOffset(page, size);
  return {
    offset: Number.isNaN(offset) || offset < 0 ? 0 : offset,
    size,
  };
};

export const version = process.env.NEXT_PUBLIC_FE_VERSION;
export const apiHost = process.env.NEXT_PUBLIC_API_HOST;

export const currentEncryptedKey =
  "bddbc0d46a0666ce57a466168d99c1830b0c65e052d77188f2cbfc3f6486588c";

export const GOOGLE_ANALYTICS_ID = "G-M82QX9RT4L";

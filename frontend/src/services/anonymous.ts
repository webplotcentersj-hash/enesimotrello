const STORAGE_KEY = 'anonymous_user_id';

const generateUUID = (): string => {
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, (c) => {
    const r = (Math.random() * 16) | 0;
    const v = c === 'x' ? r : (r & 0x3) | 0x8;
    return v.toString(16);
  });
};

export const ensureAnonymousUserId = (): string => {
  if (typeof window === 'undefined') {
    return '';
  }

  let userId = window.localStorage.getItem(STORAGE_KEY);

  if (!userId) {
    userId = generateUUID();
    window.localStorage.setItem(STORAGE_KEY, userId);
  }

  return userId;
};

export const getAnonymousDisplayName = (userId: string): string => {
  if (!userId) {
    return 'Guest';
  }

  return `Guest ${userId.slice(0, 8).toUpperCase()}`;
};



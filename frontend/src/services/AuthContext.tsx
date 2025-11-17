import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { ensureAnonymousUserId, getAnonymousDisplayName } from './anonymous.ts';

interface AuthContextType {
  anonymousUserId: string;
  displayName: string;
  loading: boolean;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};

interface AuthProviderProps {
  children: ReactNode;
}

export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
  const [anonymousUserId, setAnonymousUserId] = useState<string>('');
  const [displayName, setDisplayName] = useState<string>('Guest');
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const userId = ensureAnonymousUserId();
    setAnonymousUserId(userId);
    setDisplayName(getAnonymousDisplayName(userId));
    setLoading(false);
  }, []);

  const value = {
    anonymousUserId,
    displayName,
    loading,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};


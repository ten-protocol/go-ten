import { useState, useEffect } from 'react';

/**
 * A hook for working with localStorage that provides type safety
 * @param key The localStorage key to manage
 * @param initialValue The initial value to use if no value is found in localStorage
 * @returns A tuple containing the current value and functions to update and remove it
 */
export function useLocalStorage<T>(
  key: string,
  initialValue: T
): [T, (value: T) => void, () => void] {
  // Get the stored value from localStorage, or use the initial value
  const readValue = (): T => {
    // Check if we're in a browser environment
    if (typeof window === 'undefined') {
      return initialValue;
    }

    try {
      const item = window.localStorage.getItem(key);
      return item ? (JSON.parse(item) as T) : initialValue;
    } catch (error) {
      console.warn(`Error reading localStorage key "${key}":`, error);
      return initialValue;
    }
  };

  // State to hold the current value
  const [storedValue, setStoredValue] = useState<T>(readValue);

  // Update localStorage when the state changes
  const setValue = (value: T) => {
    try {
      // Save state
      setStoredValue(value);
      
      // Save to localStorage
      if (typeof window !== 'undefined') {
        window.localStorage.setItem(key, JSON.stringify(value));
        // Dispatch an event so other instances of the hook know to update
        window.dispatchEvent(new Event('local-storage'));
      }
    } catch (error) {
      console.warn(`Error setting localStorage key "${key}":`, error);
    }
  };

  // Remove the item from localStorage
  const removeValue = () => {
    try {
      if (typeof window !== 'undefined') {
        window.localStorage.removeItem(key);
        // Reset to initial value
        setStoredValue(initialValue);
        // Dispatch an event so other instances know to update
        window.dispatchEvent(new Event('local-storage'));
      }
    } catch (error) {
      console.warn(`Error removing localStorage key "${key}":`, error);
    }
  };

  // Listen for changes to this localStorage item in other windows/tabs
  useEffect(() => {
    const handleStorageChange = () => {
      setStoredValue(readValue());
    };
    
    // Listen for the custom event and storage events
    window.addEventListener('local-storage', handleStorageChange);
    window.addEventListener('storage', handleStorageChange);
    
    return () => {
      window.removeEventListener('local-storage', handleStorageChange);
      window.removeEventListener('storage', handleStorageChange);
    };
  }, [key]);

  return [storedValue, setValue, removeValue];
} 
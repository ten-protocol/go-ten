'use client';

import { useState, useEffect, useRef } from 'react';
import { useQuery } from '@tanstack/react-query';
import { useRouter } from 'next/navigation';
import { searchRecords } from '@/api/search';
import { SearchResponse, SearchResult } from '@/src/types/interfaces/SearchInterfaces';
import { ResponseDataInterface } from '@repo/ui/lib/types/common';
import { pageLinks } from '@/src/routes';
import { pathToUrl } from '@/src/routes/router';
import { Search, Loader2, X } from '@repo/ui/components/shared/react-icons';
import SearchResultItem from './search-result-item';

export default function SearchBar() {
    const [query, setQuery] = useState('');
    const [debouncedQuery, setDebouncedQuery] = useState('');
    const [isDropdownOpen, setIsDropdownOpen] = useState(false);
    const [focusedIndex, setFocusedIndex] = useState(-1);
    const searchRef = useRef<HTMLDivElement>(null);
    const inputRef = useRef<HTMLInputElement>(null);
    const router = useRouter();

    // Debounce the search query
    useEffect(() => {
        const timer = setTimeout(() => {
            setDebouncedQuery(query);
        }, 1000);

        return () => clearTimeout(timer);
    }, [query]);

    // Close dropdown when clicking outside
    useEffect(() => {
        const handleClickOutside = (event: MouseEvent) => {
            if (searchRef.current && !searchRef.current.contains(event.target as Node)) {
                setIsDropdownOpen(false);
                setFocusedIndex(-1);
            }
        };

        document.addEventListener('mousedown', handleClickOutside);
        return () => document.removeEventListener('mousedown', handleClickOutside);
    }, []);

    // Fetch search results
    const { data: searchResponse, isLoading, error, isFetching } = useQuery({
        queryKey: ['search', debouncedQuery],
        queryFn: () => searchRecords(debouncedQuery),
        enabled: debouncedQuery.length >= 2,
        staleTime: 5 * 60 * 1000, // 5 minutes
    });

    // Extract search results from the response
    const searchResults = searchResponse?.result;

    // Debug logging
    console.log('Search state:', { query, debouncedQuery, isLoading, isFetching, hasResults: !!searchResults?.ResultsData?.length });

    // Handle input change
    const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const value = e.target.value;
        setQuery(value);
        setIsDropdownOpen(value.length >= 2);
        setFocusedIndex(-1);
    };

    // Handle clear button click
    const handleClear = () => {
        setQuery('');
        setDebouncedQuery('');
        setIsDropdownOpen(false);
        setFocusedIndex(-1);
        inputRef.current?.focus();
    };

    // Handle keyboard navigation
    const handleKeyDown = (e: React.KeyboardEvent) => {
        if (!searchResults?.ResultsData) return;

        switch (e.key) {
            case 'ArrowDown':
                e.preventDefault();
                setFocusedIndex(prev => 
                    prev < searchResults.ResultsData.length - 1 ? prev + 1 : prev
                );
                break;
            case 'ArrowUp':
                e.preventDefault();
                setFocusedIndex(prev => prev > 0 ? prev - 1 : -1);
                break;
            case 'Enter':
                e.preventDefault();
                if (focusedIndex >= 0 && searchResults.ResultsData[focusedIndex]) {
                    handleResultClick(searchResults.ResultsData[focusedIndex]);
                }
                break;
            case 'Escape':
                setIsDropdownOpen(false);
                setFocusedIndex(-1);
                inputRef.current?.blur();
                break;
        }
    };

    // Handle result click
    const handleResultClick = (result: SearchResult) => {
        let route = '';
        
        switch (result.type) {
            case 'transaction':
                route = pathToUrl(pageLinks.txByHash, { hash: result.hash });
                break;
            case 'batch':
                route = pathToUrl(pageLinks.batchByHash, { hash: result.hash });
                break;
            case 'rollup':
                route = pathToUrl(pageLinks.rollupByHash, { hash: result.hash });
                break;
            default:
                return;
        }

        router.push(route);
        setIsDropdownOpen(false);
        setQuery('');
        setFocusedIndex(-1);
    };

    return (
        <div className="relative w-full" ref={searchRef}>
            <div className="relative">
                <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-5 h-5" />
                <input
                    ref={inputRef}
                    type="text"
                    placeholder="Search transactions, batches, or rollups..."
                    value={query}
                    onChange={handleInputChange}
                    onKeyDown={handleKeyDown}
                    onFocus={() => query.length >= 2 && setIsDropdownOpen(true)}
                    className="w-full pl-10 pr-4 py-3 border border-white/5 focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none transition-all duration-200"
                />
                {(isLoading || isFetching) && (
                    <Loader2 className="absolute right-3 top-1/2 transform -translate-y-1/2 text-blue-500 w-5 h-5 animate-spin" />
                )}
                {query.length > 0 && !isLoading && !isFetching && (
                    <button
                        onClick={handleClear}
                        className="absolute right-3 top-1/2 transform -translate-y-1/2 text-gray-400 hover:text-gray-600 transition-colors duration-200 p-1 rounded-full hover:bg-gray-100"
                        type="button"
                        aria-label="Clear search"
                    >
                        <X className="w-4 h-4" />
                    </button>
                )}
            </div>

            {/* Dropdown Results */}
            {isDropdownOpen && query.length >= 2 && (
                <div className="absolute top-full left-0 right-0 mt-2 bg-white border border-gray-200 rounded-lg shadow-lg z-50 min-h-[100px]">
                    {error ? (
                        <div className="p-4 text-red-500 text-center border-b border-gray-100">
                            Error loading search results
                        </div>
                    ) : (isLoading || isFetching) ? (
                        <div className="p-6 text-center border-b border-gray-100">
                            <Loader2 className="w-8 h-8 animate-spin mx-auto mb-3 text-blue-500" />
                            <p className="text-gray-600 font-medium">Searching...</p>
                            <p className="text-sm text-gray-500 mt-1">Looking for "{debouncedQuery}"</p>
                        </div>
                    ) : searchResults?.ResultsData && searchResults.ResultsData.length > 0 ? (
                        <div className="max-h-96 overflow-y-auto">
                            {searchResults.ResultsData.map((result: SearchResult, index: number) => (
                                <SearchResultItem
                                    key={`${result.type}-${result.hash}`}
                                    result={result}
                                    index={index}
                                    focusedIndex={focusedIndex}
                                    onResultClick={handleResultClick}
                                    onMouseEnter={setFocusedIndex}
                                />
                            ))}
                            {searchResults.Total > searchResults.ResultsData.length && (
                                <div className="p-3 text-center text-sm text-gray-500 border-t border-gray-100 bg-gray-50">
                                    Showing {searchResults.ResultsData.length} of {searchResults.Total} results
                                </div>
                            )}
                        </div>
                    ) : debouncedQuery.length >= 2 && !isLoading && !isFetching ? (
                        <div className="p-4 text-gray-500 text-center border-b border-gray-100">
                            No results found for &quot;{debouncedQuery}&quot;
                        </div>
                    ) : query !== debouncedQuery ? (
                        <div className="p-6 text-center border-b border-gray-100">
                            <Loader2 className="w-8 h-8 animate-spin mx-auto mb-3 text-blue-500" />
                            <p className="text-gray-600 font-medium">Searching...</p>
                            <p className="text-sm text-gray-500 mt-1">Looking for "{query}"</p>
                        </div>
                    ) : (
                        <div className="p-4 text-gray-400 text-center">
                            <p className="text-sm">Type to search...</p>
                        </div>
                    )}
                </div>
            )}
        </div>
    );
}
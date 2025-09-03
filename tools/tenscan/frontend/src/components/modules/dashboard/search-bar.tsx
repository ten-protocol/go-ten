'use client';

import React, { useState, useEffect, useRef } from 'react';
import { useQuery } from '@tanstack/react-query';
import { useRouter } from 'next/navigation';
import { searchRecords } from '@/api/search';
import {  SearchResult } from '@/src/types/interfaces/SearchInterfaces';
import { pageLinks } from '@/src/routes';
import { pathToUrl } from '@/src/routes/router';
import { Search, Loader2, X } from '@repo/ui/components/shared/react-icons';
import SearchResultItem from './search-result-item';
import {Input} from "@/src/components/ui/input";

export default function SearchBar() {
    const [query, setQuery] = useState('');
    const [debouncedQuery, setDebouncedQuery] = useState('');
    const [isDropdownOpen, setIsDropdownOpen] = useState(false);
    const [focusedIndex, setFocusedIndex] = useState(-1);
    const searchRef = useRef<HTMLDivElement>(null);
    const inputRef = useRef<HTMLInputElement>(null);
    const router = useRouter();

    useEffect(() => {
        const timer = setTimeout(() => {
            setDebouncedQuery(query);
        }, 1000);

        return () => clearTimeout(timer);
    }, [query]);


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

    const { data: searchResponse, isLoading, error, isFetching } = useQuery({
        queryKey: ['search', debouncedQuery],
        queryFn: () => searchRecords(debouncedQuery),
        enabled: debouncedQuery.length >= 2,
        staleTime: 5 * 60 * 1000, // 5 minutes
    });

    const searchResults = searchResponse?.result;

    const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const value = e.target.value;
        setQuery(value);
        setIsDropdownOpen(value.length >= 2);
        setFocusedIndex(-1);
    };

    const handleClear = () => {
        setQuery('');
        setDebouncedQuery('');
        setIsDropdownOpen(false);
        setFocusedIndex(-1);
        inputRef.current?.focus();
    };

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
        <div className="relative w-full max-w-2xl mx-auto p-4" ref={searchRef}>
            <div className="search-container-shape absolute inset-[1px] pointer-events-none">
                <div className="absolute inset-0 animate-scan-overlay"/>
            </div>


            <div className="relative">
                <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-white/70 w-5 h-5 z-10"/>
                <Input
                    ref={inputRef}
                    type="text"
                    placeholder="Search transactions, batches, or rollups..."
                    value={query}
                    onChange={handleInputChange}
                    onKeyDown={handleKeyDown}
                    onFocus={() => query.length >= 2 && setIsDropdownOpen(true)}
                    className="w-full p-6 pl-10 bg-white/5"
                />
                {(isLoading || isFetching) && (
                    <Loader2
                        className="absolute right-3 top-1/2 transform -translate-y-1/2 text-accent w-5 h-5 animate-spin"/>
                )}
                {query.length > 0 && !isLoading && !isFetching && (
                    <button
                        onClick={handleClear}
                        className="absolute right-3 top-1/2 transform -translate-y-1/2 text-gray-400 hover:text-gray-600 transition-colors duration-200 p-1 rounded-full hover:bg-gray-100"
                        type="button"
                        aria-label="Clear search"
                    >
                        <X className="w-4 h-4"/>
                    </button>
                )}
            </div>

            {isDropdownOpen && query.length >= 2 && (
                <div
                    className="absolute top-full left-0 right-0 mt-2 bg-background/80 backdrop-blur border border-white/5 shadow-lg z-50 min-h-[60px]">
                    <div className="absolute inset-0 animate-scan-overlay opacity-10 pointer-events-none"/>


                    {error ? (
                        <div className="p-4 text-red-500 text-center">
                            Error loading search results
                        </div>
                    ) : (isLoading || isFetching) ? (
                        <div className="p-6 text-center">
                            <Loader2 className="w-8 h-8 animate-spin mx-auto mb-3 text-white"/>
                            <p className="text-gray-600 font-medium">Searching...</p>
                            <p className="text-sm text-gray-500 mt-1">Looking for &quot;{debouncedQuery}&quot;</p>
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
                                <div
                                    className="p-3 text-center text-sm text-white">
                                    Showing {searchResults.ResultsData.length} of {searchResults.Total} results
                                </div>
                            )}
                        </div>
                    ) : debouncedQuery.length >= 2 && !isLoading && !isFetching ? (
                        <div className="p-4 text-white/80 text-center">
                            No results found for &quot;{debouncedQuery}&quot;
                        </div>
                    ) : query !== debouncedQuery ? (
                        <div className="p-6 text-center">
                            <Loader2 className="w-8 h-8 animate-spin mx-auto mb-3 text-accent"/>
                            <p className="text-white/80 font-medium">Searching...</p>
                            <p className="text-sm text-white/60 mt-1">Looking for &quot;{query}&quot;</p>
                        </div>
                    ) : (
                        <div className="p-4 text-white/60 text-center">
                            <p className="text-sm">Type to search...</p>
                        </div>
                    )}
                </div>
            )}
        </div>
    );
}
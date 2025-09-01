import { SearchResult } from '@/src/types/interfaces/SearchInterfaces';
import { Hash, Layers, FileText } from '@repo/ui/components/shared/react-icons';

interface SearchResultItemProps {
    result: SearchResult;
    index: number;
    focusedIndex: number;
    onResultClick: (result: SearchResult) => void;
    onMouseEnter: (index: number) => void;
}

export default function SearchResultItem({
    result,
    index,
    focusedIndex,
    onResultClick,
    onMouseEnter,
}: SearchResultItemProps) {
    // Get icon for result type
    const getResultIcon = (type: string) => {
        switch (type) {
            case 'transaction':
                return <FileText className="w-4 h-4 text-blue-500" />;
            case 'batch':
                return <Layers className="w-4 h-4 text-green-500" />;
            case 'rollup':
                return <Hash className="w-4 h-4 text-purple-500" />;
            default:
                return <Hash className="w-4 h-4 text-gray-500" />;
        }
    };

    // Format timestamp
    const formatTimestamp = (timestamp: number) => {
        return new Date(timestamp * 1000).toLocaleDateString();
    };

    return (
        <div
            className={`p-4 cursor-pointer hover:bg-gray-50 transition-colors duration-150 border-b border-gray-100 last:border-b-0 ${
                index === focusedIndex ? 'bg-blue-50 border-l-4 border-l-blue-500' : ''
            }`}
            onClick={() => onResultClick(result)}
            onMouseEnter={() => onMouseEnter(index)}
        >
            <div className="flex items-center space-x-3">
                {getResultIcon(result.type)}
                <div className="flex-1 min-w-0">
                    <div className="flex items-center space-x-2">
                        <span className="text-sm font-medium text-gray-900 capitalize">
                            {result.type}
                        </span>
                        {result.height && (
                            <span className="text-xs text-gray-500">
                                Height: {result.height.toString()}
                            </span>
                        )}
                        {result.sequence && (
                            <span className="text-xs text-gray-500">
                                Seq: {result.sequence.toString()}
                            </span>
                        )}
                    </div>
                    <div className="text-sm text-gray-600 font-mono truncate">
                        {result.hash}
                    </div>
                    {result.timestamp && (
                        <div className="text-xs text-gray-500">
                            {formatTimestamp(result.timestamp)}
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
}

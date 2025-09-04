import { SearchResult } from '@/src/types/interfaces/SearchInterfaces';
import {Hash, Layers, FileText, CubeIcon} from '@repo/ui/components/shared/react-icons';

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
    const getResultIcon = (type: string) => {
        switch (type) {
            case 'transaction':
                return <FileText className="w-4 h-4 text-blue-300" />;
            case 'batch':
                return <Layers className="w-4 h-4 text-red-400" />;
            case 'rollup':
                return <CubeIcon className="w-4 h-4 text-purple-300" />;
            default:
                return <Hash className="w-4 h-4 text-gray-500" />;
        }
    };

    const formatTimestamp = (timestamp: number) => {
        return `${new Date(timestamp * 1000).toLocaleTimeString()} - ${new Date(timestamp * 1000).toLocaleDateString()}`
    };

    return (
        <div
            className={`p-4 cursor-pointer hover:bg-white/10 hover:pl-3 hover:border-l-4 border-white/50 transition-colors duration-150 ${
                index === focusedIndex ? 'bg-white/10' : ''
            }`}
            onClick={() => onResultClick(result)}
            onMouseEnter={() => onMouseEnter(index)}
        >
            <div className="flex items-center space-x-3">
                {getResultIcon(result.type)}
                <div className="flex-1 min-w-0">
                    <div className="flex items-center space-x-2">
                        <span className="text-sm font-medium text-white capitalize">
                            {result.type}
                        </span>
                        {result.height && (
                            <span className="text-xs text-white/60">
                                Height: {result.height.toString()}
                            </span>
                        )}
                        {result.sequence && (
                            <span className="text-xs text-white/60">
                                Seq: {result.sequence.toString()}
                            </span>
                        )}
                    </div>
                    <div className="text-sm text-white/80 font-mono truncate">
                        {result.hash}
                    </div>
                    {result.timestamp && (
                        <div className="text-xs text-white/60">
                            {formatTimestamp(result.timestamp)}
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
}

import { useMutation } from '@tanstack/react-query';
import { getUserID } from '@/api/gateway';
import { useLocalStorage } from 'usehooks-ts';

export function useGetUserID() {
    const [tenToken] = useLocalStorage<string | null>('ten_token', null);
    
    const mutation = useMutation({
        mutationFn: () => {
            if (!tenToken) {
                throw new Error('TEN token not found. Please authenticate first.');
            }
            return getUserID(tenToken);
        },
    });

    const getUserId = async () => {
        try {
            const result = await mutation.mutateAsync();
            return result;
        } catch (error) {
            console.error('Error getting user ID:', error);
            throw error;
        }
    };

    return {
        getUserId,
        isLoading: mutation.isPending,
        error: mutation.error,
        isSuccess: mutation.isSuccess,
        data: mutation.data,
    };
} 
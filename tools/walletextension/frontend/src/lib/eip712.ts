import { TypedDataDomain, TypedDataParameter } from 'viem';
import { tenChainIDDecimal } from '@/lib/constants';

export const generateEIP712 = (
    token: `0x${string}`
): {
    domain: TypedDataDomain;
    types: Record<string, TypedDataParameter[]>;
    primaryType: string;
    message: Record<string, string>;
} => ({
    types: {
        EIP712Domain: [
            { name: 'name', type: 'string' },
            { name: 'version', type: 'string' },
            { name: 'chainId', type: 'uint256' },
            { name: 'verifyingContract', type: 'address' },
        ],
        Authentication: [{ name: 'Encryption Token', type: 'address' }],
    },
    primaryType: 'Authentication',
    domain: {
        name: 'Ten',
        version: '1.0',
        chainId: tenChainIDDecimal,
        verifyingContract: '0x0000000000000000000000000000000000000000',
    },
    message: {
        'Encryption Token': token,
    },
});

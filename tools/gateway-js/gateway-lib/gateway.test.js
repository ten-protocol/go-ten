const axios = require('axios');
const Gateway = require('./gateway.js');

// Mocking axios module
jest.mock('axios');

// Mocking the global fetch function
global.fetch = jest.fn();

describe('Gateway', () => {
    let gateway;
    const httpURL = 'http://example.com';
    const wsURL = 'ws://example.com';
    const provider = {}; // Placeholder for the provider if it's needed in future tests.

    beforeEach(() => {
        gateway = new Gateway(httpURL, wsURL, provider);
    });

    it('should join successfully', async () => {
        axios.get.mockResolvedValue({
            status: 200,
            data: 'testUserID',
        });

        await gateway.join();

        expect(gateway.userId).toBe('testUserID');
    });

    it('should throw error on unsuccessful join', async () => {
        axios.get.mockRejectedValue(new Error('Network error'));

        await expect(gateway.join()).rejects.toThrow('Failed to get userID. Error: Network error');
    });

    it('should register account successfully', async () => {
        axios.post.mockResolvedValue({
            status: 200,
            text: jest.fn().mockResolvedValue('Account registered')
        });

        gateway = new Gateway(httpURL, wsURL, {
            request: jest.fn().mockResolvedValue("mockSignature")
        })

        const result = await gateway.registerAccount( 'address');

        expect(result).toBe('Account registered');
    });

    it('should throw error on unsuccessful account registration', async () => {
        gateway = new Gateway(httpURL, wsURL, {
            request: jest.fn().mockRejectedValue(new Error('Signature error')),
        })

        await expect(gateway.registerAccount('address')).rejects.toThrow('Failed to sign message. Error: Signature error');
    });
});

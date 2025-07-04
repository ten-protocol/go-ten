# Ten Gateway

Ten Gateway is a Next.js and Tailwind CSS-powered application.

## Folder Structure

```
📁 Ten Gateway
├── 📁 api - Contains server-side code, such as API routes or server logic
├── 📁 public - This directory is used to serve static assets. Files inside this directory can be referenced in your code with a URL path
├── 📁 src - Main source code directory for this project
│   ├── 📁 components - Contains reusable React components used throughout the application
│   ├── 📁 pages - Typically used for Next.js pages. Each .tsx or .js file in this directory becomes a route in your application
│   ├── 📁 hooks - Custom React hooks that can be shared and reused across components
│   ├── 📁 lib - Utility functions or modules that provide common functionalities across the application
│   ├── 📁 routes - Route-related logic or configuration can be placed in this directory
│   ├── 📁 services - Used for services that interact with external APIs or handle other data-related tasks
│   └── 📁 types - Type definitions (.d.ts files or TypeScript files) for TypeScript, describing the shape of data and objects used in the application
└── 📁 styles - Global styles, stylesheets, or styling-related configurations for this project
```

## Getting Started

1. **Clone the Repository:**

    ```bash
    git clone https://github.com/ten-protocol/go-ten.git
    cd go-ten/tools/walletextension/frontend
    ```

2. **Install Dependencies:**

    ```bash
    pnpm install
    ```

3. **Configure Environment Variables:**
   Create a `.env.local` file in the root directory of the project and add the following environment variables:

    ```bash
    # Gateway URL - this will be exposed to the client
    NEXT_PUBLIC_API_GATEWAY_URL=********

    # HMAC Secret - this will NOT be exposed to the client (server-side only)
    # Must be a 32-byte hex-encoded string that matches the backend --hmacSecret flag
    HMAC_SECRET=your_32_byte_hex_encoded_secret_here
    ```

    Possible values for `NEXT_PUBLIC_API_GATEWAY_URL` are:

    - `https://uat-testnet.ten.xyz`
    - `https://sepolia-testnet.ten.xyz`
    - `https://dev-testnet.ten.xyz`

    **Important:** The `HMAC_SECRET` must match the same value used in the backend `--hmacSecret` flag. This secret is used for HMAC generation in the GetUserID functionality and is never exposed to the client-side JavaScript.

4. **Run the Development Server:**
    ```bash
    pnpm run dev
    ```

## Security Notes

- The `HMAC_SECRET` environment variable is **server-side only** and will never be exposed to the client
- This secret must match the backend configuration for proper HMAC verification
- Never commit the actual secret value to version control
- Use different secrets for different environments (development, staging, production)

## Usage

### GetUserID Functionality

The frontend includes GetUserID functionality that allows authenticated users to retrieve their user ID:

```typescript
import { useGetUserID } from '@/hooks/useGetUserID';

function MyComponent() {
    const { getUserId, isLoading, error, data } = useGetUserID();

    const handleGetUserId = async () => {
        try {
            const userId = await getUserId();
            console.log('User ID:', userId);
        } catch (error) {
            console.error('Failed to get user ID:', error);
        }
    };

    return (
        <button onClick={handleGetUserId} disabled={isLoading}>
            {isLoading ? 'Loading...' : 'Get User ID'}
        </button>
    );
}
```

### How GetUserID Works

1. **Client generates timestamp** and requests HMAC signature from server
2. **Server generates HMAC** using the secret (never exposed to client)
3. **Client sends RPC request** to backend with timestamp and signature
4. **Backend validates HMAC** and returns the user ID
5. **Client receives user ID** without ever seeing the secret

This approach ensures the HMAC secret remains secure while providing the GetUserID functionality.

## Built With

- [Next.js](https://nextjs.org/)
- [Tailwind CSS](https://tailwindcss.com/)
- [Shadcn-UI](https://shadcn.com/)
- [TypeScript](https://www.typescriptlang.org/)

## Contributing

Contributions are welcome! Please follow our [contribution guidelines](/docs/_docs/community/contributions.md).

## License

This project is licensed under the [GNU Affero General Public License v3.0](/LICENSE).

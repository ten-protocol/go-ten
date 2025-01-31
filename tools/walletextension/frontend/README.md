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
   NEXT_PUBLIC_API_GATEWAY_URL=********
   ```
   
   Possible values for `NEXT_PUBLIC_API_GATEWAY_URL` are:
   - `https://uat-testnet.ten.xyz`
   - `https://sepolia-testnet.ten.xyz`
   - `https://dev-testnet.ten.xyz`

4. **Run the Development Server:**
   ```bash
   pnpm run dev
   ```

   The application will be accessible at [http://localhost:3000](http://localhost:3000).

## Usage

- Connect to Ten Testnet using the button in the top right corner of the application or on the homepage
- You can request tokens from the Discord bot by typing `!faucet <your address>` in the #faucet channel
- You can also revoke accounts by clicking the "Revoke Accounts" button on the homepage

## Built With

- [Next.js](https://nextjs.org/)
- [Tailwind CSS](https://tailwindcss.com/)
- [Shadcn-UI](https://shadcn.com/)
- [TypeScript](https://www.typescriptlang.org/)


## Contributing

Contributions are welcome! Please follow our [contribution guidelines](/docs/_docs/community/contributions.md).

## License

This project is licensed under the [GNU Affero General Public License v3.0](/LICENSE).
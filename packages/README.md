# Monorepo Structure for TEN Frontend Projects

This repository uses a monorepo setup to manage multiple frontend projects and shared resources such as components, utilities, and hooks. The structure enhances code reuse, maintainability, and collaboration across projects like **Tenscan**, **Gateway**, **Bridge**, and more.

## Folder Structure

```bash
📁 packages
├── 📁 apis - Server-side logic, API routes, and backend services
│   ├── 📁 .config - Configuration files for the APIs
│   ├── 📁 src - Source files for API logic
│   └── 📁 storage - Storage-related logic or utilities
├── 📁 eslint-config - Centralized ESLint configuration for all frontend projects
├── 📁 shared - Reusable components, hooks, and utilities shared across frontend apps
│   └── 📁 src - Main directory containing shared code
├── 📁 typescript-config - Centralized TypeScript configurations
│   ├── 📄 base.json - Base TypeScript configuration for general projects
│   ├── 📄 nextjs.json - Configuration specific to Next.js projects
│   └── 📄 react-library.json - Configuration for React libraries
└── 📁 ui - 
    ├── 📁 api - API logic consumed by the frontend
    ├── 📁 components - Reusable React components used in the UI
    ├── 📁 hooks - Custom hooks used across the frontend
    ├── 📁 lib - Utility functions used across the frontend
    ├── 📁 public - Static files such as images and assets
    ├── 📁 services - External service interactions like APIs
    ├── 📁 routes - Routing configuration and route-related logic
    └── 📁 stores - Global state mgt
```

## Getting Started

1. **Clone the Repository:**

   ```bash
   git clone https://github.com/ten-protocol/go-ten.git
   ```

2. **Install Dependencies:**

   ```bash
   pnpm install
   ```

3. **Navigate to the Project:**

   ```bash
   Tenscan: cd tools/tenscan/frontend
   Gateway: cd tools/walletextension/frontend
   Bridge: cd tools/bridge-frontend
   ```

4. **Run the Project:**

   ```bash
    pnpm dev
    ```


## Built With

- [Next.js](https://nextjs.org/)
- [Tailwind CSS](https://tailwindcss.com/)
- [TypeScript](https://www.typescriptlang.org/)

## Contributing

Contributions are welcome! Follow our [contribution guidelines](/docs/_docs/community/contributions.md).

## License

This project is licensed under the [GNU Affero General Public License v3.0](/LICENSE).
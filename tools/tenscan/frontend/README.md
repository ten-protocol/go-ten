# Tenscan: Blockchain Explorer for The Encryption Network (TEN)

Tenscan is a Next.js and Tailwind CSS-powered application that serves as a Blockchain Explorer for The Encryption Network (TEN). This explorer allows users to interact with and explore the blocks, transactions, batches, resources, and personal data on the TEN blockchain. Tenscan is built using the Shadcn-UI Component Library for a consistent and visually appealing user interface.

## Folder Structure

```
📁 Tenscan
├── 📁 api - Contains server-side code, such as API routes or server logic
├── 📁 pages - Typically used for Next.js pages. Each .tsx or .js file in this directory becomes a route in your application
├── 📁 public - This directory is used to serve static assets. Files inside this directory can be referenced in your code with a URL path
├── 📁 src - Main source code directory for this project
│   ├── 📁 components - Contains reusable React components used throughout the application
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
   cd go-ten/tools/tenscan/frontend
   ```

2. **Install Dependencies:**
   ```bash
   pnpm install
   ```

3. **Run the Development Server:**
   ```bash
   pnpm run dev
   ```

   The application will be accessible at [http://localhost:3000](http://localhost:3000).

## Usage

- Visit the different sections of the explorer through the navigation links in the UI.
- Explore the different blocks, transactions, batches, resources, and personal data on the TEN.
- View the details of each batch by their hash.

## Built With

- [Next.js](https://nextjs.org/)
- [Tailwind CSS](https://tailwindcss.com/)
- [Shadcn-UI](https://shadcn.com/)
- [TypeScript](https://www.typescriptlang.org/)


## Contributing

Contributions are welcome! Please follow our [contribution guidelines](/docs/_docs/community/contributions.md).

## License

This project is licensed under the [GNU Affero General Public License v3.0](/LICENSE).
# Tenscan: Blockchain Explorer for The Encryption Network (TEN)

Tenscan is a Next.js and Tailwind CSS-powered application that serves as a Blockchain Explorer for The Encryption Network (TEN). This explorer allows users to interact with and explore the blocks, transactions, batches, resources, and personal data on the TEN blockchain. Tenscan is built using the Shadcn-UI Component Library for a consistent and visually appealing user interface.

## Folder Structure

```
📁 Tenscan
├── 📁 api
├── 📁 pages
│   ├── 📁 batches
│   │   ├── 📄 [hash].tsx
│   │   └── 📄 index.tsx
│   ├── 📁 blocks
│   │   └── 📄 index.tsx
│   ├── 📁 personal
│   │   └── 📄 index.tsx
│   ├── 📁 resources
│   │   └── 📄 index.tsx
│   └── 📁 transactions
│       └── 📄 index.tsx
├── 📁 public
├── 📁 styles
│   ├── 📁 fonts
│   └── 📄 global.css
├── 📁 src
│   ├── 📁 components
│   │   ├── 📁 layouts
│   │   ├── 📁 modules
│   │   │   ├── 📁 batches
│   │   │   ├── 📁 blocks
│   │   │   ├── 📁 common
│   │   │   ├── 📁 dashboard
│   │   │   ├── 📁 personal
│   │   │   ├── 📁 resources
│   │   │   └── 📁 transactions
│   │   ├── 📁 providers
│   │   └── 📁 ui
│   ├── 📁 hooks
│   ├── 📁 lib
│   │   ├── 📄 constants
│   │   └── 📄 utils
│   ├── 📁 routes
│   ├── 📁 services
│   └── 📁 types
└── 📁 styles
    ├── 📁 fonts
    └── 📄 global.css
```

## Getting Started

1. **Clone the Repository:**
   ```bash
   git clone https://github.com/ten-protocol/go-ten.git
   cd go-ten/tools/obscuroscan_v3/frontend
   ```

2. **Install Dependencies:**
   ```bash
   npm install
   ```

3. **Run the Development Server:**
   ```bash
   npm run dev
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
# Monorepo Structure for TEN Frontend Projects

This repository uses a monorepo setup to manage multiple frontend projects and shared resources such as components, utilities, and hooks. The structure enhances code reuse, maintainability, and collaboration across projects like **Tenscan**, **Gateway**, **Bridge**, and more.

## Folder Structure

```bash
📁 packages
├── 📁 apis - Server-side logic, API routes, and backend services
│   ├── 📁 .config - Configuration files for the APIs
│   ├── 📁 src - Source files for API logic
│   ├── 📁 storage - Storage-related logic or utilities
├── 📁 eslint-config - Centralized ESLint configuration for all frontend projects
├── 📁 shared - Reusable components, hooks, and utilities shared across frontend apps
│   ├── 📁 src - Main directory containing shared code
├── 📁 typescript-config - Centralized TypeScript configurations
│   ├── 📄 base.json - Base TypeScript configuration for general projects
│   ├── 📄 nextjs.json - Configuration specific to Next.js projects
│   ├── 📄 react-library.json - Configuration for React libraries
├── 📁 ui - 
│   ├── 📁 api - API logic consumed by the frontend
│   ├── 📁 components - Reusable React components used in the UI
│   ├── 📁 hooks - Custom hooks used across the frontend
│   ├── 📁 lib - Utility functions used across the frontend
│   ├── 📁 public - Static files such as images and assets
│   ├── 📁 services - External service interactions like APIs
│   ├── 📁 routes - Routing configuration and route-related logic
│   ├── 📁 stores - Global state mgt
```
/** @type {import("eslint").Linter.Config} */
module.exports = {
  root: true,
  rules: {
    "no-redeclare": "off",
  },
  extends: ["@repo/eslint-config/react-internal.js"],
  parser: "@typescript-eslint/parser",
  parserOptions: {
    project: "packages/ui/tsconfig.lint.json",
  },
};
